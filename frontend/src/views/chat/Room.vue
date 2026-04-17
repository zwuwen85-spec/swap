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
              <!-- 对方头像 -->
              <el-avatar
                v-if="message.sender_id !== currentUserId"
                :size="40"
                :src="otherUser.avatar"
                class="message-avatar"
              >
                {{ (otherUser.nickname || otherUser.username)?.[0] || '?' }}
              </el-avatar>

              <div class="message-content">
                <div class="message-bubble">
                  {{ message.content }}
                </div>
                <div class="message-time">{{ formatTime(message.create_time) }}</div>
              </div>

              <!-- 我方头像 -->
              <el-avatar
                v-if="message.sender_id === currentUserId"
                :size="40"
                :src="userStore.avatar"
                class="message-avatar my-avatar"
              >
                {{ (userStore.username)?.[0] || '我' }}
              </el-avatar>
            </div>
          </div>

          <!-- 输入区域 -->
          <div class="input-area">
            <el-input
              v-model="messageContent"
              type="textarea"
              :rows="3"
              placeholder="请输入消息..."
              @keydown.enter.exact.prevent="handleSend"
              resize="none"
            />
            <div class="input-actions">
              <el-button
                type="primary"
                @click="handleSend"
                :disabled="!messageContent.trim()"
                :loading="sending"
              >
                发送(Enter)
              </el-button>
            </div>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getMessageList, sendMessage, checkOnline } from '@/api/message'
import { getUserInfoById } from '@/api/user'
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

// 监听全局 WebSocket 传来的最新消息
watch(() => userStore.latestMessage, (newMsg) => {
  if (newMsg && (newMsg.sender_id === otherUserId.value || newMsg.receiver_id === otherUserId.value)) {
    // 将新消息追加到当前聊天窗口
    messageList.value.push({
      id: newMsg.id || Date.now(), // WebSocket消息可能没有完整ID，用时间戳兜底
      sender_id: newMsg.sender_id,
      receiver_id: newMsg.receiver_id,
      content: newMsg.content,
      type: newMsg.type || 1,
      is_read: 1, // 当前正在聊天界面，视为已读
      create_time: newMsg.timestamp || Date.now() / 1000
    })
    
    // 乐观扣除全局未读数
    if (userStore.unreadMessageCount > 0) {
      userStore.unreadMessageCount -= 1
    }

    nextTick(() => {
      scrollToBottom()
    })
  }
})

// 监听全局 WebSocket 传来的最新在线状态
watch(() => userStore.onlineStatus, (status) => {
  if (status && status.user_id === otherUserId.value) {
    isOnline.value = status.online
  }
}, { deep: true })

// 监听路由参数变化，确保能重新获取数据
watch(() => route.params.userId, (newId) => {
  if (newId) {
    otherUserId.value = parseInt(newId)
    initRoom()
  }
})

const initRoom = async () => {
  // 1. 获取对方用户信息
  try {
    const userRes = await getUserInfoById(otherUserId.value)
    if (userRes.code === 0 && userRes.data) {
      otherUser.value = userRes.data
    }
  } catch (err) {
    console.error('获取对方用户信息失败:', err)
  }

  // 2. 获取对方在线状态
  try {
    const onlineRes = await checkOnline(otherUserId.value)
    if (onlineRes.code === 0) {
      isOnline.value = onlineRes.data.is_online
    }
  } catch (err) {
    console.error('获取对方在线状态失败:', err)
  }

  // 3. 获取消息历史
  fetchMessages()
}

// 获取消息列表
const fetchMessages = async () => {
  try {
    const res = await getMessageList(otherUserId.value)
    // 后端返回的是按时间倒序，为了在聊天窗口正常显示，需要反转为正序（最新在最下）
    const list = res.data.list ? res.data.list.reverse() : []
    messageList.value = list

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
  initRoom()
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
  overflow: hidden;
}

.el-container {
  height: 100%;
  overflow: hidden;
}

.el-header {
  background-color: #409eff;
  color: white;
  padding: 0 20px;
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
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0;
  width: 50%;
  min-width: 600px;
  margin: 0 auto;
  overflow: hidden;
  min-height: 0;
}

.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f5f7fa;
  min-height: 0;
  box-shadow: 0 0 16px rgba(0, 0, 0, 0.04);
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px 30px;
  background: white;
  border-bottom: 1px solid #ebeef5;
  flex-shrink: 0;
  position: relative; /* 确保返回按钮可以相对此容器绝对定位 */
}

.chat-header .el-button {
  /* 返回按钮固定在最左侧的留白区域中 */
  position: absolute;
  left: -5px; /* 利用负距离将其挤出原本的对齐线，保证头像是真正的对齐起点 */
  border: none;
  background: transparent;
  box-shadow: none;
}

.chat-header .el-avatar {
  /* 清除我刚才加的错误 margin-left */
  margin-left: 0;
}

.user-info h3 {
  margin: 0 0 5px 0;
  font-size: 16px;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px 30px;
  background-color: #f5f7fa;
}

/* 隐藏滚动条但保留滚动功能 (Webkit浏览器) */
.messages::-webkit-scrollbar {
  width: 0px;
  background: transparent;
}

.messages {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}

.message-item {
  display: flex;
  margin-bottom: 24px;
  align-items: flex-start;
  width: 100%;
}

.message-item.sent {
  justify-content: flex-end;
}

.message-item.received {
  justify-content: flex-start;
}

.message-avatar {
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  border: 1px solid #f0f2f5;
  margin-top: 2px; /* 使头像与气泡顶部更好地对齐 */
}

.message-item.received .message-avatar {
  margin-right: 12px;
}

.message-item.sent .my-avatar {
  margin-left: 12px;
}

.message-content {
  max-width: 65%;
  display: flex;
  flex-direction: column;
}

.message-item.sent .message-content {
  align-items: flex-end;
}

.message-item.received .message-content {
  align-items: flex-start;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 12px;
  word-break: break-word;
  font-size: 15px;
  line-height: 1.5;
  position: relative;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

/* 对方的气泡 */
.message-item.received .message-bubble {
  background-color: white;
  border: 1px solid #ebeef5;
  border-top-left-radius: 4px;
  color: #303133;
}

/* 我方的气泡 */
.message-item.sent .message-bubble {
  background-color: #409eff;
  color: white;
  border-top-right-radius: 4px;
}

.message-time {
  font-size: 12px;
  color: #b1b3b8;
  margin-top: 6px;
  padding: 0 4px;
}

.input-area {
  padding: 15px 30px; /* 修改为 30px，与上方内容区对齐 */
  background: white;
  border-top: 1px solid #ebeef5; /* 恢复边框 */
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex-shrink: 0;
  box-shadow: none; /* 去除阴影，让其不再显得过于悬浮 */
}

.input-area :deep(.el-textarea__inner) {
  border: none;
  box-shadow: none;
  background-color: transparent;
  padding: 5px 0;
  resize: none;
}

.input-area :deep(.el-textarea__inner:focus) {
  box-shadow: none;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
</style>
