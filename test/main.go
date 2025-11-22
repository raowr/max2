package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// 全局统计信息
var (
	successCount int64
	failCount    int64
	totalSent    int64
	totalRecv    int64
	activeConns  int64
	startTime    time.Time
)

// 消息结构体
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
	Name string `json:"name"`
}

// 客户端配置
type Config struct {
	ServerAddr      string
	NumConnections  int
	Concurrency     int
	MessageInterval time.Duration
	TestDuration    time.Duration
	EnablePprof     bool
	PprofPort       int
	ClientName      string // 客户端标识
}

// 连接管理器
type ConnectionManager struct {
	config      *Config
	connections []*websocket.Conn
	mu          sync.RWMutex
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

func main() {
	config := parseFlags()

	if config.EnablePprof {
		go startPprof(config.PprofPort)
	}

	manager := NewConnectionManager(config)

	// 处理退出信号
	setupSignalHandler(manager)

	// 启动测试
	if err := manager.Run(); err != nil {
		log.Fatalf("测试失败: %v", err)
	}
}

// 解析命令行参数
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.ServerAddr, "server", "ws://localhost:8080/ws", "WebSocket服务器地址")
	flag.IntVar(&config.NumConnections, "connections", 1000, "要创建的连接数量")
	flag.IntVar(&config.Concurrency, "concurrency", 100, "并发创建连接数")
	flag.DurationVar(&config.MessageInterval, "interval", 10*time.Second, "消息发送间隔")
	flag.DurationVar(&config.TestDuration, "duration", 1*time.Minute, "测试持续时间")
	flag.BoolVar(&config.EnablePprof, "pprof", false, "启用性能分析")
	flag.IntVar(&config.PprofPort, "pprof-port", 6060, "pprof端口")
	flag.StringVar(&config.ClientName, "name", "ws-client", "客户端标识名称")

	flag.Parse()
	return config
}

// 启动性能分析
func startPprof(port int) {
	log.Printf("启动性能分析服务器 :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// 创建连接管理器
func NewConnectionManager(config *Config) *ConnectionManager {
	return &ConnectionManager{
		config:      config,
		connections: make([]*websocket.Conn, 0, config.NumConnections),
		stopChan:    make(chan struct{}),
	}
}

// 设置信号处理器
func setupSignalHandler(manager *ConnectionManager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("收到退出信号，正在关闭连接...")
		close(manager.stopChan)
		manager.Stop()
		os.Exit(0)
	}()
}

// 运行测试
func (m *ConnectionManager) Run() error {
	startTime = time.Now()
	log.Printf("开始WebSocket压力测试")
	log.Printf("目标服务器: %s", m.config.ServerAddr)
	log.Printf("连接数量: %d", m.config.NumConnections)
	log.Printf("并发数: %d", m.config.Concurrency)
	log.Printf("测试时长: %v", m.config.TestDuration)
	log.Printf("客户端名称: %s", m.config.ClientName)
	log.Printf("消息格式: %s", `{ type: 'heartbeat', data: 'ping' }`)

	// 启动统计报告
	go m.reportStats()

	// 创建连接
	if err := m.createConnections(); err != nil {
		return err
	}

	log.Printf("成功创建 %d 个连接，开始消息测试", atomic.LoadInt64(&successCount))

	// 启动消息发送
	m.startMessageSending()

	// 等待测试完成或停止信号
	select {
	case <-time.After(m.config.TestDuration):
		log.Println("测试时间到，正在停止...")
	case <-m.stopChan:
		log.Println("收到停止信号，正在停止...")
	}

	m.Stop()
	return nil
}

// 创建连接
func (m *ConnectionManager) createConnections() error {
	log.Println("开始创建连接...")

	connChan := make(chan int, m.config.NumConnections)
	for i := 0; i < m.config.NumConnections; i++ {
		connChan <- i
	}
	close(connChan)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, m.config.Concurrency)

	for i := 0; i < m.config.Concurrency; i++ {
		wg.Add(1)
		go m.connectionWorker(connChan, semaphore, &wg)
	}

	wg.Wait()
	log.Printf("连接创建完成: 成功=%d, 失败=%d", successCount, failCount)

	if atomic.LoadInt64(&successCount) == 0 {
		return fmt.Errorf("所有连接尝试都失败了")
	}

	return nil
}

// 连接工作协程
func (m *ConnectionManager) connectionWorker(connChan <-chan int, semaphore chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range connChan {
		semaphore <- struct{}{}

		conn, err := m.createSingleConnection(id)
		if err != nil {
			atomic.AddInt64(&failCount, 1)
			if atomic.LoadInt64(&failCount)%100 == 0 {
				log.Printf("连接失败数达到 %d", atomic.LoadInt64(&failCount))
			}
		} else {
			atomic.AddInt64(&successCount, 1)
			m.mu.Lock()
			m.connections = append(m.connections, conn)
			m.mu.Unlock()
			atomic.AddInt64(&activeConns, 1)
		}

		<-semaphore

		// 显示进度
		total := atomic.LoadInt64(&successCount) + atomic.LoadInt64(&failCount)
		if total%1000 == 0 {
			elapsed := time.Since(startTime)
			rate := float64(total) / elapsed.Seconds()
			log.Printf("进度: %d/%d (%.1f conn/s)", total, m.config.NumConnections, rate)
		}
	}
}

// 创建单个连接
func (m *ConnectionManager) createSingleConnection(id int) (*websocket.Conn, error) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}

	conn, _, err := dialer.Dial(m.config.ServerAddr, nil)
	if err != nil {
		return nil, err
	}

	// 启动消息接收协程
	go m.receiveMessages(conn, id)

	return conn, nil
}

// 启动消息发送
func (m *ConnectionManager) startMessageSending() {
	ticker := time.NewTicker(m.config.MessageInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.broadcastMessage()
		case <-m.stopChan:
			return
		}
	}
}

// 生成指定的消息内容
func generateMessage(clientName string, connID int) ([]byte, error) {
	msg := Message{
		Type: "heartbeat",
		Data: "ping",
		Name: fmt.Sprintf("%s-conn-%d", clientName, connID),
	}

	return json.Marshal(msg)
}

// 广播消息到所有连接
func (m *ConnectionManager) broadcastMessage() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.connections) == 0 {
		return
	}

	var wg sync.WaitGroup

	// 限制并发发送数
	semaphore := make(chan struct{}, 1000)

	for i, conn := range m.connections {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(c *websocket.Conn, connID int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			message, err := generateMessage(m.config.ClientName, connID)
			if err != nil {
				log.Printf("生成消息失败: %v", err)
				return
			}

			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				atomic.AddInt64(&failCount, 1)
				log.Printf("发送消息失败: %v", err)
			} else {
				atomic.AddInt64(&totalSent, 1)
			}
		}(conn, i)
	}

	wg.Wait()
}

// 接收消息
func (m *ConnectionManager) receiveMessages(conn *websocket.Conn, id int) {
	defer atomic.AddInt64(&activeConns, -1)

	for {
		select {
		case <-m.stopChan:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("连接 %d 读取消息失败: %v", id, err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				atomic.AddInt64(&totalRecv, 1)

				// 可选：记录接收到的消息内容
				if atomic.LoadInt64(&totalRecv)%1000 == 0 {
					log.Printf("接收到消息: %s", string(message))
				}
			}
		}
	}
}

// 停止所有连接
func (m *ConnectionManager) Stop() {
	log.Println("正在关闭所有连接...")

	m.mu.Lock()
	defer m.mu.Unlock()

	var wg sync.WaitGroup
	for _, conn := range m.connections {
		wg.Add(1)
		go func(c *websocket.Conn) {
			defer wg.Done()
			c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			c.Close()
		}(conn)
	}

	wg.Wait()
	m.connections = m.connections[:0]
	log.Println("所有连接已关闭")

	// 打印最终统计
	m.printFinalStats()
}

// 报告统计信息
func (m *ConnectionManager) reportStats() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			success := atomic.LoadInt64(&successCount)
			fail := atomic.LoadInt64(&failCount)
			active := atomic.LoadInt64(&activeConns)
			sent := atomic.LoadInt64(&totalSent)
			recv := atomic.LoadInt64(&totalRecv)

			log.Printf("状态: 成功=%d, 失败=%d, 活跃=%d, 发送=%d, 接收=%d, 时长=%.1fs",
				success, fail, active, sent, recv, elapsed)

			if elapsed > 0 {
				log.Printf("速率: 连接=%.1f/s, 消息=%.1f/s",
					float64(success+fail)/elapsed, float64(sent)/elapsed)
			}

		case <-m.stopChan:
			return
		}
	}
}

// 打印最终统计信息
func (m *ConnectionManager) printFinalStats() {
	totalTime := time.Since(startTime)
	success := atomic.LoadInt64(&successCount)
	fail := atomic.LoadInt64(&failCount)
	sent := atomic.LoadInt64(&totalSent)
	recv := atomic.LoadInt64(&totalRecv)

	fmt.Println("\n=== 测试结果统计 ===")
	fmt.Printf("测试时长: %v\n", totalTime)
	fmt.Printf("连接统计: 成功=%d, 失败=%d, 成功率=%.2f%%\n",
		success, fail, float64(success)/float64(success+fail)*100)
	fmt.Printf("消息统计: 发送=%d, 接收=%d, 收发比=%.2f%%\n",
		sent, recv, float64(recv)/float64(sent)*100)
	fmt.Printf("平均速率: %.1f 消息/秒\n", float64(sent)/totalTime.Seconds())
}
