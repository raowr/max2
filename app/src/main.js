import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router' // 引入路由配置
import { websocket } from './utils/websocket'
import { storage } from './utils/storage' // 导入封装的存储工具

const app = createApp(App)

// 添加全局websocket
app.config.globalProperties.$websocket = websocket

// 确保 DOM 完全加载后再初始化
function initApp() {
  // 使用路由
  app.use(router)
  
  // 挂载应用
  app.mount('#app')
  
  // 连接WebSocket（根据你的实际地址修改）
  const wsUrl = import.meta.env.VITE_WS_URL;
  console.log('WebSocket URL:', wsUrl);
  
  // 延迟连接 WebSocket，确保应用先加载
  setTimeout(() => {
    if (wsUrl) {
      console.log('user_id:', storage.local.get("user_id"))
// 初始化 WebSocket 配置（全局一次）
      websocket.init({
        url: wsUrl+"?user_id="+storage.local.get("user_id"), // 后端地址
        reconnectInterval: 5000, // 5秒重连一次
        heartbeatInterval: 10000 // 20秒一次心跳
      });

// 启动连接（可在登录后再调用，这里直接启动作为示例）
      websocket.connect();
    } else {
      console.warn('VITE_WS_URL is not defined, using default WebSocket URL')
      // 可以在这里设置一个默认的 WebSocket URL
      websocket.connect('ws://127.0.0.1:8000/enter')
    }
  }, 100)
}

// 根据环境选择初始化方式
if (window.plus) {
  // 在 5+ Runtime 环境中
  document.addEventListener('plusready', initApp)
} else {
  // 在普通浏览器环境中
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initApp)
  } else {
    initApp()
  }
}

// 添加全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue error:', err)
  console.log('Instance:', instance)
  console.log('Info:', info)
}

// 可选：添加全局未处理 Promise  rejection 处理
window.addEventListener('unhandledrejection', (event) => {
  console.error('Unhandled promise rejection:', event.reason)
})