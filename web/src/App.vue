<script setup>
import {  RouterView } from 'vue-router'
</script>

<template>
<div>
  <RouterView />
  <GlobalToast />
</div>

</template>

<script>
import GlobalToast from '@/components/Toast.vue'
export default {
  components: { GlobalToast }
}
</script>

<style scoped>

</style>
