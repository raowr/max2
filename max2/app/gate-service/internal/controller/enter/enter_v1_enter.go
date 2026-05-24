package enter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	actionv1 "gate-service/api/action/v1"
	v1 "gate-service/api/enter/v1"

	"gate-service/internal/consts"
	"gate-service/internal/controller/action"
	"gate-service/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

func (c *ControllerV1) Enter(ctx context.Context, req *v1.EnterReq) (res *v1.EnterRes, err error) {
	var (
		r          = g.RequestFromCtx(ctx)
		ws         *websocket.Conn
		wsUpGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 从配置读取 debug 模式（推荐：配置文件中配置 app.debug = true/false）
				// MustGet: 键不存在时 panic；若想避免 panic，用 Get().Bool() 并处理默认值
				//isDebug := g.Cfg().MustGet(ctx, "app.debug").Bool()
				//if isDebug {
				//	return true // 调试模式允许所有跨域
				//}
				//// 2. Android 原生客户端（Origin 为空字符串）
				//origin := r.Header.Get("Origin")
				//allowedOrigins := map[string]bool{
				//	allowedOrigin: true, // 浏览器前端域名
				//	"": true,            // 允许空 Origin（Android 原生）
				//}
				//return allowedOrigins[origin]
				return true
			},
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				g.Log().Errorf(ctx, "WebSocket升级失败: %v", reason)
				w.WriteHeader(status)
				w.Write([]byte(reason.Error()))
			},
		}
	)

	// 获取或生成用户ID
	userName := r.GetQuery("user_name").String()
	token := r.GetQuery("token").String()

	//通过username从缓存获取token
	// userToken := service.Cache().Get(ctx, userName).String()
	userInfo := service.Cache().Get(ctx, "user:"+userName).Bytes()
	entUser := Users{}
	if userInfo == nil {
		g.Log().Errorf(ctx, "用户不存在")
		return
	}
	//缓存有，验证密码正确
	if err = json.Unmarshal(userInfo, &entUser); err != nil {
		g.Log().Errorf(ctx, "用户信息json错误")
		return
	}
	if entUser.Token != token {
		g.Log().Errorf(ctx, "token不存在")
		ws.Close()
		return
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JwtKey), nil
	})
	if err != nil || !jwtToken.Valid {
		g.Log().Errorf(ctx, "token验证失败: %v", err)
		ws.Close()
		return
	}

	// 升级HTTP连接为WebSocket
	ws, err = wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Errorf(ctx, "WebSocket升级失败: %v", err)
		r.Response.WriteHeader(http.StatusInternalServerError)
		r.Response.Write([]byte("连接建立失败"))
		return
	}

	// 创建连接专用上下文
	connCtx, cancel := context.WithCancel(context.Background())

	// 创建客户端实例（增大消息缓冲）
	client := &Client{
		conn:      ws,
		userName:  userName, // 初始化关闭标记为0（未关闭）
		cancel:    cancel,
		subClient: g.Redis(),
		pubClient: g.Redis(),
		heartbeat: time.Now(),
	}

	// 加锁更新客户端连接（重连逻辑）
	clientsMu.Lock()
	if oldClient, err := getClient(userName); err == nil && oldClient != nil {
		if oldClient.cancel != nil {
			oldClient.cancel() // 触发旧连接的 <-ctx.Done()
		}
		if oldClient.conn != nil {
			oldClient.conn.Close() // 关闭旧连接
		}

		g.Log().Infof(ctx, "用户 %s 重连，关闭旧连接", userName)
	}

	if err = addClient(client); err != nil {
		g.Log().Errorf(ctx, "添加客户端到缓存失败: %v", err)
		ws.Close()
		return
	}
	clientsMu.Unlock()

	g.Log().Infof(ctx, "用户 %s 连接成功", userName)

	// 启动协程（增加退出日志）
	go client.readLoop(connCtx)
	go client.writeLoop(connCtx)
	go client.heartbeatCheck(connCtx)

	// 返回空响应
	return &v1.EnterRes{}, nil
}

// 读消息循环：处理客户端消息
func (c *Client) readLoop(ctx context.Context) {
	defer c.closeConnection("读循环退出")

	for {
		// 监听 ctx 是否已关闭 (调用 Close 方法时触发)
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "Casbin 停止监听")
			return
		default:
			// 继续执行后续功能
		}
		// 读取客户端消息
		mt, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				g.Log().Errorf(ctx, "用户 %s 消息读取错误: %v", c.userName, err)
			}
			g.Log().Infof(ctx, "用户 %s 连接已关闭，读循环退出", c.userName)
			return // 连接断开，退出循环
		}

		// 只要成功接收到消息，无论后续处理是否成功，都更新心跳
		c.heartbeat = time.Now()

		// 解析消息（严格错误处理）
		var msg *actionv1.SendActionReq
		if err := gconv.Struct(data, &msg); err != nil {
			errMsg := &actionv1.SendActionReq{
				Type: consts.ChatTypeError,
				Data: fmt.Sprintf("无效消息格式: %v", err),
				From: c.userName,
			}
			c.safeSendMessage(ctx, errMsg)
			continue
		}
		g.Log().Infof(ctx, "用户 %s 接收消息: %s（类型: %d）", c.userName, data, mt)
		// 根据消息类型处理业务
		switch msg.Type {
		case consts.Heartbeat:
			// 1. 心跳消息在这里直接处理
			c.handleHeartbeat(ctx, msg.Data)

		case "":
			// 2. 没有类型时返回错误
			errMsg := &actionv1.SendActionReq{
				Type: consts.ChatTypeError,
				Data: "消息类型为空",
				From: c.userName,
			}
			c.safeSendMessage(ctx, errMsg)

		default:
			// 3. 其他消息类型发送到 sendAction
			// 验证是否为有效类型
			validTypes := map[string]bool{
				consts.InitRoom:   true,
				consts.CreateRoom: true,
				consts.JoinRoom:   true,
				consts.LeaveRoom:  true,
				consts.Play:       true,
				consts.PlayCard:   true,
				consts.GetInfo:    true,
			}

			if validTypes[msg.Type] {
				msg.From = c.userName
				// 有效类型：发送到 sendAction
				action.SendAction(msg)
			} else {
				// 无效类型：返回错误
				errMsg := &actionv1.SendActionReq{
					Type: consts.ChatTypeError,
					Data: "未知消息类型: " + msg.Type,
					From: c.userName,
				}
				c.safeSendMessage(ctx, errMsg)
			}
		}
	}
}

// 写消息循环：发送消息到客户端
func (c *Client) writeLoop(ctx context.Context) {
	defer g.Log().Infof(ctx, "用户 %s writeLoop退出", c.userName)
	var err error
	c.sub, _, err = c.subClient.Subscribe(ctx, consts.PlayerMsgPrefix+c.userName)
	if err != nil {
		g.Log().Error(ctx, "writeLoop 订阅失败:", err)
		return
	}
	// 确保在函数退出时关闭订阅
	// defer func() {
	// 	_ = c.sub.Close(ctx)
	// }()

	// 循环接收消息
	for {
		// 监听 ctx 是否已关闭 (调用 Close 方法时触发)
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "writeLoop 停止监听")
			return
		default:
			// 继续执行后续功能
		}

		// Receive 接收消息
		msg, err := c.sub.Receive(ctx)
		if err != nil {
			// 如果是 Context 被取消导致的错误，直接退出
			if ctx.Err() != nil {
				g.Log().Info(ctx, "ctx Watcher 接收消息错误:", ctx.Err())
				return
			}
			// 其它错误（如网络断开），可以考虑简单的重试或记录日志
			g.Log().Error(ctx, "Casbin Watcher 接收消息错误:", err)
			// 如果出错直接退出，等待下一次重启或手动干预
			return
		}
		msgData, success := ParseRedisSubscribeMessage(msg.String(), c.userName, ctx)
		if !success {
			continue
		}

		// 检查连接是否已关闭
		c.mutex.RLock()
		conn := c.conn
		c.mutex.RUnlock()
		if conn == nil {
			g.Log().Infof(ctx, "用户 %s 连接已关闭，跳过消息发送", c.userName)
			return
		}
		msgDataBytes, err := gjson.Encode(msgData)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", c.userName, err)
		}
		g.Log().Infof(ctx, "用户 %s 发送消息: %s", c.userName, string(msgDataBytes))
		if err := c.conn.WriteMessage(ghttp.WsMsgText, msgDataBytes); err != nil {
			g.Log().Errorf(ctx, " writeLoop 用户 %s 消息发送失败: %v,消息内容: %s", c.userName, err, string(msgDataBytes))
			return
		}
	}
}

// 心跳检测：超时断开连接
func (c *Client) heartbeatCheck(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	defer g.Log().Infof(ctx, "用户 %s 心跳检测退出", c.userName)
	g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
	for {
		select {
		case <-ticker.C:
			g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
			if time.Since(c.heartbeat) > 60*time.Second {
				g.Log().Infof(ctx, "用户 %s 心跳超时（60秒），断开连接", c.userName)
				//c.conn.Close()
				c.closeConnection("心跳超时")
				return
			}
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s 心跳检测被取消: %v", c.userName, ctx.Err())
			return
		}
	}
}

// 处理心跳
func (c *Client) handleHeartbeat(ctx context.Context, data string) {
	if data == "ping" {
		msgData := &actionv1.SendActionReq{
			Type: consts.Heartbeat,
			Data: "pong",
			From: c.userName,
		}
		c.safeSendMessage(ctx, msgData)
	}
}

// 安全发送消息的函数
func (c *Client) safeSendMessage(ctx context.Context, msg *actionv1.SendActionReq) {
	_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+c.userName, c.encodeMessage(ctx, msg))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", c.userName, err)
	}
}

// 消息编码（处理错误）
func (c *Client) encodeMessage(ctx context.Context, msg *actionv1.SendActionReq) []byte {
	msgBytes, err := gjson.Encode(msg)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", c.userName, err)
		return []byte(`{"type":"error","data":"消息编码失败"}`)
	}
	return msgBytes
}

// 关闭客户端连接
func (c *Client) closeConnection(reason string) {

	// 原子操作检查是否已关闭
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		g.Log().Infof(context.Background(), "用户 %s 连接已关闭，跳过重复关闭", c.userName)
		return
	}

	rmMu.Lock()
	defer rmMu.Unlock()
	logCtx := context.Background()
	g.Log().Infof(logCtx, "用户 %s 关闭连接，原因: %s", c.userName, reason)

	// 关闭订阅者
	if c.sub != nil {
		_ = c.sub.Close(logCtx)
		c.sub = nil
	}

	// 先取消上下文
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil // 设置为 nil，避免重复调用
	}

	// 然后关闭连接（防止重复关闭）
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil // 关键：设置为 nil，防止重复关闭
	}

	// 最后清理资源
	clientsMu.Lock()
	err := removeClient(c.userName)
	if err != nil {
		g.Log().Infof(logCtx, "用户 %s 删除客户端失败", c.userName)
	}
	clientsMu.Unlock()
}
