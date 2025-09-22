import { createRouter, createWebHistory } from 'vue-router'
import Index from '../views/Index.vue'
import Room from '../views/Room.vue'
import Game from '../views/Game.vue'

const routes = [
      {
    path: '/',  // 根路径
    name: 'index',
    component: Index  // 对应到首页组件
  },
  {
    path: '/room',
    name: 'Room',
    component: Room
  },
    {
    path: '/game',
    name: 'Game',
    component: Game
  }
  // 其他路由...
]

const router = createRouter({
  // Vite 中使用 import.meta.env.BASE_URL
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router