import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import Loading from '../views/Loading.vue'
import Index from '../views/Index.vue'
import Room from '../views/Room.vue'
import Game from '../views/Download.vue'
import Download from '../views/Game.vue'

const routes = [
    {
      path: '/',  // 根路径
      name: 'Loading',
      component: Loading  // 对应到首页组件
    },
    {
      path: '/index',  // 根路径
      name: 'Index',
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
  },
    {
    path: '/download',
    name: 'Download',
    component: Download
  }
  // 其他路由...
]

const router = createRouter({
  // Vite 中使用 import.meta.env.BASE_URL
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes
})

export default router