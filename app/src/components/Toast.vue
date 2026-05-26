<template>
  <div v-if="visible" class="global-toast" :class="position">
    {{ message }}
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const visible = ref(false)
const message = ref('游戏一定很精彩，但是今日你在线已达到8小时,请注意休息！')
const position = ref('top') // 'top', 'center'

let timer = null

const show = (text, duration = 4000, pos = 'top') => {
  if (timer) clearTimeout(timer)
  message.value = text
  position.value = pos
  visible.value = true
  timer = setTimeout(() => {
    visible.value = false
  }, duration)
}

// 监听全局事件
onMounted(() => {
  window.addEventListener('show-toast', (e) => {
    show(e.detail.text, e.detail.duration, e.detail.position)
  })
})

onUnmounted(() => {
  window.removeEventListener('show-toast')
})
</script>

<style scoped>
.global-toast {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  color: gold;
  padding: 14px 24px;
  border-radius: 12px;
  z-index: 9999;
  white-space: nowrap;
  font-size: 15px;
  font-weight: 500;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.35), 0 4px 12px rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: toastInOut 3s forwards;
}
.global-toast.top {
  top: 32px;
}
.global-toast.center {
  top: 50%;
  transform: translate(-50%, -50%);
}
@keyframes toastInOut {
  0% {
    opacity: 0;
    transform: translateX(-50%) translateY(-20px) scale(0.9);
  }
  15% {
    opacity: 1;
    transform: translateX(-50%) translateY(0) scale(1);
  }
  85% {
    opacity: 1;
    transform: translateX(-50%) translateY(0) scale(1);
  }
  100% {
    opacity: 0;
    transform: translateX(-50%) translateY(-10px) scale(0.95);
    visibility: hidden;
  }
}
</style>