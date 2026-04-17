<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'

const route = useRoute()
const userStore = useUserStore()

// 监听路由变化，自动刷新未读消息数
watch(
  () => route.path,
  () => {
    if (userStore.isLogin) {
      userStore.fetchUnreadCount()
    }
  }
)

onMounted(() => {
  if (userStore.isLogin) {
    userStore.fetchUnreadCount()
    userStore.initWebSocket()
  }
  
  // 定时刷新未读消息数
  setInterval(() => {
    if (userStore.isLogin) {
      userStore.fetchUnreadCount()
    }
  }, 30000)
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  min-height: 100vh;
}
</style>
