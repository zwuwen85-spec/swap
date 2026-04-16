<template>
  <div class="goods-list">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods" class="active">商品列表</router-link>
            <template v-if="userStore.isLogin">
              <router-link to="/exchange">我的交换</router-link>
              <router-link to="/my-goods">我的发布</router-link>
              <router-link to="/publish">发布商品</router-link>
              <router-link to="/profile">个人中心</router-link>
              <a @click="handleLogout">退出</a>
            </template>
            <template v-else>
              <router-link to="/login">登录</router-link>
              <router-link to="/register">注册</router-link>
            </template>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <!-- 筛选栏 -->
        <div class="filter-section">
          <el-card class="filter-card" shadow="never">
            <div class="filter-content">
              <div class="filter-header">
                <div class="filter-title">
                  <el-icon class="title-icon"><Filter /></el-icon>
                  <span>商品筛选</span>
                </div>
                <el-button
                  v-if="hasActiveFilters"
                  type="warning"
                  size="default"
                  plain
                  @click="resetFilters"
                >
                  <el-icon><RefreshLeft /></el-icon>
                  重置筛选
                </el-button>
              </div>

              <el-row :gutter="24" class="filter-row">
                <!-- 搜索框 -->
                <el-col :xs="24" :sm="24" :md="12" :lg="12">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><Search /></el-icon>
                      搜索商品
                    </label>
                    <el-input
                      v-model="searchKeyword"
                      placeholder="输入关键词搜索想要的商品..."
                      @keyup.enter="handleSearch"
                      clearable
                      size="large"
                      class="search-input"
                    >
                      <template #append>
                        <el-button @click="handleSearch" type="primary" size="large">
                          <el-icon><Search /></el-icon>
                          搜索
                        </el-button>
                      </template>
                    </el-input>
                  </div>
                </el-col>

                <!-- 分类筛选 -->
                <el-col :xs="12" :sm="12" :md="6" :lg="6">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><Grid /></el-icon>
                      商品分类
                    </label>
                    <el-select
                      v-model="filters.category_id"
                      placeholder="选择分类"
                      clearable
                      @change="handleFilterChange"
                      size="large"
                      class="filter-select"
                    >
                      <el-option label="全部分类" :value="0" />
                      <el-option
                        v-for="category in categories"
                        :key="category.id"
                        :label="category.name"
                        :value="category.id"
                      />
                    </el-select>
                  </div>
                </el-col>

                <!-- 类型筛选 -->
                <el-col :xs="12" :sm="12" :md="6" :lg="6">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><ShoppingCart /></el-icon>
                      交易类型
                    </label>
                    <el-select
                      v-model="filters.type"
                      placeholder="选择类型"
                      clearable
                      @change="handleFilterChange"
                      size="large"
                      class="filter-select"
                    >
                      <el-option label="全部类型" :value="0" />
                      <el-option label="售卖" :value="1" />
                      <el-option label="交换" :value="2" />
                      <el-option label="均可" :value="3" />
                    </el-select>
                  </div>
                </el-col>
              </el-row>

              <el-row :gutter="24" class="filter-row-secondary">
                <!-- 成色筛选 -->
                <el-col :xs="12" :sm="12" :md="6" :lg="6">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><Star /></el-icon>
                      商品成色
                    </label>
                    <el-select
                      v-model="filters.condition"
                      placeholder="选择成色"
                      clearable
                      @change="handleFilterChange"
                      size="large"
                      class="filter-select"
                    >
                      <el-option label="全部成色" :value="0" />
                      <el-option label="全新" :value="1" />
                      <el-option label="九成新" :value="2" />
                      <el-option label="八成新" :value="3" />
                      <el-option label="七成新" :value="4" />
                    </el-select>
                  </div>
                </el-col>

                <!-- 排序 -->
                <el-col :xs="12" :sm="12" :md="6" :lg="6">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><Sort /></el-icon>
                      排序方式
                    </label>
                    <el-select
                      v-model="filters.sort"
                      placeholder="选择排序"
                      @change="handleFilterChange"
                      size="large"
                      class="filter-select"
                    >
                      <el-option label="最新发布" value="time_desc" />
                      <el-option label="最早发布" value="time_asc" />
                      <el-option label="价格从低到高" value="price_asc" />
                      <el-option label="价格从高到低" value="price_desc" />
                    </el-select>
                  </div>
                </el-col>

                <!-- 快速筛选标签 -->
                <el-col :xs="24" :sm="24" :md="12" :lg="12">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><PriceTag /></el-icon>
                      快速筛选
                    </label>
                    <div class="quick-filters">
                      <el-tag
                        :type="filters.type === 1 ? 'danger' : 'info'"
                        effect="plain"
                        @click="filters.type = filters.type === 1 ? 0 : 1; handleFilterChange()"
                        class="quick-tag"
                      >
                        想要购买
                      </el-tag>
                      <el-tag
                        :type="filters.type === 2 ? 'success' : 'info'"
                        effect="plain"
                        @click="filters.type = filters.type === 2 ? 0 : 2; handleFilterChange()"
                        class="quick-tag"
                      >
                        想要交换
                      </el-tag>
                      <el-tag
                        :type="filters.condition === 1 ? 'warning' : 'info'"
                        effect="plain"
                        @click="filters.condition = filters.condition === 1 ? 0 : 1; handleFilterChange()"
                        class="quick-tag"
                      >
                        仅看全新
                      </el-tag>
                    </div>
                  </div>
                </el-col>
              </el-row>

              <!-- 活跃筛选标签 -->
              <div v-if="hasActiveFilters" class="active-tags">
                <span class="tags-label">
                  <el-icon><Filter /></el-icon>
                  当前筛选：
                </span>
                <el-tag
                  v-if="searchKeyword"
                  closable
                  @close="searchKeyword = ''; handleSearch()"
                  class="filter-tag"
                  type="primary"
                >
                  搜索: {{ searchKeyword }}
                </el-tag>
                <el-tag
                  v-if="filters.category_id"
                  closable
                  @close="filters.category_id = 0; handleFilterChange()"
                  class="filter-tag"
                  type="success"
                >
                  分类: {{ getCategoryName(filters.category_id) }}
                </el-tag>
                <el-tag
                  v-if="filters.type"
                  closable
                  @close="filters.type = 0; handleFilterChange()"
                  class="filter-tag"
                  type="warning"
                >
                  类型: {{ getTypeName(filters.type) }}
                </el-tag>
                <el-tag
                  v-if="filters.condition"
                  closable
                  @close="filters.condition = 0; handleFilterChange()"
                  class="filter-tag"
                  type="info"
                >
                  成色: {{ getConditionText(filters.condition) }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </div>

        <!-- 商品列表 -->
        <div class="goods-wrapper" v-loading="loading">
          <el-empty v-if="goodsList.length === 0 && !loading" description="暂无商品" />

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
                    <!-- 类型标签 -->
                    <div v-if="goods.type === 2" class="type-badge">交换</div>
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
                    <el-tag v-else type="success" size="small">交换</el-tag>
                  </div>

                  <!-- 元数据 -->
                  <div class="meta-row">
                    <span class="condition-tag">{{ getConditionText(goods.condition) }}</span>
                    <span class="divider">|</span>
                    <span>{{ goods.view_count }} 浏览</span>
                  </div>

                  <!-- 用户信息 -->
                  <div class="user-row">
                    <el-avatar :size="20" :src="goods.user?.avatar">
                      {{ goods.user?.username?.[0] }}
                    </el-avatar>
                    <span class="username">{{ goods.user?.nickname || goods.user?.username }}</span>
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
            :page-sizes="[20, 40, 60, 80]"
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getGoodsList, getCategories } from '@/api/goods'
import { ElMessage } from 'element-plus'
import { Picture, Search, Filter, RefreshLeft, Grid, ShoppingCart, Star, Sort, PriceTag } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 搜索关键词
const searchKeyword = ref('')

// 筛选条件
const filters = reactive({
  category_id: 0,
  type: 0,
  condition: 0,
  sort: 'time_desc'
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20
})

// 数据
const loading = ref(false)
const goodsList = ref([])
const total = ref(0)
const categories = ref([])

// 计算是否有活跃筛选
const hasActiveFilters = computed(() => {
  return searchKeyword.value ||
         filters.category_id !== 0 ||
         filters.type !== 0 ||
         filters.condition !== 0
})

// 获取商品列表
const fetchGoodsList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...filters,
      keyword: searchKeyword.value || undefined
    }

    const res = await getGoodsList(params)
    goodsList.value = res.data.list
    total.value = res.data.total

    // 处理图片数据
    goodsList.value.forEach(goods => {
      goods.getImages = () => {
        if (typeof goods.images === 'string') {
          try {
            return JSON.parse(goods.images)
          } catch {
            return []
          }
        }
        return goods.images || []
      }
    })
  } catch (error) {
    console.error('获取商品列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res.data
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 筛选改变
const handleFilterChange = () => {
  pagination.page = 1
  fetchGoodsList()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchGoodsList()
}

// 分页改变
const handlePageChange = (page) => {
  pagination.page = page
  fetchGoodsList()
  // 滚动到顶部
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleSizeChange = (size) => {
  pagination.page_size = size
  pagination.page = 1
  fetchGoodsList()
}

// 跳转到详情
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

// 获取分类名称
const getCategoryName = (categoryId) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category?.name || '未知'
}

// 获取类型名称
const getTypeName = (type) => {
  const map = {
    1: '售卖',
    2: '交换',
    3: '均可'
  }
  return map[type] || '未知'
}

// 重置筛选
const resetFilters = () => {
  searchKeyword.value = ''
  filters.category_id = 0
  filters.type = 0
  filters.condition = 0
  filters.sort = 'time_desc'
  handleFilterChange()
}

// 初始化
onMounted(() => {
  fetchCategories()
  fetchGoodsList()
})
</script>

<style scoped>
.goods-list {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  width: 100%;
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
  width: 80%;
  margin: 0 auto;
  padding: 0 20px;
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

.nav a:hover,
.nav a.active {
  opacity: 0.8;
}

.el-main {
  width: 80%;
  max-width: 80%;
  margin: 0 auto;
  padding: 30px 20px;
}

/* 筛选栏 */
.filter-section {
  margin-bottom: 30px;
}

.filter-card {
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(64, 158, 255, 0.2);
  background: linear-gradient(135deg, #ffffff 0%, #f8fbff 100%);
  transition: all 0.3s;
}

.filter-card:hover {
  box-shadow: 0 8px 32px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
  border-color: rgba(64, 158, 255, 0.3);
}

.filter-content {
  padding: 14px 24px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 10px;
  border-bottom: 2px solid rgba(64, 158, 255, 0.15);
}

.filter-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.title-icon {
  font-size: 24px;
  color: #409eff;
}

.filter-row {
  margin-bottom: 14px;
}

.filter-row-secondary {
  margin-bottom: 0;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 2px;
}

.filter-label .el-icon {
  font-size: 18px;
  color: #409eff;
}

.search-input {
  width: 100%;
}

.search-input .el-input__wrapper {
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;
  padding: 5px 12px;
  height: 40px;
  font-size: 15px;
}

.search-input .el-input__wrapper:hover,
.search-input .el-input__wrapper.is-focus {
  box-shadow: 0 4px 20px rgba(64, 158, 255, 0.25);
  border-color: #409eff;
}

.filter-select {
  width: 100%;
}

.filter-select .el-input__wrapper {
  border-radius: 10px;
  transition: all 0.3s;
  height: 40px;
  font-size: 15px;
}

.filter-select:hover .el-input__wrapper {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.15);
}

.quick-filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.quick-tag {
  cursor: pointer;
  transition: all 0.3s;
  user-select: none;
  padding: 6px 16px;
  font-size: 14px;
  height: 40px;
  display: inline-flex;
  align-items: center;
  border-radius: 20px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

.quick-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);
}

.quick-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.active-tags {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 14px;
  padding: 10px 16px;
  background: linear-gradient(135deg, #f0f7ff 0%, #e6f4ff 100%);
  border-radius: 12px;
  border: 1px dashed rgba(64, 158, 255, 0.3);
  flex-wrap: wrap;
}

.tags-label {
  font-size: 15px;
  color: #606266;
  font-weight: 500;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 6px;
}

.tags-label .el-icon {
  color: #409eff;
  font-size: 16px;
}

.filter-tag {
  font-weight: 500;
  padding: 6px 14px;
  font-size: 14px;
  border: none;
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

.filter-tag:hover {
  opacity: 0.85;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

/* 商品网格 */
.goods-wrapper {
  min-height: 400px;
  margin-top: 10px;
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
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
  background: #ffffff;
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
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
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
  line-height: 1.4;
}

.title:hover {
  color: #409eff;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
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

.user-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
}

.username {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100px;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 30px 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .goods-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
  }

  .filter-card {
    width: 100%;
  }

  .filter-content {
    padding: 12px 16px;
  }

  .filter-header {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
    margin-bottom: 16px;
  }

  .filter-title {
    font-size: 18px;
  }

  .filter-row,
  .filter-row-secondary {
    margin-bottom: 12px;
  }

  .filter-item {
    margin-bottom: 8px;
  }

  .active-tags {
    flex-direction: column;
    align-items: flex-start;
    padding: 10px;
    margin-top: 12px;
  }

  .quick-filters {
    gap: 8px;
  }

  .quick-tag {
    padding: 5px 12px;
    font-size: 13px;
    height: auto;
  }

  .image-section {
    height: 180px;
  }

  .content-section {
    padding: 12px;
  }

  .title {
    font-size: 15px;
  }

  .price {
    font-size: 22px;
  }
}
</style>
