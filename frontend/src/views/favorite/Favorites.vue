<template>
  <div class="favorites">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat">消息</router-link>
            <router-link to="/publish">发布商品</router-link>
            <router-link to="/profile">个人中心</router-link>
            <router-link to="/favorites" class="active">我的收藏</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <el-card class="favorites-card">
          <template #header>
            <div class="card-header">
              <span>我的收藏</span>
              <span class="count">共 {{ total }} 件商品</span>
            </div>
          </template>

          <!-- 收藏列表 -->
          <div v-loading="loading" class="favorites-list">
            <el-empty v-if="list.length === 0 && !loading" description="暂无收藏商品" />

            <el-row :gutter="20">
              <el-col
                v-for="item in list"
                :key="item.id"
                :xs="24"
                :sm="12"
                :md="8"
                :lg="6"
                class="goods-item-col"
              >
                <el-card class="goods-card" shadow="hover">
                  <!-- 商品图片 -->
                  <div class="goods-images" @click="goDetail(item.goods.id)">
                    <el-image
                      v-if="item.goods && item.goods.getImages && item.goods.getImages().length > 0"
                      :src="item.goods.getImages()[0]"
                      fit="cover"
                      class="goods-image"
                      :preview-src-list="item.goods.getImages()"
                    >
                      <template #error>
                        <div class="image-error">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <div v-else class="no-image">
                      <el-icon><Picture /></el-icon>
                    </div>
                    <div v-if="item.goods?.type === 2" class="exchange-tag">交换</div>
                    <!-- 状态标签 -->
                    <div v-if="item.goods?.status !== 1" :class="['status-tag', getStatusClass(item.goods?.status)]">
                      {{ getStatusText(item.goods?.status) }}
                    </div>
                  </div>

                  <!-- 商品信息 -->
                  <div class="goods-info">
                    <div class="goods-title" @click="goDetail(item.goods.id)">
                      {{ item.goods?.title }}
                    </div>
                    <div class="goods-meta">
                      <span class="price" v-if="item.goods">
                        <template v-if="item.goods.type === 1 || item.goods.type === 3">
                          ¥{{ item.goods.price }}
                        </template>
                        <el-tag v-else type="success" size="small">交换</el-tag>
                      </span>
                      <span class="condition">
                        {{ getConditionText(item.goods?.condition) }}
                      </span>
                    </div>
                    <div class="goods-footer">
                      <div class="user-info">
                        <el-avatar :size="24" :src="item.goods?.user?.avatar">
                          {{ item.goods?.user?.username?.[0] }}
                        </el-avatar>
                        <span class="username">{{ item.goods?.user?.nickname || item.goods?.user?.username }}</span>
                      </div>
                      <el-button
                        type="danger"
                        size="small"
                        :icon="Delete"
                        @click.stop="handleRemove(item.goods.id)"
                      >
                        取消收藏
                      </el-button>
                    </div>
                    <div class="favorite-time">
                      收藏于 {{ formatTime(item.create_time) }}
                    </div>
                  </div>
                </el-card>
              </el-col>
            </el-row>
          </div>

          <!-- 分页 -->
          <div v-if="total > 0" class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[12, 24, 48]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="fetchFavorites"
              @current-change="fetchFavorites"
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
import { getFavoriteList, removeFavorite } from '@/api/favorite'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Picture } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const list = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const defaultImage = 'https://via.placeholder.com/300x200?text=No+Image'

// 获取收藏列表
const fetchFavorites = async () => {
  loading.value = true
  try {
    const res = await getFavoriteList({
      page: currentPage.value,
      page_size: pageSize.value
    })

    // 处理图片数据
    const processedList = (res.data.list || []).map(item => {
      return {
        ...item,
        goods: {
          ...item.goods,
          getImages: () => {
            if (!item.goods) return []
            if (typeof item.goods.images === 'string') {
              try {
                return JSON.parse(item.goods.images)
              } catch {
                return []
              }
            }
            return item.goods.images || []
          }
        }
      }
    })

    list.value = processedList
    total.value = res.data.total
  } catch (error) {
    console.error('获取收藏列表失败:', error)
    ElMessage.error('获取收藏列表失败')
  } finally {
    loading.value = false
  }
}

// 取消收藏
const handleRemove = async (goodsId) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await removeFavorite(goodsId)
    ElMessage.success('取消收藏成功')

    // 重新获取列表
    await fetchFavorites()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消收藏失败:', error)
      ElMessage.error('取消收藏失败')
    }
  }
}

// 查看商品详情
const goDetail = (goodsId) => {
  router.push(`/goods/${goodsId}`)
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
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
  fetchFavorites()
})
</script>

<style scoped>
.favorites {
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
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  width: 100%;
}

.favorites-card {
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

.favorites-list {
  margin-top: 20px;
}

.goods-item-col {
  margin-bottom: 20px;
}

.goods-card {
  height: 100%;
  transition: all 0.3s;
  overflow: hidden;
}

.goods-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.goods-images {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 4px;
  margin-bottom: 12px;
  cursor: pointer;
}

.goods-images:hover .goods-image {
  opacity: 0.8;
}

.goods-images {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 4px;
  margin-bottom: 12px;
}

.goods-image {
  width: 100%;
  height: 100%;
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

.image-error,
.no-image {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 48px;
}

.exchange-tag {
  position: absolute;
  top: 10px;
  right: 10px;
  background: #67c23a;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-tag {
  position: absolute;
  top: 10px;
  left: 10px;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
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

.goods-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.goods-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  min-height: 44px;
  cursor: pointer;
}

.goods-title:hover {
  color: #409eff;
}

.goods-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.price {
  font-size: 20px;
  font-weight: bold;
  color: #f56c6c;
}

.condition {
  font-size: 12px;
  color: #909399;
  background: #f0f0f0;
  padding: 4px 10px;
  border-radius: 12px;
}

.goods-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-size: 14px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
}

.favorite-time {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
