<template>
  <div class="my-comments">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat">消息</router-link>
            <router-link to="/favorites">我的收藏</router-link>
            <router-link to="/my-comments" class="active">我的评论</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <el-card class="comments-card">
          <template #header>
            <div class="card-header">
              <span>我的评论</span>
              <span class="count">共 {{ total }} 条</span>
            </div>
          </template>

          <!-- 评论列表 -->
          <div v-loading="loading">
            <el-empty v-if="comments.length === 0 && !loading" description="暂无评论" />

            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="goods-info" @click="goGoodsDetail(comment.goods.id)">
                <div class="goods-image-wrapper">
                  <el-image
                    v-if="comment.goods && comment.goods.getImages && comment.goods.getImages().length > 0"
                    :src="comment.goods.getImages()[0]"
                    fit="cover"
                    class="goods-image"
                  >
                    <template #error>
                      <div class="image-slot">
                        <el-icon><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                  <div v-else class="no-image">
                    <el-icon><Picture /></el-icon>
                  </div>
                  <!-- 状态标签 -->
                  <div v-if="comment.goods?.status !== 1" :class="['status-tag', getStatusClass(comment.goods?.status)]">
                    {{ getStatusText(comment.goods?.status) }}
                  </div>
                </div>
                <div class="goods-detail">
                  <div class="goods-title">{{ comment.goods?.title }}</div>
                  <div class="goods-meta">
                    <span class="goods-price" v-if="comment.goods?.price">
                      ¥{{ comment.goods.price }}
                    </span>
                    <span class="condition">{{ getConditionText(comment.goods?.condition) }}</span>
                  </div>
                </div>
              </div>

              <div class="comment-content">
                <div class="comment-header">
                  <el-rate v-model="comment.rating" disabled size="small" />
                  <span class="time">{{ formatTime(comment.create_time) }}</span>
                </div>
                <div class="comment-text">{{ comment.content }}</div>
                <div class="comment-actions">
                  <el-button type="danger" link @click="handleDelete(comment.id)">
                    删除
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
              @size-change="fetchComments"
              @current-change="fetchComments"
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
import { getMyComments, deleteComment } from '@/api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Picture } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const comments = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 获取我的评论
const fetchComments = async () => {
  loading.value = true
  try {
    const res = await getMyComments({
      page: currentPage.value,
      page_size: pageSize.value
    })

    // 处理图片数据
    const processedComments = (res.data.list || []).map(comment => {
      return {
        ...comment,
        goods: comment.goods ? {
          ...comment.goods,
          getImages: () => {
            if (typeof comment.goods.images === 'string') {
              try {
                return JSON.parse(comment.goods.images)
              } catch {
                return []
              }
            }
            return comment.goods.images || []
          }
        } : null
      }
    })

    comments.value = processedComments
    total.value = res.data.total
  } catch (error) {
    console.error('获取评论列表失败:', error)
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
  }
}

// 删除评论
const handleDelete = async (commentId) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteComment(commentId)
    ElMessage.success('删除成功')

    // 刷新列表
    await fetchComments()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除评论失败:', error)
      ElMessage.error('删除评论失败')
    }
  }
}

// 查看商品详情
const goGoodsDetail = (goodsId) => {
  router.push(`/goods/${goodsId}`)
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    0: '已下架',
    1: '在售',
    2: '已售出',
    3: '已交换'
  }
  return map[status] || '未知'
}

// 获取状态样式类
const getStatusClass = (status) => {
  const map = {
    0: 'status-offline',
    1: 'status-online',
    2: 'status-sold',
    3: 'status-exchanged'
  }
  return map[status] || ''
}

// 获取成色文本
const getConditionText = (condition) => {
  const map = {
    1: '全新',
    2: '9成新',
    3: '8成新',
    4: '7成新'
  }
  return map[condition] || '未知'
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  fetchComments()
})
</script>

<style scoped>
.my-comments {
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

.comments-card {
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

.comment-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.goods-info {
  display: flex;
  gap: 15px;
  width: 300px;
  cursor: pointer;
}

.goods-image-wrapper {
  position: relative;
  width: 80px;
  height: 80px;
  flex-shrink: 0;
}

.goods-image {
  width: 80px;
  height: 80px;
  border-radius: 4px;
}

.status-tag {
  position: absolute;
  top: 4px;
  left: 4px;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.status-offline {
  background: #909399;
}

.status-online {
  background: #67c23a;
}

.status-sold {
  background: #f56c6c;
}

.status-exchanged {
  background: #e6a23c;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 30px;
}

.no-image {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 80px;
  height: 80px;
  border-radius: 4px;
  flex-shrink: 0;
  background: #f5f7fa;
  color: #909399;
  font-size: 30px;
}

.goods-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.goods-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  margin-bottom: 8px;
}

.goods-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.goods-price {
  font-size: 16px;
  font-weight: bold;
  color: #f56c6c;
}

.condition {
  font-size: 11px;
  color: #909399;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 10px;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.time {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.comment-text {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 10px;
}

.comment-actions {
  display: flex;
  gap: 10px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
