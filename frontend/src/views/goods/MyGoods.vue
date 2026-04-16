<template>
  <div class="my-goods">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/my-goods" class="active">我的发布</router-link>
            <router-link to="/publish">发布商品</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <!-- 页面标题栏 -->
        <div class="page-header">
          <div class="header-left">
            <h2>我的发布</h2>
            <p class="subtitle">管理你发布的所有商品</p>
          </div>
          <el-button type="primary" size="large" @click="goToPublish">
            <el-icon><Plus /></el-icon>
            发布新商品
          </el-button>
        </div>

        <!-- 统计卡片 -->
        <div class="stats-cards">
          <div
            v-for="(item, index) in statusStats"
            :key="index"
            class="stat-card"
            :class="{ active: statusFilter === item.value }"
            @click="handleStatusClick(item.value)"
          >
            <div class="stat-icon" :style="{ backgroundColor: item.color }">
              <el-icon :size="24">
                <component :is="item.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ item.count }}</div>
              <div class="stat-label">{{ item.label }}</div>
            </div>
          </div>
          <div class="stat-card total" @click="handleStatusClick(-1)">
            <div class="stat-icon total">
              <el-icon :size="24"><Grid /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ allGoodsTotal }}</div>
              <div class="stat-label">全部商品</div>
            </div>
          </div>
        </div>

        <!-- 商品网格 -->
        <div v-loading="loading" class="goods-wrapper">
          <el-empty v-if="goodsList.length === 0 && !loading" description="暂无商品">
            <el-button type="primary" @click="goToPublish">发布第一个商品</el-button>
          </el-empty>

          <div v-else class="goods-grid">
            <div
              v-for="goods in goodsList"
              :key="goods.id"
              class="goods-item"
            >
              <!-- 商品卡片 -->
              <el-card class="goods-card" shadow="hover">
                <!-- 图片区域 -->
                <div class="image-section" @click="goDetail(goods.id)">
                  <div v-if="goods.getImages && goods.getImages().length > 0" class="image-wrapper">
                    <el-image
                      :src="goods.getImages()[0]"
                      fit="cover"
                      class="goods-image"
                    >
                      <template #error>
                        <div class="image-error">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <!-- 状态标签 -->
                    <div class="status-badge" :class="getStatusClass(goods.status)">
                      {{ getStatusText(goods.status) }}
                    </div>
                    <!-- 类型标签 -->
                    <div v-if="goods.type === 1" class="type-badge sale">售卖</div>
                    <div v-else-if="goods.type === 2" class="type-badge">交换</div>
                    <div v-else-if="goods.type === 3" class="type-badge both">均可</div>
                  </div>
                  <div v-else class="image-wrapper no-image" @click="goDetail(goods.id)">
                    <el-icon><Picture /></el-icon>
                  </div>
                </div>

                <!-- 内容区域 -->
                <div class="content-section">
                  <!-- 标题 -->
                  <h3 class="title" @click="goDetail(goods.id)">{{ goods.title }}</h3>

                  <!-- 价格 -->
                  <div class="price-row">
                    <template v-if="goods.type === 1 || goods.type === 3">
                      <span class="price">¥{{ goods.price }}</span>
                      <span v-if="goods.original_price" class="original-price">¥{{ goods.original_price }}</span>
                    </template>
                    <el-tag v-else type="success" size="large">交换</el-tag>
                  </div>

                  <!-- 元数据 -->
                  <div class="meta-row">
                    <span class="condition-tag">{{ getConditionText(goods.condition) }}</span>
                    <span class="divider">|</span>
                    <span>{{ goods.view_count }} 浏览</span>
                    <span class="divider">|</span>
                    <span>{{ goods.favorite_count || 0 }} 收藏</span>
                  </div>

                  <!-- 操作栏 -->
                  <div class="action-bar">
                    <el-dropdown trigger="click" @command="(cmd) => handleCommand({ action: cmd, goods })">
                      <el-button size="small" type="primary" plain>
                        操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item command="edit">
                            <el-icon><Edit /></el-icon>
                            <span>编辑商品</span>
                          </el-dropdown-item>
                          <el-dropdown-item command="share">
                            <el-icon><Share /></el-icon>
                            <span>分享商品</span>
                          </el-dropdown-item>
                          <el-dropdown-item divided v-if="goods.status === 1" command="markSold">
                            <el-icon><CircleCheck /></el-icon>
                            <span>标记为已售出</span>
                          </el-dropdown-item>
                          <el-dropdown-item v-if="goods.status === 1" command="offSale">
                            <el-icon><Download /></el-icon>
                            <span>下架商品</span>
                          </el-dropdown-item>
                          <el-dropdown-item v-if="goods.status === 0 || goods.status === 2" command="onSale">
                            <el-icon><Upload /></el-icon>
                            <span>重新上架</span>
                          </el-dropdown-item>
                          <el-dropdown-item divided command="delete" class="danger">
                            <el-icon><Delete /></el-icon>
                            <span>删除商品</span>
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </div>
              </el-card>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="total > 0" class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.page_size"
            :total="total"
            :page-sizes="[12, 24, 36, 48]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getMyGoods, deleteGoods, updateGoods } from '@/api/goods'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Picture,
  Edit,
  Download,
  Upload,
  Delete,
  ArrowDown,
  Grid,
  CircleCheck,
  CircleClose,
  Warning,
  Share
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const goodsList = ref([])
const total = ref(0)  // 当前筛选条件的总数
const allGoodsTotal = ref(0)  // 所有商品的总数
const statusFilter = ref(-1)

const pagination = reactive({
  page: 1,
  page_size: 12
})

// 状态统计
const statusStats = ref([
  { label: '已上架', value: 1, count: 0, color: '#67c23a', icon: 'CircleCheck' },
  { label: '已售出', value: 2, count: 0, color: '#909399', icon: 'CircleClose' },
  { label: '已下架', value: 0, count: 0, color: '#e6a23c', icon: 'Warning' },
  { label: '已交换', value: 3, count: 0, color: '#409eff', icon: 'Grid' }
])

// 获取我的商品列表
const fetchGoodsList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size
    }

    if (statusFilter.value >= 0) {
      params.status = statusFilter.value
    }

    const res = await getMyGoods(params)
    const list = res.data.list || []

    // 处理图片数据
    goodsList.value = list.map(goods => {
      return {
        ...goods,
        getImages: () => {
          if (typeof goods.images === 'string') {
            try {
              return JSON.parse(goods.images)
            } catch {
              return []
            }
          }
          return goods.images || []
        }
      }
    })

    total.value = res.data.total
  } catch (error) {
    console.error('获取商品列表失败:', error)
    ElMessage.error('获取商品列表失败')
  } finally {
    loading.value = false
  }
}

// 更新状态统计
const updateStatusStats = async () => {
  try {
    const statuses = [1, 2, 0, 3]  // 已上架、已售、已下架、已交换
    let totalCount = 0

    for (let i = 0; i < statuses.length; i++) {
      const res = await getMyGoods({
        page: 1,
        page_size: 1,
        status: statuses[i]
      })
      const count = res.data.total
      statusStats.value[i].count = count
      totalCount += count
    }

    // 计算所有商品总数
    allGoodsTotal.value = totalCount
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 状态改变
const handleStatusClick = (status) => {
  statusFilter.value = status
  pagination.page = 1
  fetchGoodsList()
}

// 分页改变
const handlePageChange = (page) => {
  pagination.page = page
  fetchGoodsList()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleSizeChange = (size) => {
  pagination.page_size = size
  pagination.page = 1
  fetchGoodsList()
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    0: '已下架',
    1: '已上架',
    2: '已售出',
    3: '已交换'
  }
  return map[status] || '未知'
}

// 获取状态样式类
const getStatusClass = (status) => {
  const map = {
    0: 'offsale',
    1: 'selling',
    2: 'sold',
    3: 'exchanged'
  }
  return map[status] || ''
}

// 获取成色文本
const getConditionText = (condition) => {
  const map = {
    1: '全新',
    2: '九成新',
    3: '八成新',
    4: '七成新'
  }
  return map[condition] || '未知'
}

// 操作菜单处理
const handleCommand = async ({ action, goods }) => {
  switch (action) {
    case 'edit':
      router.push(`/publish?edit=${goods.id}`)
      break

    case 'share':
      // 复制链接到剪贴板
      const url = `${window.location.origin}/goods/${goods.id}`
      navigator.clipboard.writeText(url).then(() => {
        ElMessage.success('链接已复制到剪贴板')
      }).catch(() => {
        ElMessage.error('复制失败，请手动复制链接')
      })
      break

    case 'markSold':
      try {
        await ElMessageBox.confirm('确定要将该商品标记为已售出吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'success'
        })

        await updateGoods({
          id: goods.id,
          status: 2  // 已售出状态
        })

        ElMessage.success('商品已标记为已售出')
        fetchGoodsList()
        await updateStatusStats()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('标记失败:', error)
          ElMessage.error('标记失败')
        }
      }
      break

    case 'offSale':
      try {
        await ElMessageBox.confirm('确定要下架该商品吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        await updateGoods({
          id: goods.id,
          status: 0  // 下架状态
        })

        ElMessage.success('商品已下架')
        fetchGoodsList()
        await updateStatusStats()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('下架失败:', error)
          ElMessage.error('下架失败')
        }
      }
      break

    case 'onSale':
      try {
        await updateGoods({
          id: goods.id,
          status: 1
        })

        ElMessage.success('商品已上架')
        fetchGoodsList()
        await updateStatusStats()
      } catch (error) {
        console.error('上架失败:', error)
        ElMessage.error('上架失败')
      }
      break

    case 'delete':
      try {
        await ElMessageBox.confirm('确定要删除该商品吗？删除后无法恢复！', '警告', {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'error'
        })

        await deleteGoods(goods.id)

        ElMessage.success('商品已删除')
        fetchGoodsList()
        await updateStatusStats()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          ElMessage.error('删除失败')
        }
      }
      break
  }
}

// 跳转到发布页面
const goToPublish = () => {
  router.push('/publish')
}

// 跳转到商品详情
const goDetail = (id) => {
  router.push(`/goods/${id}`)
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(async () => {
  // 先获取统计数据（包括所有商品总数）
  await updateStatusStats()
  // 然后获取商品列表
  await fetchGoodsList()
})
</script>

<style scoped>
.my-goods {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
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
  transition: all 0.3s;
}

.nav a:hover {
  opacity: 0.8;
}

.nav a.active {
  font-weight: bold;
}

.el-main {
  max-width: 1400px;
  margin: 0 auto;
  padding: 30px 20px;
}

/* 页面标题栏 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 20px 30px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.header-left h2 {
  margin: 0 0 5px 0;
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.subtitle {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 15px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 15px;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.stat-card.active {
  border: 2px solid #409eff;
  background: #ecf5ff;
}

.stat-card.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-icon.total {
  background: rgba(255, 255, 255, 0.2);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-card.total .stat-value {
  color: white;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
}

.stat-card.total .stat-label {
  color: rgba(255, 255, 255, 0.9);
}

/* 商品网格 */
.goods-wrapper {
  min-height: 400px;
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.goods-item {
  height: 100%;
}

.goods-card {
  height: 100%;
  border: none;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
}

.goods-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

/* 图片区域 */
.image-section {
  position: relative;
  width: 100%;
  height: 220px;
  overflow: hidden;
  cursor: pointer;
}

.image-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
}

.goods-image {
  width: 100%;
  height: 100%;
  transition: transform 0.3s;
}

.image-section:hover .goods-image {
  transform: scale(1.05);
}

.no-image {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  color: #909399;
  font-size: 48px;
}

.image-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  color: #909399;
  font-size: 48px;
}

/* 状态标签 */
.status-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  color: white;
  backdrop-filter: blur(4px);
}

.status-badge.selling {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
}

.status-badge.sold {
  background: linear-gradient(135deg, #909399 0%, #b3b3b3 100%);
}

.status-badge.offsale {
  background: linear-gradient(135deg, #e6a23c 0%, #f0a165 100%);
}

.status-badge.exchanged {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
}

/* 类型标签 */
.type-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  color: white;
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  backdrop-filter: blur(4px);
}

.type-badge.sale {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
}

.type-badge.both {
  background: linear-gradient(135deg, #e6a23c 0%, #f0a165 100%);
}

/* 内容区域 */
.content-section {
  padding: 16px;
}

.title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 12px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  color: #303133;
  transition: color 0.3s;
}

.title:hover {
  color: #409eff;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 12px;
}

.price {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
}

.original-price {
  font-size: 14px;
  color: #909399;
  text-decoration: line-through;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #909399;
  margin-bottom: 12px;
}

.condition-tag {
  padding: 3px 10px;
  background: linear-gradient(135deg, #f0f0f0 0%, #e8e8e8 100%);
  border-radius: 12px;
  font-weight: 500;
}

.divider {
  color: #dcdfe6;
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: flex-end;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 30px 0;
}

/* 下拉菜单样式 */
.danger {
  color: #f56c6c;
}

:deep(.el-dropdown-menu__item.danger:hover) {
  background-color: #fef0f0;
  color: #f56c6c;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
}

:deep(.el-dropdown-menu__item .el-icon) {
  font-size: 16px;
}

/* 响应式 */
@media (max-width: 1200px) {
  .stats-cards {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .page-header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }

  .goods-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  }
}
</style>
