<template>
  <div class="app-container">
    <header class="app-header">
      <h1>{{ appName }}</h1>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const appName = ref('My Vue App')
</script>

<style scoped>
.app-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f5f5;
}

.app-header {
  background: #4CAF50;
  color: white;
  padding: 12px 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  /* 适配刘海屏 */
  padding-top: max(12px, env(safe-area-inset-top));
}

.app-header h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.app-main {
  flex: 1;
  overflow-y: auto;
  /* 适配底部安全区域 */
  padding-bottom: env(safe-area-inset-bottom);
}
</style>