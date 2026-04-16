<template>
  <div class="goods-detail">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">发现</router-link>
            <router-link v-if="userStore.isLogin" to="/profile">个人中心</router-link>
            <router-link v-else to="/login">登录</router-link>
          </div>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main v-loading="loading" class="detail-main">
        <div v-if="goods" class="detail-container">
          <!-- 主卡片 -->
          <el-card class="main-card" shadow="never">
            <!-- 图片展示 -->
            <div class="image-section">
              <div class="image-gallery">
                <div class="main-image-container" v-if="imageList.length > 0">
                  <el-image
                    :src="imageList[currentImageIndex]"
                    fit="contain"
                    class="main-image"
                    :preview-src-list="imageList"
                    :initial-index="currentImageIndex"
                  />
                  <div class="image-actions">
                    <el-button
                      circle
                      @click="previousImage"
                      :disabled="currentImageIndex === 0"
                      class="image-nav-btn"
                    >
                      <el-icon><ArrowLeft /></el-icon>
                    </el-button>
                    <el-button
                      circle
                      @click="nextImage"
                      :disabled="currentImageIndex === imageList.length - 1"
                      class="image-nav-btn"
                    >
                      <el-icon><ArrowRight /></el-icon>
                    </el-button>
                  </div>
                </div>
                <div v-else class="no-image">
                  <div class="no-image-content">
                    <el-icon :size="100"><Picture /></el-icon>
                    <p>暂无图片</p>
                  </div>
                </div>

                <!-- 缩略图 -->
                <div class="thumbnail-container" v-if="imageList.length > 1">
                  <div
                    v-for="(image, index) in imageList"
                    :key="index"
                    class="thumbnail"
                    :class="{ active: index === currentImageIndex }"
                    @click="currentImageIndex = index"
                  >
                    <el-image :src="image" fit="cover" />
                  </div>
                </div>
              </div>
            </div>

            <!-- 商品信息 -->
            <div class="info-section">
              <!-- 标题和类型标签 -->
              <div class="header-section">
                <h1 class="title">{{ goods.title }}</h1>
                <div class="tags">
                  <el-tag v-if="goods.type === 1" type="danger" effect="dark" size="small">出售</el-tag>
                  <el-tag v-if="goods.type === 2" type="success" effect="dark" size="small">交换</el-tag>
                  <el-tag v-if="goods.type === 3" type="warning" effect="dark" size="small">可售可换</el-tag>
                  <el-tag type="info" effect="plain" size="small">{{ getConditionText(goods.condition) }}</el-tag>
                </div>
              </div>

              <!-- 价格卡片 -->
              <div class="price-row">
                <div class="price-card">
                  <div v-if="goods.type === 1 || goods.type === 3" class="price-section">
                    <div class="price-value">
                      <span class="currency">¥</span>
                      <span class="amount">{{ goods.price }}</span>
                    </div>
                    <div v-if="goods.original_price" class="original-price">
                      原价 ¥{{ goods.original_price }}
                      <span class="discount">
                        {{ Math.round((1 - goods.price / goods.original_price) * 100) }}% OFF
                      </span>
                    </div>
                  </div>
                  <div v-else class="exchange-only">
                    <el-icon><Switch /></el-icon>
                    <span>仅支持交换</span>
                  </div>
                </div>
              </div>

              <!-- 支持交换提示条 -->
              <div v-if="goods.type === 2 || goods.type === 3" class="exchange-tip-bar">
                <div class="tip-content">
                  <el-icon class="tip-icon"><Switch /></el-icon>
                  <div class="tip-text-group">
                    <span class="tip-title">本商品支持交换</span>
                    <span class="tip-desc">您可以使用自己的商品与卖家进行交换</span>
                  </div>
                  <el-tag type="success" effect="light" size="small">交换模式</el-tag>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="action-buttons-row">
                <el-button
                  type="primary"
                  @click="handleWant"
                  class="main-action-btn"
                >
                  <el-icon><ChatDotSquare /></el-icon>
                  我想要
                </el-button>
                <el-button
                  v-if="goods.type === 2 || goods.type === 3"
                  @click="handleExchange"
                  class="main-action-btn exchange-btn"
                >
                  <el-icon><Switch /></el-icon>
                  发起交换
                </el-button>
                <el-button
                  @click="handleFavorite"
                  :class="['main-action-btn', 'favorite-btn', { active: isFavorited }]"
                >
                  <el-icon><Star /></el-icon>
                  {{ isFavorited ? '已收藏' : '收藏' }}
                </el-button>
                <el-button
                  @click="showReportDialog = true"
                  class="action-icon-btn"
                >
                  <el-icon><Warning /></el-icon>
                  举报
                </el-button>
              </div>

              <!-- 商品属性和卖家信息 -->
              <el-row :gutter="20" class="details-row">
                <!-- 商品属性 -->
                <el-col :span="12">
                  <div class="details-card">
                    <div class="details-card-header">
                      <el-icon><BoxIcon /></el-icon>
                      <span>商品信息</span>
                    </div>
                    <div class="attribute-list">
                      <div class="attribute-item">
                        <span class="attr-label">分类</span>
                        <span class="attr-value">{{ goods.category?.name || '未分类' }}</span>
                      </div>
                      <div class="attribute-item">
                        <span class="attr-label">发布时间</span>
                        <span class="attr-value">{{ formatTime(goods.create_time) }}</span>
                      </div>
                      <div class="attribute-item">
                        <span class="attr-label">浏览次数</span>
                        <span class="attr-value">{{ goods.view_count }}</span>
                      </div>
                      <div class="attribute-item">
                        <span class="attr-label">收藏次数</span>
                        <span class="attr-value">{{ goods.favorite_count }}</span>
                      </div>
                    </div>
                  </div>
                </el-col>

                <!-- 发布者信息 -->
                <el-col :span="12">
                  <div class="details-card seller-card">
                    <div class="details-card-header">
                      <el-icon><User /></el-icon>
                      <span>发布者</span>
                      <span class="seller-rating">
                        <el-icon><StarFilled /></el-icon>
                        {{ sellerCredit }}分
                      </span>
                    </div>
                    <div
                      class="seller-info"
                      @click="goToUserDetail"
                      :class="{ 'clickable': hasValidSeller }"
                      :title="hasValidSeller ? '点击查看发布者详情' : '用户信息不可用'"
                    >
                      <el-avatar :size="48" :src="goods.user?.avatar" class="seller-avatar">
                        {{ goods.user?.username?.[0] || '?' }}
                      </el-avatar>
                      <div class="seller-detail">
                        <div class="name">
                          {{ goods.user?.nickname || goods.user?.username || '未知用户' }}
                        </div>
                        <div class="seller-meta">
                          <span v-if="goods.user?.school" class="meta-text">
                            <el-icon><School /></el-icon>
                            {{ goods.user.school }}
                          </span>
                          <span v-if="goods.user?.student_id" class="meta-text">
                            <el-icon><Postcard /></el-icon>
                            {{ goods.user.student_id }}
                          </span>
                          <span v-if="goods.user?.gender" class="meta-text">
                            <el-icon><User /></el-icon>
                            {{ getGenderText(goods.user.gender) }}
                          </span>
                        </div>
                        <el-rate
                          v-model="sellerCredit"
                          disabled
                          show-score
                          text-color="#ff9900"
                          score-template="{value}"
                          size="small"
                        />
                      </div>
                      <el-icon v-if="hasValidSeller" class="arrow-icon"><ArrowRight /></el-icon>
                    </div>
                  </div>
                </el-col>
              </el-row>

              <!-- 商品描述 -->
              <div class="description-section" v-if="goods.description">
                <div class="section-header">
                  <el-icon><Document /></el-icon>
                  <span>商品描述</span>
                </div>
                <div class="description-text">
                  {{ goods.description }}
                </div>
              </div>

              <!-- 交易地点 -->
              <div class="location-section" v-if="goods.location">
                <div class="section-header">
                  <el-icon><Location /></el-icon>
                  <span>交易地点</span>
                </div>
                <div class="location-content">
                  <el-icon class="location-icon"><Location /></el-icon>
                  <span class="location-text">{{ goods.location }}</span>
                </div>
              </div>
            </div>
          </el-card>

          <!-- 评论列表 -->
          <CommentList v-if="goods" :goods-id="goods.id" class="comment-section" />
        </div>
      </el-main>
    </el-container>

    <!-- 举报对话框 -->
    <ReportDialog
      v-if="goods"
      v-model="showReportDialog"
      :target-type="1"
      :target-id="goods.id"
      :target-name="goods.title"
      @success="handleReportSuccess"
    />

    <!-- 我想要对话框 -->
    <el-dialog
      v-model="showWantDialog"
      title="购买商品"
      width="500px"
      class="custom-dialog"
    >
      <el-form :model="wantForm" label-width="80px">
        <el-form-item label="商品">
          <div class="dialog-goods-info">
            <div class="dialog-goods-title">{{ goods?.title }}</div>
            <div v-if="goods?.type === 1 || goods?.type === 3" class="dialog-goods-price">
              ¥{{ goods?.price }}
            </div>
          </div>
        </el-form-item>
        <el-form-item label="留言">
          <el-input
            v-model="wantForm.message"
            type="textarea"
            :rows="4"
            placeholder="给卖家留个言吧，说明你的购买意向..."
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showWantDialog = false">取消</el-button>
          <el-button type="primary" @click="submitWant" :loading="submitting">
            确认购买
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 发起交换对话框 -->
    <el-dialog
      v-model="showExchangeDialog"
      title="发起交换"
      width="600px"
      class="custom-dialog"
    >
      <el-form :model="exchangeForm" label-width="100px">
        <el-form-item label="目标商品">
          <div class="exchange-goods-card">
            <el-image
              v-if="imageList.length > 0"
              :src="imageList[0]"
              fit="cover"
              class="exchange-goods-image"
            />
            <div class="exchange-goods-info">
              <div class="exchange-goods-title">{{ goods?.title }}</div>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="我的商品">
          <el-select
            v-model="exchangeForm.my_goods_id"
            placeholder="选择要交换的商品"
            style="width: 100%"
            size="large"
          >
            <el-option
              v-for="item in myGoodsList"
              :key="item.id"
              :label="item.title"
              :value="item.id"
            >
              <div class="select-option">
                <el-image
                  v-if="item.images && item.getImages && item.getImages().length > 0"
                  :src="item.getImages()[0]"
                  fit="cover"
                  class="select-option-image"
                />
                <span class="select-option-title">{{ item.title }}</span>
              </div>
            </el-option>
          </el-select>
          <div class="tip">如果没有可交换的商品，请先发布商品</div>
        </el-form-item>
        <el-form-item label="留言">
          <el-input
            v-model="exchangeForm.message"
            type="textarea"
            :rows="4"
            placeholder="介绍一下你的商品，说明交换意向..."
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showExchangeDialog = false" size="large">取消</el-button>
          <el-button type="primary" @click="submitExchange" :loading="submitting" size="large">
            发起交换
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 未找到 -->
    <el-empty v-if="!goods && !loading" description="商品不存在" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getGoodsDetail } from '@/api/goods'
import { addFavorite, removeFavorite, checkFavorite } from '@/api/favorite'
import { createExchange } from '@/api/exchange'
import { getMyGoods } from '@/api/goods'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  ArrowRight,
  Picture,
  ChatDotSquare,
  Switch,
  Star,
  Location,
  Warning,
  Document,
  User,
  Box as BoxIcon,
  Clock,
  View,
  StarFilled,
  School,
  Postcard
} from '@element-plus/icons-vue'
import CommentList from '@/components/CommentList.vue'
import ReportDialog from '@/components/ReportDialog.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const goods = ref(null)
const isFavorited = ref(false)
const showReportDialog = ref(false)
const showWantDialog = ref(false)
const showExchangeDialog = ref(false)
const submitting = ref(false)
const myGoodsList = ref([])
const currentImageIndex = ref(0)

const sellerCredit = ref(5) // 默认信誉分

const wantForm = ref({
  message: ''
})

const exchangeForm = ref({
  my_goods_id: null,
  message: ''
})

// 图片列表
const imageList = computed(() => {
  if (!goods.value) return []
  if (typeof goods.value.images === 'string') {
    try {
      return JSON.parse(goods.value.images)
    } catch {
      return []
    }
  }
  return goods.value.images || []
})

// 是否有有效的用户信息
const hasValidSeller = computed(() => {
  return goods.value?.user && goods.value.user.id > 0
})

// 获取商品详情
const fetchGoodsDetail = async () => {
  loading.value = true
  try {
    const goodsId = route.params.id
    const res = await getGoodsDetail(goodsId)
    goods.value = res.data.goods
    isFavorited.value = res.data.is_favorited

    // 设置信誉分
    if (goods.value.user?.credit_score) {
      sellerCredit.value = Math.floor(goods.value.user.credit_score / 20)
    }
  } catch (error) {
    console.error('获取商品详情失败:', error)
    ElMessage.error('获取商品详情失败')
  } finally {
    loading.value = false
  }
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

// 获取性别文本
const getGenderText = (gender) => {
  const map = {
    1: '男',
    2: '女'
  }
  return map[gender] || '未知'
}

// 获取类型文本
const getTypeText = (type) => {
  const map = {
    1: '售卖',
    2: '交换',
    3: '均可'
  }
  return map[type] || '未知'
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp || timestamp === 0) return '暂无'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取性别显示信息
const getUserGenderDisplay = (gender) => {
  if (!gender || gender === 0) {
    return { text: '未设置', type: 'info' }
  }
  const map = {
    1: { text: '男', type: 'primary' },
    2: { text: '女', type: 'danger' }
  }
  return map[gender] || { text: '未知', type: 'info' }
}

// 返回
const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/')
}

// 跳转到用户详情页
const goToUserDetail = () => {
  console.log('goToUserDetail 被调用', {
    hasValidSeller: hasValidSeller.value,
    user: goods.value?.user
  })

  if (!hasValidSeller.value) {
    ElMessage.warning('用户信息加载中或不存在')
    return
  }

  const userId = goods.value.user.id
  console.log('准备跳转到用户详情页，userId:', userId)

  router.push(`/user/${userId}`)
}

// 图片导航
const previousImage = () => {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value--
  }
}

const nextImage = () => {
  if (currentImageIndex.value < imageList.value.length - 1) {
    currentImageIndex.value++
  }
}

// 我想要
const handleWant = async () => {
  if (!userStore.isLogin) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  // 检查是否是自己的商品
  if (goods.value.user_id === userStore.userId) {
    ElMessage.warning('不能购买自己的商品')
    return
  }

  // 打开购买对话框
  wantForm.value.message = ''
  showWantDialog.value = true
}

// 发起交换
const handleExchange = async () => {
  if (!userStore.isLogin) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  // 检查是否是自己的商品
  if (goods.value.user_id === userStore.userId) {
    ElMessage.warning('不能和自己交换商品')
    return
  }

  // 获取我的商品列表（只获取状态为在售的商品）
  try {
    const res = await getMyGoods({
      page: 1,
      page_size: 100,
      status: 1
    })

    // 处理图片数据并过滤掉当前商品
    myGoodsList.value = (res.data.list || [])
      .filter(item => item.id !== goods.value.id)
      .map(item => {
        return {
          ...item,
          getImages: () => {
            if (typeof item.images === 'string') {
              try {
                return JSON.parse(item.images)
              } catch {
                return []
              }
            }
            return item.images || []
          }
        }
      })

    if (myGoodsList.value.length === 0) {
      ElMessage.warning('您还没有可交换的商品，请先发布商品')
      router.push('/publish')
      return
    }

    // 打开交换对话框
    exchangeForm.value.my_goods_id = null
    exchangeForm.value.message = ''
    showExchangeDialog.value = true
  } catch (error) {
    console.error('获取我的商品失败:', error)
    ElMessage.error('获取我的商品失败')
  }
}

// 提交购买请求
const submitWant = async () => {
  try {
    submitting.value = true

    await createExchange({
      goods_id: goods.value.id,
      type: 1, // 购买
      message: wantForm.value.message
    })

    ElMessage.success('购买请求已发送给卖家')
    showWantDialog.value = false
    wantForm.value.message = ''
  } catch (error) {
    console.error('购买失败:', error)
    ElMessage.error(error.response?.data?.message || '购买失败')
  } finally {
    submitting.value = false
  }
}

// 提交交换请求
const submitExchange = async () => {
  if (!exchangeForm.value.my_goods_id) {
    ElMessage.warning('请选择要交换的商品')
    return
  }

  try {
    submitting.value = true

    await createExchange({
      goods_id: goods.value.id,
      my_goods_id: exchangeForm.value.my_goods_id,
      type: 2, // 交换
      message: exchangeForm.value.message
    })

    ElMessage.success('交换请求已发送给卖家')
    showExchangeDialog.value = false
    exchangeForm.value.my_goods_id = null
    exchangeForm.value.message = ''
  } catch (error) {
    console.error('发起交换失败:', error)
    ElMessage.error(error.response?.data?.message || '发起交换失败')
  } finally {
    submitting.value = false
  }
}

// 收藏
const handleFavorite = async () => {
  if (!userStore.isLogin) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  try {
    if (isFavorited.value) {
      // 取消收藏
      await removeFavorite(route.params.id)
      isFavorited.value = false
      goods.value.favorite_count--
      ElMessage.success('取消收藏成功')
    } else {
      // 添加收藏
      await addFavorite(route.params.id)
      isFavorited.value = true
      goods.value.favorite_count++
      ElMessage.success('收藏成功')
    }
  } catch (error) {
    console.error('收藏操作失败:', error)
    ElMessage.error(isFavorited.value ? '取消收藏失败' : '收藏失败')
  }
}

// 举报成功回调
const handleReportSuccess = () => {
  ElMessage.success('举报成功，感谢您的反馈')
}

onMounted(() => {
  fetchGoodsDetail()
})
</script>

<style scoped>
.goods-detail {
  min-height: 100vh;
  background-color: #f5f7fa;
  display: flex;
  flex-direction: column;
}

/* 头部导航 */
.el-header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 0 30px;
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
  font-weight: 600;
  color: #409eff;
}

.nav {
  display: flex;
  gap: 25px;
}

.nav a {
  color: #606266;
  text-decoration: none;
  cursor: pointer;
  font-size: 15px;
  transition: color 0.3s;
}

.nav a:hover {
  color: #409eff;
}

/* 主内容区 */
.el-main {
  width: 100%;
  height: calc(100vh - 60px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.detail-main {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.detail-container {
  width: 80vw;
  height: 80vh;
  max-width: 80vw;
  max-height: 80vh;
  overflow-y: auto;
  padding-right: 5px;
}

.detail-container::-webkit-scrollbar {
  width: 6px;
}

.detail-container::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.detail-container::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

/* 主卡片 */
.main-card {
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 20px;
  width: 100%;
}

.main-card::-webkit-scrollbar {
  width: 6px;
}

.main-card::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.main-card::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}

/* 图片区域 */
.image-section {
  margin-bottom: 15px;
  width: 100%;
}

.image-gallery {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  width: 100%;
  min-width: 100%;
}

.main-image-container {
  position: relative;
  width: 100%;
  min-width: 100%;
  height: 380px;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

:deep(.el-image) {
  width: 100% !important;
  max-width: 100% !important;
  height: 100% !important;
  min-width: 100% !important;
}

:deep(.el-image__inner) {
  width: 100% !important;
  max-width: 100% !important;
  height: 100% !important;
  object-fit: contain;
}

.main-image {
  width: 100%;
  height: 100%;
}

.image-actions {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 10px;
  z-index: 10;
}

.image-nav-btn {
  background: rgba(255, 255, 255, 0.95);
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  width: 40px;
  height: 40px;
  flex-shrink: 0;
}

.image-nav-btn:hover {
  background: #fff;
  transform: scale(1.05);
}

.no-image {
  position: absolute;
  bottom: 15px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
}

.image-nav-btn {
  background: rgba(255, 255, 255, 0.95);
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  width: 36px;
  height: 36px;
}

.image-nav-btn:hover {
  background: #fff;
  transform: scale(1.05);
}

.no-image {
  height: 450px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
}

.no-image-content {
  text-align: center;
  color: #909399;
}

.no-image-content p {
  margin-top: 10px;
  font-size: 14px;
}

/* 缩略图 */
.thumbnail-container {
  display: flex;
  gap: 8px;
  padding: 10px;
  background: #fff;
  overflow-x: auto;
}

.thumbnail {
  width: 60px;
  height: 60px;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
  flex-shrink: 0;
}

.thumbnail:hover {
  border-color: #c0c4cc;
  transform: translateY(-2px);
}

.thumbnail.active {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.thumbnail .el-image {
  width: 100%;
  height: 100%;
}

/* 商品信息区域 */
.info-section {
  padding: 0 20px;
  width: 100%;
  min-width: 100%;
}

.header-section {
  margin-bottom: 15px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 10px 0;
  color: #303133;
  line-height: 1.4;
}

.tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

/* 价格卡片行 */
.price-row {
  margin-bottom: 18px;
  width: 100%;
}

.price-card {
  padding: 18px;
  background: linear-gradient(135deg, #fff7f0 0%, #ffe8dc 100%);
  border-radius: 10px;
  border: 1px solid #ffe4d1;
  width: 100%;
}

/* 价格和操作按钮行 */
.price-action-row {
  margin-bottom: 18px;
  width: 100%;
}

.price-card {
  padding: 18px;
  background: linear-gradient(135deg, #fff7f0 0%, #ffe8dc 100%);
  border-radius: 10px;
  border: 1px solid #ffe4d1;
}

.price-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 15px;
}

.price-section {
  flex: 1;
}

.price-value {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.currency {
  font-size: 20px;
  font-weight: 600;
  color: #f56c6c;
}

.amount {
  font-size: 34px;
  font-weight: 700;
  color: #f56c6c;
}

.original-price {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
  font-size: 13px;
  color: #909399;
}

.discount {
  background: #f56c6c;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.exchange-only {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  padding: 20px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  color: #0ea5e9;
  border-radius: 8px;
}

.exchange-only .el-icon {
  font-size: 24px;
}

/* 操作按钮行 */
.action-buttons-row {
  display: flex;
  gap: 12px;
  margin-bottom: 18px;
}

.main-action-btn {
  flex: 1;
  height: 48px;
  border-radius: 8px;
  font-weight: 500;
  font-size: 15px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.favorite-btn {
  flex: 0 0 auto;
  min-width: 120px;
}

.favorite-btn.active {
  background: #f56c6c;
  border-color: #f56c6c;
  color: white;
}

.action-icon-btn {
  flex: 0 0 auto;
  min-width: 110px;
  height: 48px;
  background: #fef0f0;
  border: 1px solid #fde2e2;
  color: #f56c6c;
  border-radius: 8px;
}

.action-icon-btn:hover {
  background: #f56c6c;
  border-color: #f56c6c;
  color: white;
}

/* 详情行 */
.details-row {
  margin-bottom: 0;
}

/* 交换提示条 */
.exchange-tip-bar {
  margin-bottom: 18px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-radius: 10px;
  border: 1px solid #bae6fd;
  box-shadow: 0 2px 8px rgba(56, 189, 248, 0.1);
  width: 100%;
}

/* 操作按钮行 */
.action-buttons-row {
  display: flex;
  gap: 12px;
  margin-bottom: 18px;
  width: 100%;
}

.tip-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.tip-icon {
  font-size: 32px;
  color: #0ea5e9;
  flex-shrink: 0;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.tip-text-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tip-title {
  font-size: 16px;
  font-weight: 600;
  color: #0369a1;
}

.tip-desc {
  font-size: 13px;
  color: #0c4a6e;
}

/* 详情卡片 */
.details-card {
  padding: 16px;
  background: #fff;
  border-radius: 10px;
  border: 1px solid #ebeef5;
  transition: all 0.3s;
}

.details-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.details-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 14px;
  color: #303133;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 2px solid #f5f7fa;
}

.seller-card .details-card-header {
  justify-content: space-between;
}

.seller-rating {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 13px;
  color: #ff9900;
  font-weight: 600;
}

/* 属性列表 */
.attribute-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.attribute-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.attr-label {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.attr-value {
  font-size: 13px;
  color: #303133;
  font-weight: 600;
}

/* 卖家信息 */
.seller-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: all 0.3s;
}

.seller-info.clickable {
  cursor: pointer;
  pointer-events: auto;
  user-select: none;
}

.seller-info.clickable:hover {
  background: #ecf5ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.seller-info.clickable:hover .arrow-icon {
  color: #409eff;
  transform: translateX(3px);
}

.seller-info.clickable:active {
  transform: translateY(0);
}

.seller-avatar {
  flex-shrink: 0;
}

.seller-detail {
  flex: 1;
  min-width: 0;
}

.name {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 5px;
}

.seller-meta {
  display: flex;
  gap: 12px;
  margin-bottom: 5px;
  flex-wrap: wrap;
}

.meta-text {
  font-size: 12px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 3px;
}

.meta-text .el-icon {
  font-size: 13px;
}

.arrow-icon {
  color: #c0c4cc;
  font-size: 16px;
  transition: all 0.3s;
}

/* 描述和地点区域 */
.description-section,
.location-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 15px;
  color: #303133;
  margin-bottom: 10px;
}

.description-text {
  font-size: 14px;
  line-height: 1.7;
  color: #606266;
  white-space: pre-wrap;
  padding: 12px 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.location-section .location-content {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #606266;
  padding: 10px 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.location-icon {
  font-size: 18px;
  color: #409eff;
}

.location-text {
  flex: 1;
}

/* 评论区域 */
.comment-section {
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  border-radius: 12px;
  width: 100%;
}

/* 对话框样式 */
.price-highlight {
  margin-left: 10px;
  font-size: 16px;
  font-weight: 600;
  color: #f56c6c;
}

.tip {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
}

:deep(.custom-dialog) {
  border-radius: 8px;
}

:deep(.custom-dialog .el-dialog__header) {
  background: #409eff;
  color: white;
  padding: 16px 20px;
  margin: 0;
}

:deep(.custom-dialog .el-dialog__title) {
  color: white;
  font-size: 16px;
  font-weight: 600;
}

:deep(.custom-dialog .el-dialog__body) {
  padding: 20px;
}

:deep(.custom-dialog .el-dialog__footer) {
  padding: 16px 20px;
  background: #f5f7fa;
  border-top: 1px solid #ebeef5;
}

.dialog-goods-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 6px;
}

.dialog-goods-title {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.dialog-goods-price {
  font-size: 18px;
  color: #f56c6c;
  font-weight: 600;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.exchange-goods-card {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 6px;
}

.exchange-goods-image {
  width: 60px;
  height: 60px;
  border-radius: 6px;
  flex-shrink: 0;
}

.exchange-goods-info {
  flex: 1;
  display: flex;
  align-items: center;
}

.exchange-goods-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.select-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.select-option-image {
  width: 32px;
  height: 32px;
  border-radius: 4px;
  flex-shrink: 0;
}

.select-option-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 响应式 */
@media (max-width: 768px) {
  .el-main {
    padding: 15px;
  }

  .detail-container {
    width: 95vw;
    height: 90vh;
    max-width: 95vw;
    max-height: 90vh;
  }

  .action-buttons-row {
    flex-wrap: wrap;
    gap: 10px;
  }

  .main-action-btn {
    flex: 1 1 45%;
    min-width: 45%;
    height: 44px;
    font-size: 14px;
  }

  .favorite-btn {
    min-width: 100px;
    height: 44px;
  }

  .action-icon-btn {
    min-width: 90px;
    height: 44px;
  }

  .title {
    font-size: 20px;
  }

  .amount {
    font-size: 28px;
  }

  .main-image-container {
    height: 280px;
  }

  .no-image {
    height: 280px;
  }

  .exchange-tip-bar {
    padding: 12px 15px;
  }

  .tip-icon {
    font-size: 24px;
  }

  .tip-title {
    font-size: 14px;
  }

  .tip-desc {
    font-size: 12px;
  }
}
</style>
