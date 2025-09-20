class WebSocketService {
  constructor() {
    this.socket = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.listeners = new Map()
  }

  connect(url) {
    if (this.socket) return

    this.socket = new WebSocket(url)
    
    this.socket.onopen = () => {
      this.reconnectAttempts = 0
      this.trigger('connected')
    }

    this.socket.onmessage = (event) => {
      this.trigger('message', JSON.parse(event.data))
    }

    this.socket.onclose = () => {
      this.trigger('disconnected')
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        setTimeout(() => {
          this.reconnectAttempts++
          this.connect(url)
        }, 3000)
      }
    }
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event).add(callback)
  }

  off(event, callback) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).delete(callback)
    }
  }

  trigger(event, ...args) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => callback(...args))
    }
  }

  send(data) {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(data))
    }
  }

  close() {
    this.socket?.close()
    this.socket = null
    this.reconnectAttempts = 0
  }
}

export const websocket = new WebSocketService()
