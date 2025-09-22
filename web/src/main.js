import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router' // 引入路由配置
import { websocket } from './utils/websocket'

const app = createApp(App)

// 添加全局websocket
app.config.globalProperties.$websocket = websocket
// 连接WebSocket（根据你的实际地址修改）
websocket.connect('ws://127.0.0.1:8000/enter')

app.use(router) // 使用路由
.mount('#app')
