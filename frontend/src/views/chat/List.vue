<template>
  <div class="chat-list">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat" class="active">消息</router-link>
            <router-link to="/publish">发布商品</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <el-card class="chat-card">
          <!-- 标题和未读数 -->
          <div class="header">
            <h2>消息列表</h2>
            <el-badge :value="totalUnread" :hidden="totalUnread === 0">
              <span>未读消息</span>
            </el-badge>
          </div>

          <!-- 会话列表 -->
          <div class="conversation-list" v-loading="loading">
            <el-empty v-if="conversations.length === 0 && !loading" description="暂无消息" />

            <div
              v-for="conv in conversations"
              :key="conv.user_id"
              class="conversation-item"
              @click="goChat(conv)"
            >
              <el-avatar :size="50" :src="conv.avatar">
                {{ conv.username?.[0] || '?' }}
              </el-avatar>
              <div class="conv-info">
                <div class="username">{{ conv.nickname || conv.username }}</div>
                <div class="last-message">{{ conv.last_message || '暂无消息' }}</div>
                <div class="meta">
                  <span class="time">{{ formatTime(conv.update_time) }}</span>
                  <el-badge v-if="conv.unread_count > 0" :value="conv.unread_count" />
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getConversations, getUnreadCount } from '@/api/message'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const conversations = ref([])
const totalUnread = ref(0)

// 获取会话列表
const fetchConversations = async () => {
  loading.value = true
  try {
    const res = await getConversations()
    conversations.value = res.data
  } catch (error) {
    console.error('获取会话列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取未读数
const fetchUnreadCount = async () => {
  try {
    const res = await getUnreadCount()
    totalUnread.value = res.data.total
  } catch (error) {
    console.error('获取未读数失败:', error)
  }
}

// 进入聊天
const goChat = (conv) => {
  router.push(`/chat/${conv.user_id}`)
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = (now - date) / 1000 // 秒

  if (diff < 60) {
    return '刚刚'
  } else if (diff < 3600) {
    return Math.floor(diff / 60) + '分钟前'
  } else if (diff < 86400) {
    return Math.floor(diff / 3600) + '小时前'
  } else if (diff < 604800) {
    return Math.floor(diff / 86400) + '天前'
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

onMounted(() => {
  fetchConversations()
  fetchUnreadCount()

  // 定时刷新未读数
  const interval = setInterval(() => {
    fetchUnreadCount()
  }, 30000) // 30秒

  // 组件卸载时清除定时器
  return () => {
    clearInterval(interval)
  }
})
</script>

<style scoped>
.chat-list {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.el-header {
  background-color: #409eff;
  color: white;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1400px;
  margin: 0 auto;
}

.logo {
  font-size: 20px;
  margin: 0;
  cursor: pointer;
}

.nav {
  display: flex;
  gap: 20px;
}

.nav a {
  color: white;
  text-decoration: none;
  cursor: pointer;
}

.nav a.active {
  font-weight: bold;
}

.el-main {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.chat-card {
  padding: 20px;
}

.chat-card .header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.conversation-list {
  min-height: 400px;
}

.conversation-item {
  display: flex;
  gap: 15px;
  padding: 15px 10px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.conversation-item:hover {
  background-color: #f5f7fa;
}

.conv-info {
  flex: 1;
  min-width: 0;
}

.username {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 5px;
}

.last-message {
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 5px;
}

.time {
  font-size: 12px;
  color: #909399;
}
</style>
