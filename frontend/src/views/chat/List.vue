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
            <router-link to="/chat" class="active">
              <el-badge :value="userStore.unreadMessageCount" :hidden="userStore.unreadMessageCount === 0" :max="99" class="nav-badge">
                消息
              </el-badge>
            </router-link>
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
            <el-badge :value="userStore.unreadMessageCount" :hidden="userStore.unreadMessageCount === 0">
              <span class="unread-text">未读消息</span>
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
              <div class="avatar-wrapper">
                <el-avatar :size="54" :src="conv.avatar" class="conv-avatar">
                  {{ (conv.nickname || conv.username)?.[0] || '?' }}
                </el-avatar>
              </div>
              <div class="conv-content">
                <div class="conv-header">
                  <span class="username">{{ conv.nickname || conv.username }}</span>
                  <span class="time">{{ formatTime(conv.update_time) }}</span>
                </div>
                <div class="conv-footer">
                  <span class="last-message">{{ conv.last_message || '暂无消息' }}</span>
                  <div v-if="conv.unread_count > 0" class="msg-badge-wrapper">
                    <span class="custom-badge">{{ conv.unread_count > 99 ? '99+' : conv.unread_count }}</span>
                  </div>
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
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getConversations } from '@/api/message'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const conversations = ref([])

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

// 监听全局 WebSocket 传来的最新消息，以便实时刷新会话列表
watch(() => userStore.latestMessage, (newMsg) => {
  if (newMsg) {
    fetchConversations()
  }
})

// 获取未读数
const fetchUnreadCount = async () => {
  await userStore.fetchUnreadCount()
}

// 进入聊天
const goChat = (conv) => {
  // 乐观更新未读数量
  if (conv.unread_count > 0) {
    userStore.unreadMessageCount -= conv.unread_count
    conv.unread_count = 0
  }
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
  userStore.fetchUnreadCount()
})
</script>

<style scoped>
.chat-list {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
}

.el-container {
  height: 100%;
}

.el-header {
  background-color: #409eff;
  color: white;
  padding: 0 20px;
  height: 60px;
  flex-shrink: 0;
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
  width: 50%;
  min-width: 600px;
  margin: 0 auto;
  padding: 0;
  height: calc(100vh - 60px);
}

.chat-card {
  height: 100%;
  border-radius: 0;
  overflow: hidden;
  box-shadow: 0 0 16px rgba(0, 0, 0, 0.04) !important;
  display: flex;
  flex-direction: column;
  border: none;
}

/* 覆盖 el-card 默认 padding */
:deep(.el-card__body) {
  padding: 0;
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.chat-card .header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f2f5;
  background-color: #fff;
  flex-shrink: 0;
}

.chat-card .header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
  font-weight: 600;
}

.unread-text {
  color: #909399;
  font-size: 13px;
  margin-right: 8px;
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
  background: #fff;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.conversation-item::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 94px;
  right: 0;
  height: 1px;
  background-color: #f0f2f5;
}

.conversation-item:last-child::after {
  display: none;
}

.conversation-item:hover {
  background-color: #f5f7fa;
}

.avatar-wrapper {
  margin-right: 16px;
  flex-shrink: 0;
}

.conv-avatar {
  border: 1px solid #ebeef5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.conv-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.conv-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 6px;
}

.username {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 70%;
}

.time {
  font-size: 12px;
  color: #909399;
  flex-shrink: 0;
}

.conv-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.last-message {
  font-size: 14px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  margin-right: 16px;
}

.msg-badge-wrapper {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.custom-badge {
  background-color: #f56c6c;
  color: #fff;
  font-size: 12px;
  height: 18px;
  line-height: 18px;
  padding: 0 6px;
  border-radius: 9px;
  font-weight: bold;
  text-align: center;
  white-space: nowrap;
  box-shadow: 0 0 0 1px #fff;
}
</style>
