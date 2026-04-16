<template>
  <div class="notifications">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat">消息</router-link>
            <router-link to="/notifications" class="active">通知</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <el-card class="notification-card">
          <template #header>
            <div class="card-header">
              <div class="left">
                <span class="title">通知中心</span>
                <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="badge">
                  <span>未读通知</span>
                </el-badge>
              </div>
              <div class="right">
                <el-button size="small" @click="handleMarkAllRead" :disabled="unreadCount === 0">
                  全部标记为已读
                </el-button>
                <el-button size="small" @click="handleClearRead" type="danger" plain>
                  清空已读
                </el-button>
              </div>
            </div>
          </template>

          <!-- 类型筛选 -->
          <div class="filter-bar">
            <el-radio-group v-model="filterType" @change="fetchNotifications">
              <el-radio-button :label="0">全部</el-radio-button>
              <el-radio-button :label="1">系统通知</el-radio-button>
              <el-radio-button :label="2">交换通知</el-radio-button>
              <el-radio-button :label="3">评论通知</el-radio-button>
            </el-radio-group>

            <el-radio-group v-model="filterRead" @change="fetchNotifications" class="read-filter">
              <el-radio-button label="">全部</el-radio-button>
              <el-radio-button label="0">未读</el-radio-button>
              <el-radio-button label="1">已读</el-radio-button>
            </el-radio-group>
          </div>

          <!-- 通知列表 -->
          <div v-loading="loading" class="notification-list">
            <el-empty v-if="notifications.length === 0 && !loading" description="暂无通知" />

            <div
              v-for="notif in notifications"
              :key="notif.id"
              :class="['notification-item', notif.is_read === 0 ? 'unread' : '']"
              @click="handleClick(notif)"
            >
              <div class="notif-icon">
                <el-icon v-if="notif.type === 1" class="system-icon"><Bell /></el-icon>
                <el-icon v-else-if="notif.type === 2" class="exchange-icon"><Switch /></el-icon>
                <el-icon v-else-if="notif.type === 3" class="comment-icon"><ChatDotSquare /></el-icon>
              </div>

              <div class="notif-content">
                <div class="notif-header">
                  <span class="notif-title">{{ notif.title }}</span>
                  <span class="notif-time">{{ formatTime(notif.create_time) }}</span>
                </div>
                <div class="notif-text">{{ notif.content }}</div>
              </div>

              <div class="notif-actions" @click.stop>
                <el-button
                  v-if="notif.is_read === 0"
                  type="primary"
                  link
                  size="small"
                  @click="handleMarkRead(notif.id)"
                >
                  标为已读
                </el-button>
                <el-button
                  type="danger"
                  link
                  size="small"
                  @click="handleDelete(notif.id)"
                >
                  删除
                </el-button>
              </div>
            </div>

            <!-- 分页 -->
            <el-pagination
              v-if="total > 0"
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              class="pagination"
              @size-change="fetchNotifications"
              @current-change="fetchNotifications"
            />
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import {
  getNotificationList,
  getUnreadCount,
  markAsRead,
  markAllAsRead,
  deleteNotification,
  clearReadNotifications
} from '@/api/notification'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Bell, Switch, ChatDotSquare } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const notifications = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const unreadCount = ref(0)
const filterType = ref(0)
const filterRead = ref('')

let refreshTimer = null

// 获取通知列表
const fetchNotifications = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    if (filterType.value > 0) {
      params.type = filterType.value
    }

    if (filterRead.value !== '') {
      params.is_read = parseInt(filterRead.value)
    }

    const res = await getNotificationList(params)
    notifications.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取通知列表失败:', error)
    ElMessage.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

// 获取未读数量
const fetchUnreadCount = async () => {
  try {
    const res = await getUnreadCount()
    unreadCount.value = res.data.count
  } catch (error) {
    console.error('获取未读数量失败:', error)
  }
}

// 点击通知
const handleClick = (notif) => {
  // 标记为已读
  if (notif.is_read === 0) {
    handleMarkRead(notif.id)
  }

  // 跳转链接
  if (notif.link) {
    router.push(notif.link)
  }
}

// 标记单条为已读
const handleMarkRead = async (notifId) => {
  try {
    await markAsRead(notif.id || notifId)
    ElMessage.success('标记成功')

    // 刷新列表和未读数
    await fetchNotifications()
    await fetchUnreadCount()
  } catch (error) {
    console.error('标记失败:', error)
    ElMessage.error('标记失败')
  }
}

// 全部标记为已读
const handleMarkAllRead = async () => {
  try {
    await markAllAsRead()
    ElMessage.success('全部标记为已读')

    // 刷新列表和未读数
    await fetchNotifications()
    await fetchUnreadCount()
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error('操作失败')
  }
}

// 删除通知
const handleDelete = async (notifId) => {
  try {
    await ElMessageBox.confirm('确定要删除这条通知吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteNotification(notifId)
    ElMessage.success('删除成功')

    // 刷新列表和未读数
    await fetchNotifications()
    await fetchUnreadCount()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 清空已读通知
const handleClearRead = async () => {
  try {
    await ElMessageBox.confirm('确定要清空所有已读通知吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await clearReadNotifications()
    ElMessage.success('清空成功')

    // 刷新列表
    await fetchNotifications()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清空失败:', error)
      ElMessage.error('清空失败')
    }
  }
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

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchNotifications()
  fetchUnreadCount()

  // 定时刷新未读数（每30秒）
  refreshTimer = setInterval(() => {
    fetchUnreadCount()
  }, 30000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.notifications {
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
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  width: 100%;
}

.notification-card {
  min-height: 600px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.badge span {
  font-size: 14px;
  color: #606266;
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.read-filter {
  margin-left: auto;
}

.notification-list {
  min-height: 400px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 15px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.unread {
  background-color: #ecf5ff;
}

.notif-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  flex-shrink: 0;
}

.system-icon {
  font-size: 20px;
  color: #409eff;
}

.exchange-icon {
  font-size: 20px;
  color: #67c23a;
}

.comment-icon {
  font-size: 20px;
  color: #e6a23c;
}

.notif-content {
  flex: 1;
  min-width: 0;
}

.notif-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 5px;
}

.notif-title {
  font-weight: 500;
  color: #303133;
}

.notif-time {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.notif-text {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.notif-actions {
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex-shrink: 0;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
