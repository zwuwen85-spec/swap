<template>
  <div class="user-detail">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link v-if="userStore.isLogin" to="/profile">个人中心</router-link>
            <router-link v-else to="/login">登录</router-link>
          </div>
        </div>
      </el-header>

      <el-main v-loading="loading" class="detail-main">
        <div v-if="userInfo" class="detail-container">
          <!-- 用户信息卡片 -->
          <el-card class="user-card">
            <div class="user-header">
              <el-avatar :size="100" :src="userInfo.avatar">
                {{ userInfo.username?.[0] }}
              </el-avatar>
              <div class="user-basic-info">
                <h2>{{ userInfo.nickname || userInfo.username }}</h2>
                <p class="username">@{{ userInfo.username }}</p>
                <div class="credit-section">
                  <el-rate
                    v-model="creditScore"
                    disabled
                    show-score
                    text-color="#ff9900"
                    score-template="{value}"
                    size="small"
                  />
                  <span class="credit-label">信誉分：{{ userInfo.credit_score }}</span>
                </div>
              </div>
            </div>

            <el-divider class="section-divider" />

            <div class="user-details">
              <el-row :gutter="20">
                <el-col :span="12">
                  <div class="detail-section">
                    <h4 class="section-title">学籍信息</h4>
                    <div class="info-list">
                      <div class="info-item" v-if="userInfo.school">
                        <el-icon class="info-icon"><School /></el-icon>
                        <span class="info-label">学校</span>
                        <span class="info-value">{{ userInfo.school }}</span>
                      </div>
                      <div class="info-item" v-if="userInfo.student_id">
                        <el-icon class="info-icon"><Postcard /></el-icon>
                        <span class="info-label">学号</span>
                        <span class="info-value">{{ userInfo.student_id }}</span>
                      </div>
                      <div class="info-item" v-if="userInfo.major">
                        <el-icon class="info-icon"><Reading /></el-icon>
                        <span class="info-label">专业</span>
                        <span class="info-value">{{ userInfo.major }}</span>
                      </div>
                      <div class="info-item">
                        <el-icon class="info-icon"><User /></el-icon>
                        <span class="info-label">性别</span>
                        <span class="info-value">
                          {{ getGenderText(userInfo.gender) }}
                        </span>
                      </div>
                    </div>
                  </div>
                </el-col>

                <el-col :span="12">
                  <div class="detail-section">
                    <h4 class="section-title">联系方式</h4>
                    <div class="info-list">
                      <div class="info-item" v-if="userStore.isLogin && userStore.userId !== userInfo.id">
                        <el-button type="primary" size="small" @click="goToChat(userInfo.id)">
                          <el-icon><ChatDotRound /></el-icon>发私信
                        </el-button>
                      </div>
                      <div class="info-item" v-if="userInfo.phone">
                        <el-icon class="info-icon"><Phone /></el-icon>
                        <span class="info-label">手机</span>
                        <span class="info-value">{{ userInfo.phone }}</span>
                      </div>
                      <div class="info-item" v-if="userInfo.email">
                        <el-icon class="info-icon"><Message /></el-icon>
                        <span class="info-label">邮箱</span>
                        <span class="info-value">{{ userInfo.email }}</span>
                      </div>
                      <div class="info-item" v-if="userInfo.qq">
                        <el-icon class="info-icon"><ChatDotRound /></el-icon>
                        <span class="info-label">QQ</span>
                        <span class="info-value">{{ userInfo.qq }}</span>
                      </div>
                      <div class="info-item" v-if="userInfo.wechat">
                        <el-icon class="info-icon"><ChatDotRound /></el-icon>
                        <span class="info-label">微信</span>
                        <span class="info-value">{{ userInfo.wechat }}</span>
                      </div>
                    </div>
                  </div>
                </el-col>
              </el-row>

              <el-empty
                v-if="!hasContactInfo && !userInfo.school && !userInfo.student_id && !userInfo.major"
                description="该用户暂未完善个人信息"
                :image-size="100"
              />
            </div>
          </el-card>

          <!-- 用户商品列表 -->
          <el-card class="goods-card">
            <template #header>
              <div class="card-header">
                <el-icon><ShoppingBag /></el-icon>
                <h3>TA的商品</h3>
              </div>
            </template>
            <div v-if="userGoods.length > 0" class="goods-list">
              <div
                v-for="goods in userGoods"
                :key="goods.id"
                class="goods-item"
                @click="goToGoods(goods.id)"
              >
                <el-image
                  :src="getGoodsImage(goods)"
                  fit="cover"
                  class="goods-image"
                />
                <div class="goods-info">
                  <div class="goods-title">{{ goods.title }}</div>
                  <div class="goods-price">
                    <span v-if="goods.type === 1 || goods.type === 3" class="price">
                      ¥{{ goods.price }}
                    </span>
                    <el-tag v-if="goods.type === 2" type="success" size="small">交换</el-tag>
                    <el-tag v-if="goods.type === 3" type="warning" size="small">可售可换</el-tag>
                  </div>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无商品" />
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getUserInfoById } from '@/api/user'
import { getGoodsByUserId } from '@/api/goods'
import { ElMessage } from 'element-plus'
import {
  School,
  Postcard,
  Reading,
  User,
  Phone,
  Message,
  ChatDotRound,
  ShoppingBag
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const userInfo = ref(null)
const userGoods = ref([])

const creditScore = ref(5)

// 是否有联系方式
const hasContactInfo = computed(() => {
  if (!userInfo.value) return false
  return !!(
    userInfo.value.phone ||
    userInfo.value.email ||
    userInfo.value.qq ||
    userInfo.value.wechat
  )
})

// 是否有学籍信息
const hasStudentInfo = computed(() => {
  if (!userInfo.value) return false
  return !!(
    userInfo.value.school ||
    userInfo.value.student_id ||
    userInfo.value.major ||
    userInfo.value.gender !== 0
  )
})

// 获取性别文本
const getGenderText = (gender) => {
  const map = {
    1: '男',
    2: '女'
  }
  return map[gender] || '未设置'
}

// 获取性别标签类型
const getGenderTagType = (gender) => {
  const map = {
    1: 'primary',  // 男 - 蓝色
    2: 'danger'    // 女 - 红色
  }
  return map[gender] || 'info'  // 未设置 - 灰色
}

// 获取商品图片
const getGoodsImage = (goods) => {
  if (typeof goods.images === 'string') {
    try {
      const images = JSON.parse(goods.images)
      return images.length > 0 ? images[0] : ''
    } catch {
      return ''
    }
  }
  return goods.images?.[0] || ''
}

// 获取用户信息
const fetchUserInfo = async () => {
  loading.value = true
  try {
    const userId = route.params.id
    const res = await getUserInfoById(userId)
    userInfo.value = res.data

    // 设置信誉分
    if (userInfo.value.credit_score) {
      creditScore.value = Math.floor(userInfo.value.credit_score / 20)
    }

    // 获取用户的商品
    await fetchUserGoods(userId)
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

// 获取用户商品
const fetchUserGoods = async (userId) => {
  try {
    const res = await getGoodsByUserId(userId, {
      page: 1,
      page_size: 10
    })
    userGoods.value = res.data.list || []
  } catch (error) {
    console.error('获取用户商品失败:', error)
  }
}

// 跳转到商品详情
const goToGoods = (goodsId) => {
  router.push(`/goods/${goodsId}`)
}

// 跳转到私聊
const goToChat = (userId) => {
  router.push(`/chat/${userId}`)
}

const goHome = () => {
  router.push('/')
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.user-detail {
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

.el-main {
  max-width: 100%;
  width: 100%;
  padding: 20px 40px;
}

.detail-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 85%;
  margin: 0 auto;
}

.user-card {
  margin-bottom: 0;
}

.user-header {
  display: flex;
  gap: 20px;
  align-items: center;
  padding: 10px 0;
}

.user-basic-info {
  flex: 1;
}

.user-basic-info h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.username {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #909399;
}

.credit-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.credit-label {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
}

.section-divider {
  margin: 20px 0;
}

.user-details {
  padding: 10px 0;
}

.detail-section {
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: all 0.3s;
}

.detail-section:hover {
  background: #f0f2f5;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.section-title {
  margin: 0 0 15px 0;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  padding-bottom: 10px;
  border-bottom: 2px solid #e0e0e0;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: white;
  border-radius: 6px;
  transition: all 0.2s;
}

.info-item:hover {
  background: #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  transform: translateX(3px);
}

.info-icon {
  font-size: 18px;
  color: #409eff;
  flex-shrink: 0;
}

.info-label {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
  min-width: 40px;
  flex-shrink: 0;
}

.info-value {
  flex: 1;
  font-size: 14px;
  color: #303133;
  font-weight: 500;
  word-break: break-all;
}

.goods-card {
  margin-bottom: 0;
}

.goods-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.goods-item {
  cursor: pointer;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.goods-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.goods-image {
  width: 100%;
  height: 150px;
}

.goods-info {
  padding: 10px;
}

.goods-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.goods-price {
  display: flex;
  align-items: center;
  gap: 8px;
}

.price {
  font-size: 18px;
  font-weight: 600;
  color: #f56c6c;
}

/* 响应式 */
@media (max-width: 768px) {
  .el-main {
    padding: 15px 20px;
  }

  .detail-container {
    max-width: 100%;
  }

  .user-header {
    flex-direction: column;
    text-align: center;
  }

  .user-basic-info h2 {
    font-size: 20px;
  }

  .credit-section {
    justify-content: center;
  }

  .detail-section {
    margin-bottom: 15px;
  }

  .goods-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1400px) {
  .detail-container {
    max-width: 1200px;
  }
}
</style>
