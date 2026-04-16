<template>
  <div class="chat-room">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat" class="active">消息</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <div class="chat-container">
          <!-- 聊天头部 -->
          <div class="chat-header">
            <el-button @click="goBack" :icon="ArrowLeft" circle />
            <el-avatar :size="40" :src="otherUser.avatar">
              {{ otherUser.username?.[0] }}
            </el-avatar>
            <div class="user-info">
              <h3>{{ otherUser.nickname || otherUser.username }}</h3>
              <el-tag v-if="isOnline" type="success" size="small">在线</el-tag>
              <el-tag v-else type="info" size="small">离线</el-tag>
            </div>
          </div>

          <!-- 消息列表 -->
          <div class="messages" ref="messagesContainer">
            <div
              v-for="message in messageList"
              :key="message.id"
              :class="['message-item', message.sender_id === currentUserId ? 'sent' : 'received']"
            >
              <el-avatar
                v-if="message.sender_id !== currentUserId"
                :size="36"
                :src="message.sender?.avatar"
                class="message-avatar"
              >
                {{ message.sender?.username?.[0] }}
              </el-avatar>
              <div class="message-content">
                <div class="message-bubble">
                  {{ message.content }}
                </div>
                <div class="message-time">{{ formatTime(message.create_time) }}</div>
              </div>
            </div>
          </div>

          <!-- 输入区域 -->
          <div class="input-area">
            <el-input
              v-model="messageContent"
              type="textarea"
              :rows="2"
              placeholder="输入消息..."
              @keydown.enter.exact="handleSend"
              resize="none"
            />
            <el-button
              type="primary"
              @click="handleSend"
              :disabled="!messageContent.trim()"
              :loading="sending"
            >
              发送
            </el-button>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getMessageList, sendMessage } from '@/api/message'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const currentUserId = computed(() => userStore.userId)

const otherUserId = ref(null)
const otherUser = ref({})
const isOnline = ref(false)
const messageList = ref([])
const messageContent = ref('')
const sending = ref(false)
const messagesContainer = ref(null)

// 获取消息列表
const fetchMessages = async () => {
  try {
    const res = await getMessageList(otherUserId.value)
    messageList.value = res.data.list

    // 滚动到底部
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('获取消息失败:', error)
  }
}

// 发送消息
const handleSend = async () => {
  if (!messageContent.value.trim()) {
    return
  }

  sending.value = true
  try {
    // 先通过HTTP API发送（确保消息保存）
    const res = await sendMessage({
      receiver_id: otherUserId.value,
      content: messageContent.value.trim(),
      type: 1
    })

    // 添加到消息列表
    messageList.value.push({
      id: res.data.message_id,
      sender_id: currentUserId.value,
      receiver_id: otherUserId.value,
      content: messageContent.value.trim(),
      type: 1,
      is_read: 0,
      create_time: Date.now() / 1000
    })

    messageContent.value = ''
    nextTick(() => {
      scrollToBottom()
    })
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error('发送失败')
  } finally {
    sending.value = false
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 返回
const goBack = () => {
  router.back()
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
  } else {
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

onMounted(() => {
  otherUserId.value = parseInt(route.params.userId)
  fetchMessages()
})

onUnmounted(() => {
  // 清理
})
</script>

<style scoped>
.chat-room {
  height: 100vh;
  display: flex;
  flex-direction: column;
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
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px 20px;
  background: white;
  border-bottom: 1px solid #ebeef5;
}

.user-info h3 {
  margin: 0 0 5px 0;
  font-size: 16px;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.message-item {
  display: flex;
  margin-bottom: 20px;
  align-items: flex-start;
}

.message-item.sent {
  flex-direction: row-reverse;
}

.message-item.received .message-avatar {
  margin-right: 10px;
}

.message-item.sent .message-avatar {
  margin-left: 10px;
}

.message-content {
  max-width: 60%;
  display: flex;
  flex-direction: column;
}

.message-item.sent .message-content {
  align-items: flex-end;
}

.message-bubble {
  padding: 10px 15px;
  border-radius: 8px;
  word-break: break-word;
}

.message-item.received .message-bubble {
  background-color: white;
  border: 1px solid #ebeef5;
}

.message-item.sent .message-bubble {
  background-color: #409eff;
  color: white;
}

.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.input-area {
  padding: 15px 20px;
  background: white;
  border-top: 1px solid #ebeef5;
  display: flex;
  gap: 10px;
}

.input-area .el-textarea {
  flex: 1;
}
</style>
