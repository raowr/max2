// src/utils/websocket.js
class WebSocketService {
  constructor() {
    if (WebSocketService.instance) {
      return WebSocketService.instance;
    }
    WebSocketService.instance = this;

    // 初始化属性
    this.url = '';
    this.ws = null;
    this.isConnecting = false;
    this.reconnectInterval = 3000; // 重连间隔
    this.maxReconnectAttempts = -1; // 无限重连
    this.reconnectCount = 0;
    this.heartbeatInterval = 30000; // 心跳间隔（30秒）
    this.heartbeatTimer = null;
    // 回调函数队列（支持多组件注册）
    this.callbacks = {
      open: [],
      message: [],
      error: [],
      close: []
    };
  }

  // 初始化配置
  init(options) {
    this.url = options.url;
    this.reconnectInterval = options.reconnectInterval || this.reconnectInterval;
    this.maxReconnectAttempts = options.maxReconnectAttempts ?? this.maxReconnectAttempts;
    this.heartbeatInterval = options.heartbeatInterval || this.heartbeatInterval;
    return this;
  }

  // 启动连接
  connect() {
    if (this.isConnecting || !this.url) return;
    this.isConnecting = true;

    try {
      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        console.log('WebSocket 连接成功');
        this.isConnecting = false;
        this.reconnectCount = 0;
        this.startHeartbeat();
        this.callbacks.open.forEach(cb => cb());
      };

      this.ws.onmessage = (event) => {
        this.callbacks.message.forEach(cb => cb(event.data));
        this.resetHeartbeat();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket 错误:', error);
        this.isConnecting = false;
        this.callbacks.error.forEach(cb => cb(error));
        this.handleReconnect();
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket 关闭，代码:', event.code);
        this.isConnecting = false;
        this.callbacks.close.forEach(cb => cb(event));
        this.stopHeartbeat();
        if (event.code !== 1000) this.handleReconnect();
      };
    } catch (error) {
      console.error('连接失败:', error);
      this.isConnecting = false;
      this.handleReconnect();
    }
  }

  // 重连逻辑
  handleReconnect() {
    if (this.maxReconnectAttempts !== -1 && this.reconnectCount >= this.maxReconnectAttempts) {
      console.log('达到最大重连次数，停止重试');
      return;
    }
    this.reconnectCount++;
    console.log(`第 ${this.reconnectCount} 次重连...`);
    setTimeout(() => this.connect(), this.reconnectInterval);
  }

  // 心跳检测
  startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.send({ type: 'heartbeat', data: 'ping' });
      }
    }, this.heartbeatInterval);
  }

  resetHeartbeat() {
    this.stopHeartbeat();
    this.startHeartbeat();
  }

  stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  // 发送消息
  send(data) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    } else {
      console.warn('WebSocket 未连接，无法发送消息');
    }
  }

  // 手动关闭连接
  close() {
    if (this.ws) {
      this.ws.close(1000, '手动关闭');
      this.ws = null;
    }
  }

  // 注册回调
  on(event, callback) {
    if (this.callbacks[event]) {
      this.callbacks[event].push(callback);
    }
  }

  // 移除回调（避免内存泄漏）
  off(event, callback) {
    if (this.callbacks[event]) {
      this.callbacks[event] = this.callbacks[event].filter(cb => cb !== callback);
    }
  }
}

// 导出单例实例
export const websocket = new WebSocketService();