<template>
  <div class="my-reports">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat">消息</router-link>
            <router-link to="/notifications">通知</router-link>
            <router-link to="/my-reports" class="active">我的举报</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <el-card class="reports-card">
          <template #header>
            <div class="card-header">
              <span>我的举报</span>
              <span class="count">共 {{ total }} 条</span>
            </div>
          </template>

          <!-- 举报列表 -->
          <div v-loading="loading">
            <el-empty v-if="reports.length === 0 && !loading" description="暂无举报记录" />

            <div v-for="report in reports" :key="report.id" class="report-item">
              <!-- 目标对象信息 -->
              <div class="target-info">
                <div class="target-header">
                  <el-tag :type="getTargetTypeTag(report.target_type)" size="small">
                    {{ getTargetTypeText(report.target_type) }}
                  </el-tag>
                  <span class="target-name">{{ getTargetName(report) }}</span>
                </div>

                <!-- 如果是商品，显示商品信息 -->
                <div v-if="report.target_type === 1 && report.goods" class="goods-info">
                  <el-image
                    :src="report.goods.images?.[0]"
                    fit="cover"
                    class="goods-image"
                  >
                    <template #error>
                      <div class="image-slot">
                        <el-icon><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                  <span class="goods-title">{{ report.goods.title }}</span>
                </div>
              </div>

              <!-- 举报信息 -->
              <div class="report-info">
                <div class="report-header">
                  <span class="reason">举报原因：{{ report.reason }}</span>
                  <el-tag :type="getStatusType(report.status)" size="small">
                    {{ getStatusText(report.status) }}
                  </el-tag>
                </div>

                <div v-if="report.description" class="description">
                  {{ report.description }}
                </div>

                <div class="report-meta">
                  <span class="time">举报时间：{{ formatTime(report.create_time) }}</span>

                  <!-- 处理结果 -->
                  <div v-if="report.status !== 0" class="handle-result">
                    <span v-if="report.status === 1" class="success">
                      <el-icon><Check /></el-icon>
                      举报成立
                    </span>
                    <span v-else-if="report.status === 2" class="rejected">
                      <el-icon><Close /></el-icon>
                      举报驳回
                      <span v-if="report.handle_result">：{{ report.handle_result }}</span>
                    </span>
                  </div>
                </div>

                <!-- 操作按钮 -->
                <div class="actions">
                  <el-button
                    v-if="report.status === 0"
                    type="danger"
                    link
                    size="small"
                    @click="handleCancel(report.id)"
                  >
                    撤销举报
                  </el-button>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <el-pagination
              v-if="total > 0"
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 50]"
              layout="total, sizes, prev, pager, next"
              class="pagination"
              @size-change="fetchReports"
              @current-change="fetchReports"
            />
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
import { getMyReports, cancelReport } from '@/api/report'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Picture, Check, Close } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const reports = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 获取举报列表
const fetchReports = async () => {
  loading.value = true
  try {
    const res = await getMyReports({
      page: currentPage.value,
      page_size: pageSize.value
    })
    reports.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取举报列表失败:', error)
    ElMessage.error('获取举报列表失败')
  } finally {
    loading.value = false
  }
}

// 撤销举报
const handleCancel = async (reportId) => {
  try {
    await ElMessageBox.confirm('确定要撤销这条举报吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await cancelReport(reportId)
    ElMessage.success('撤销成功')

    // 刷新列表
    await fetchReports()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('撤销失败:', error)
      ElMessage.error('撤销失败')
    }
  }
}

// 获取目标类型文本
const getTargetTypeText = (type) => {
  const map = {
    1: '商品',
    2: '用户',
    3: '评论'
  }
  return map[type] || '未知'
}

// 获取目标类型标签类型
const getTargetTypeTag = (type) => {
  const map = {
    1: 'success',
    2: 'warning',
    3: 'info'
  }
  return map[type] || ''
}

// 获取目标名称
const getTargetName = (report) => {
  if (report.target_type === 1 && report.goods) {
    return report.goods.title
  } else if (report.target_type === 2 && report.target_user) {
    return report.target_user.nickname || report.target_user.username
  } else if (report.target_type === 3 && report.comment) {
    return report.comment.content.substring(0, 50)
  }
  return ''
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    0: '待处理',
    1: '已处理',
    2: '已驳回'
  }
  return map[status] || '未知'
}

// 获取状态标签类型
const getStatusType = (status) => {
  const map = {
    0: 'warning',
    1: 'success',
    2: 'info'
  }
  return map[status] || ''
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchReports()
})
</script>

<style scoped>
.my-reports {
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

.reports-card {
  min-height: 600px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.count {
  color: #909399;
  font-size: 14px;
}

.report-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
}

.report-item:last-child {
  border-bottom: none;
}

.target-info {
  width: 300px;
  flex-shrink: 0;
}

.target-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.target-name {
  font-weight: 500;
  color: #303133;
}

.goods-info {
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.goods-image {
  width: 50px;
  height: 50px;
  border-radius: 4px;
  flex-shrink: 0;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #ebeef5;
  color: #909399;
  font-size: 20px;
}

.goods-title {
  flex: 1;
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.report-info {
  flex: 1;
}

.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.reason {
  font-weight: 500;
  color: #303133;
}

.description {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 10px;
}

.report-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.time {
  font-size: 12px;
  color: #909399;
}

.handle-result {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
}

.success {
  color: #67c23a;
}

.rejected {
  color: #909399;
}

.actions {
  display: flex;
  gap: 10px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
