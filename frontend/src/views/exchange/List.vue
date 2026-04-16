<template>
  <div class="exchange-list">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/exchange" class="active">我的交换</router-link>
            <router-link to="/publish">发布商品</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
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
                  <span>交换筛选</span>
                </div>
              </div>

              <el-row :gutter="24" class="filter-row">
                <!-- 交换类型 -->
                <el-col :xs="24" :sm="12" :md="8" :lg="8">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><ChatDotRound /></el-icon>
                      交换类型
                    </label>
                    <el-radio-group
                      v-model="type"
                      @change="handleTypeChange"
                      class="type-radio-group"
                      size="default"
                    >
                      <el-radio-button value="incoming">
                        <el-icon><Download /></el-icon>
                        收到的请求 ({{ incomingCount }})
                      </el-radio-button>
                      <el-radio-button value="outgoing">
                        <el-icon><Upload /></el-icon>
                        我发起的
                      </el-radio-button>
                      <el-radio-button value="all">
                        <el-icon><List /></el-icon>
                        全部
                      </el-radio-button>
                    </el-radio-group>
                  </div>
                </el-col>

                <!-- 状态筛选 -->
                <el-col :xs="24" :sm="12" :md="8" :lg="8">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><CircleCheck /></el-icon>
                      交换状态
                    </label>
                    <el-select
                      v-model="statusFilter"
                      placeholder="全部状态"
                      clearable
                      @change="fetchList"
                      size="large"
                      class="filter-select"
                    >
                      <el-option label="全部状态" :value="-1" />
                      <el-option label="待处理" :value="0">
                        <span>待处理</span>
                        <el-tag type="warning" size="small" style="margin-left: 8px">0</el-tag>
                      </el-option>
                      <el-option label="已接受" :value="1">
                        <span>已接受</span>
                        <el-tag type="success" size="small" style="margin-left: 8px">1</el-tag>
                      </el-option>
                      <el-option label="已拒绝" :value="2">
                        <span>已拒绝</span>
                        <el-tag type="danger" size="small" style="margin-left: 8px">2</el-tag>
                      </el-option>
                      <el-option label="已取消" :value="3">
                        <span>已取消</span>
                        <el-tag type="info" size="small" style="margin-left: 8px">3</el-tag>
                      </el-option>
                      <el-option label="已完成" :value="4">
                        <span>已完成</span>
                        <el-tag size="small" style="margin-left: 8px">4</el-tag>
                      </el-option>
                    </el-select>
                  </div>
                </el-col>

                <!-- 快速筛选 -->
                <el-col :xs="24" :sm="24" :md="8" :lg="8">
                  <div class="filter-item">
                    <label class="filter-label">
                      <el-icon><Lightning /></el-icon>
                      快速筛选
                    </label>
                    <div class="quick-filters">
                      <el-tag
                        :type="statusFilter === 0 ? 'warning' : 'info'"
                        effect="plain"
                        @click="statusFilter = statusFilter === 0 ? -1 : 0; fetchList()"
                        class="quick-tag"
                      >
                        待处理
                      </el-tag>
                      <el-tag
                        :type="statusFilter === 1 ? 'success' : 'info'"
                        effect="plain"
                        @click="statusFilter = statusFilter === 1 ? -1 : 1; fetchList()"
                        class="quick-tag"
                      >
                        已接受
                      </el-tag>
                      <el-tag
                        :type="statusFilter === 4 ? 'primary' : 'info'"
                        effect="plain"
                        @click="statusFilter = statusFilter === 4 ? -1 : 4; fetchList()"
                        class="quick-tag"
                      >
                        已完成
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
                  v-if="type !== 'all'"
                  closable
                  @close="type = 'all'; handleTypeChange()"
                  class="filter-tag"
                  type="primary"
                >
                  类型: {{ type === 'incoming' ? '收到的请求' : '我发起的' }}
                </el-tag>
                <el-tag
                  v-if="statusFilter >= 0"
                  closable
                  @close="statusFilter = -1; fetchList()"
                  class="filter-tag"
                  type="success"
                >
                  状态: {{ getStatusText(statusFilter) }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </div>

        <!-- 交换列表 -->
        <div class="exchange-wrapper" v-loading="loading">
          <el-empty v-if="exchangeList.length === 0 && !loading" description="暂无交换记录" />

          <div v-else class="exchange-grid">
            <div
              v-for="exchange in exchangeList"
              :key="exchange.id"
              class="exchange-item"
            >
              <!-- 交换卡片 -->
              <el-card class="exchange-card" shadow="hover">
                <!-- 顶部状态栏 -->
                <div class="card-header">
                  <div class="status-main">
                    <el-tag :type="getStatusType(exchange.status)" size="large" effect="dark">
                      {{ exchange.get_status_text ? exchange.get_status_text() : getStatusText(exchange.status) }}
                    </el-tag>
                    <el-tag
                      :type="exchange.type === 1 ? 'danger' : 'success'"
                      effect="plain"
                      size="large"
                    >
                      {{ exchange.get_type_text ? exchange.get_type_text() : getTypeText(exchange.type) }}
                    </el-tag>
                  </div>
                  <div class="status-tip">
                    <el-icon><InfoFilled /></el-icon>
                    {{ getStatusTip(exchange) }}
                  </div>
                </div>

                <!-- 状态流程步骤条 -->
                <div class="status-progress">
                  <el-steps :active="getStatusStep(exchange.status)" finish-status="success" align-center>
                    <el-step title="发起请求" :icon="Promotion" />
                    <el-step title="等待处理" :icon="Clock" />
                    <el-step
                      :title="exchange.status >= 2 && exchange.status !== 3 ? '被拒绝' : '已接受'"
                      :icon="exchange.status === 2 ? CloseBold : Select"
                    />
                    <el-step title="完成交换" :icon="CircleCheck" />
                  </el-steps>
                </div>

                <!-- 商品图片区域 -->
                <div class="goods-images-section">
                  <!-- 目标商品 -->
                  <div class="goods-image-box" @click="goGoodsDetail(exchange.goods_id)">
                    <div class="image-label">
                      <el-icon><Aim /></el-icon>
                      目标商品
                    </div>
                    <div v-if="exchange.goods && exchange.goods.getImages && exchange.goods.getImages().length > 0" class="image-wrapper">
                      <el-image
                        :src="exchange.goods.getImages()[0]"
                        fit="cover"
                        class="goods-image"
                      >
                        <template #error>
                          <div class="image-error">
                            <el-icon><Picture /></el-icon>
                          </div>
                        </template>
                      </el-image>
                    </div>
                    <div v-else class="image-wrapper no-image">
                      <el-icon><Picture /></el-icon>
                    </div>
                    <div class="goods-title">{{ exchange.goods?.title || '未知商品' }}</div>
                    <div v-if="exchange.type === 1 || exchange.type === 3" class="goods-price">
                      ¥{{ exchange.goods?.price }}
                    </div>
                  </div>

                  <!-- 交换图标 -->
                  <div v-if="exchange.type === 2" class="exchange-icon">
                    <el-icon><Right /></el-icon>
                  </div>

                  <!-- 我的商品（物物交换） -->
                  <div v-if="exchange.type === 2" class="goods-image-box" @click="goGoodsDetail(exchange.my_goods_id)">
                    <div class="image-label my-goods">
                      <el-icon><Goods /></el-icon>
                      我的商品
                    </div>
                    <div v-if="exchange.my_goods && exchange.my_goods.getImages && exchange.my_goods.getImages().length > 0" class="image-wrapper">
                      <el-image
                        :src="exchange.my_goods.getImages()[0]"
                        fit="cover"
                        class="goods-image"
                      >
                        <template #error>
                          <div class="image-error">
                            <el-icon><Picture /></el-icon>
                          </div>
                        </template>
                      </el-image>
                    </div>
                    <div v-else class="image-wrapper no-image">
                      <el-icon><Picture /></el-icon>
                    </div>
                    <div class="goods-title">{{ exchange.my_goods?.title || '未知商品' }}</div>
                  </div>
                </div>

                <!-- 用户信息 -->
                <div class="user-section">
                  <div class="user-info-box">
                    <el-avatar :size="32" :src="getOtherUser(exchange)?.avatar">
                      {{ getOtherUser(exchange)?.username?.[0] }}
                    </el-avatar>
                    <div class="user-detail">
                      <div class="user-label">
                        {{ exchange.initiator_id === currentUserId ? '对方用户' : '发起人' }}
                      </div>
                      <div class="user-name">{{ getOtherUser(exchange)?.username || '未知用户' }}</div>
                    </div>
                  </div>
                </div>

                <!-- 留言/拒绝原因 -->
                <div v-if="exchange.message || (exchange.status === 2 && exchange.reject_reason)" class="message-section">
                  <div v-if="exchange.message" class="message-box">
                    <div class="message-label">
                      <el-icon><ChatLineSquare /></el-icon>
                      留言
                    </div>
                    <p>{{ exchange.message }}</p>
                  </div>
                  <div v-if="exchange.status === 2 && exchange.reject_reason" class="message-box reject">
                    <div class="message-label">
                      <el-icon><Warning /></el-icon>
                      拒绝原因
                    </div>
                    <p>{{ exchange.reject_reason }}</p>
                  </div>
                </div>

                <!-- 操作按钮 -->
                <div class="actions-section">
                  <!-- 待处理状态 -->
                  <template v-if="exchange.status === 0">
                    <!-- 对方发来的请求 -->
                    <template v-if="exchange.target_id === currentUserId">
                      <el-button
                        type="success"
                        size="large"
                        @click="handleAccept(exchange)"
                        class="action-btn accept"
                      >
                        <el-icon><Select /></el-icon>
                        接受请求
                      </el-button>
                      <el-button
                        type="danger"
                        size="large"
                        @click="handleReject(exchange)"
                        class="action-btn reject"
                      >
                        <el-icon><CloseBold /></el-icon>
                        拒绝请求
                      </el-button>
                    </template>
                    <!-- 我发起的请求 -->
                    <template v-else>
                      <el-button
                        size="large"
                        @click="handleCancel(exchange)"
                        class="action-btn cancel"
                      >
                        <el-icon><Close /></el-icon>
                        取消请求
                      </el-button>
                      <div class="waiting-tip">
                        <el-icon class="is-loading"><Loading /></el-icon>
                        等待对方处理
                      </div>
                    </template>
                  </template>

                  <!-- 已接受状态 -->
                  <template v-if="exchange.status === 1">
                    <el-button
                      type="primary"
                      size="large"
                      @click="handleComplete(exchange)"
                      class="action-btn complete"
                    >
                      <el-icon><CircleCheck /></el-icon>
                      确认完成交换
                    </el-button>
                    <div class="action-tip">
                      <el-icon><InfoFilled /></el-icon>
                      双方已达成一致，请确认完成
                    </div>
                  </template>

                  <!-- 已拒绝/已取消状态 -->
                  <template v-if="exchange.status === 2 || exchange.status === 3">
                    <div class="status-finished">
                      <el-icon :class="exchange.status === 2 ? 'status-rejected' : 'status-cancelled'">
                        <CloseBold />
                      </el-icon>
                      {{ exchange.status === 2 ? '此交换已被拒绝' : '此交换已取消' }}
                    </div>
                    <el-button
                      size="large"
                      @click="viewDetail(exchange)"
                      class="action-btn detail"
                    >
                      <el-icon><View /></el-icon>
                      查看详情
                    </el-button>
                  </template>

                  <!-- 已完成状态 -->
                  <template v-if="exchange.status === 4">
                    <div class="status-completed">
                      <el-icon class="status-success"><CircleCheck /></el-icon>
                      交换已完成！
                    </div>
                    <el-button
                      size="large"
                      @click="goGoodsDetail(exchange.goods_id)"
                      class="action-btn"
                    >
                      <el-icon><View /></el-icon>
                      再次查看商品
                    </el-button>
                  </template>
                </div>

                <!-- 底部时间 -->
                <div class="time-section">
                  <el-icon><Clock /></el-icon>
                  {{ formatTime(exchange.create_time) }}
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

    <!-- 拒绝对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="拒绝交换" width="500px">
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="拒绝原因">
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmReject">确认拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getExchangeList, handleExchange, getPendingCount } from '@/api/exchange'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Filter, ChatDotRound, CircleCheck, Lightning, Download, Upload, List,
  Picture, Aim, Goods, Right, ChatLineSquare, Warning, Select, CloseBold,
  Close, View, Clock, InfoFilled, Promotion, Loading
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const currentUserId = computed(() => userStore.userId)

const loading = ref(false)
const exchangeList = ref([])
const total = ref(0)
const incomingCount = ref(0)
const type = ref('incoming')
const statusFilter = ref(-1)

const pagination = reactive({
  page: 1,
  page_size: 12
})

const rejectDialogVisible = ref(false)
const rejectForm = reactive({
  exchangeId: null,
  reason: ''
})

// 计算是否有活跃筛选
const hasActiveFilters = computed(() => {
  return type.value !== 'all' || statusFilter.value >= 0
})

// 获取交换列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      type: type.value,
      page: pagination.page,
      page_size: pagination.page_size
    }

    if (statusFilter.value >= 0) {
      params.status = statusFilter.value
    }

    const res = await getExchangeList(params)
    exchangeList.value = res.data.list
    total.value = res.data.total

    // 提前获取当前用户ID，避免在forEach中访问computed的问题
    const userId = currentUserId.value

    // 处理图片数据
    exchangeList.value.forEach(exchange => {
      if (exchange.goods) {
        exchange.goods.getImages = () => {
          if (typeof exchange.goods.images === 'string') {
            try {
              return JSON.parse(exchange.goods.images)
            } catch {
              return []
            }
          }
          return exchange.goods.images || []
        }
      }
      if (exchange.my_goods) {
        exchange.my_goods.getImages = () => {
          if (typeof exchange.my_goods.images === 'string') {
            try {
              return JSON.parse(exchange.my_goods.images)
            } catch {
              return []
            }
          }
          return exchange.my_goods.images || []
        }
      }

      // 添加辅助方法
      exchange.get_status_text = () => exchange.status === 0 ? '待处理' :
                                     exchange.status === 1 ? '已接受' :
                                     exchange.status === 2 ? '已拒绝' :
                                     exchange.status === 3 ? '已取消' : '已完成'
      exchange.get_type_text = () => exchange.type === 1 ? '购买' : '交换'
    })
  } catch (error) {
    console.error('获取交换列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取待处理数量
const fetchIncomingCount = async () => {
  try {
    const res = await getPendingCount()
    incomingCount.value = res.data.count
  } catch (error) {
    console.error('获取待处理数量失败:', error)
  }
}

// 类型改变
const handleTypeChange = () => {
  pagination.page = 1
  fetchList()
}

// 分页改变
const handlePageChange = (page) => {
  pagination.page = page
  fetchList()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handleSizeChange = (size) => {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

// 获取状态标签类型
const getStatusType = (status) => {
  const map = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'info',
    4: ''
  }
  return map[status] || ''
}

const getStatusText = (status) => {
  const map = {
    0: '待处理',
    1: '已接受',
    2: '已拒绝',
    3: '已取消',
    4: '已完成'
  }
  return map[status] || '未知'
}

const getTypeText = (type) => {
  const map = {
    1: '购买',
    2: '交换'
  }
  return map[type] || '未知'
}

// 获取状态提示文本
const getStatusTip = (exchange) => {
  const isInitiator = exchange.initiator_id === currentUserId.value
  const tips = {
    0: isInitiator ? '等待对方处理您的请求' : '您有一个新的交换请求待处理',
    1: '双方已达成一致，等待完成交换',
    2: '交换请求已被拒绝',
    3: '交换请求已取消',
    4: '交换已完成！'
  }
  return tips[exchange.status] || ''
}

// 获取状态步骤进度
const getStatusStep = (status) => {
  const stepMap = {
    0: 1, // 待处理 - 显示到"等待处理"
    1: 2, // 已接受 - 显示到"已接受"
    2: 2, // 已拒绝 - 显示到"被拒绝"
    3: 1, // 已取消 - 退回到"等待处理"之前
    4: 3  // 已完成 - 显示到"完成交换"
  }
  return stepMap[status] || 0
}

// 获取对方用户
const getOtherUser = (exchange) => {
  if (exchange.initiator_id === currentUserId.value) {
    return exchange.target
  }
  return exchange.initiator
}

// 接受交换
const handleAccept = async (exchange) => {
  try {
    await ElMessageBox.confirm('确认接受此交换请求？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'info'
    })

    await handleExchange({
      exchange_id: exchange.id,
      action: 'accept'
    })

    ElMessage.success('已接受交换请求')
    fetchList()
    fetchIncomingCount()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('接受交换失败:', error)
    }
  }
}

// 拒绝交换
const handleReject = (exchange) => {
  rejectForm.exchangeId = exchange.id
  rejectForm.reason = ''
  rejectDialogVisible.value = true
}

// 确认拒绝
const confirmReject = async () => {
  try {
    await handleExchange({
      exchange_id: rejectForm.exchangeId,
      action: 'reject',
      reject_reason: rejectForm.reason
    })

    ElMessage.success('已拒绝交换请求')
    rejectDialogVisible.value = false
    fetchList()
    fetchIncomingCount()
  } catch (error) {
    console.error('拒绝交换失败:', error)
  }
}

// 取消交换
const handleCancel = async (exchange) => {
  try {
    await ElMessageBox.confirm('确认取消此交换请求？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await handleExchange({
      exchange_id: exchange.id,
      action: 'cancel'
    })

    ElMessage.success('已取消交换请求')
    fetchList()
    fetchIncomingCount()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消交换失败:', error)
    }
  }
}

// 完成交换
const handleComplete = async (exchange) => {
  try {
    // 添加调试信息
    console.log('=== 完成交换调试信息 ===')
    console.log('交换ID:', exchange.id)
    console.log('商品ID:', exchange.goods_id)
    console.log('商品信息:', exchange.goods)
    console.log('商品类型:', exchange.goods?.type, '(1=售卖, 2=交换, 3=均可)')
    console.log('交换类型:', exchange.type, '(1=购买, 2=交换)')
    console.log('========================')

    await ElMessageBox.confirm('确认完成此交换？完成后商品状态将变更为已交换。', '提示', {
      confirmButtonText: '确认完成',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await handleExchange({
      exchange_id: exchange.id,
      action: 'complete'
    })

    ElMessage.success('交换已完成')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('完成交换失败:', error)
    }
  }
}

// 查看详情
const viewDetail = (exchange) => {
  ElMessage.info('详情页开发中')
}

// 跳转到商品详情
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

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  fetchList()
  fetchIncomingCount()
})
</script>

<style scoped>
.exchange-list {
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

.type-radio-group {
  width: 100%;
  display: flex;
}

.type-radio-group :deep(.el-radio-button) {
  flex: 1;
}

.type-radio-group :deep(.el-radio-button__inner) {
  width: 100%;
  border-radius: 10px;
  padding: 12px 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-weight: 500;
  border: 2px solid #dcdfe6;
  background: white;
  color: #606266;
  transition: all 0.3s;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
}

.type-radio-group :deep(.el-radio-button__inner:hover) {
  border-color: #409eff;
  color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.type-radio-group :deep(.el-radio-button.is-active .el-radio-button__inner) {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-color: #409eff;
  color: white;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.filter-select {
  width: 100%;
}

.filter-select :deep(.el-input__wrapper) {
  border-radius: 10px;
  transition: all 0.3s;
  height: 40px;
  font-size: 15px;
}

.filter-select:hover :deep(.el-input__wrapper) {
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

/* 交换网格 */
.exchange-wrapper {
  min-height: 400px;
  margin-top: 10px;
}

.exchange-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(450px, 1fr));
  gap: 24px;
  margin-bottom: 20px;
}

.exchange-item {
  height: 100%;
}

.exchange-card {
  height: 100%;
  border: none;
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  min-height: 650px;
}

.exchange-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.card-header {
  padding: 16px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
  border-bottom: 2px solid rgba(64, 158, 255, 0.1);
  min-height: 90px;
}

.status-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.status-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  background: white;
  border-radius: 10px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.status-tip .el-icon {
  color: #409eff;
  font-size: 16px;
}

.status-progress {
  padding: 20px 20px 0;
  background: white;
}

.status-progress :deep(.el-steps) {
  margin-bottom: 0;
}

.status-progress :deep(.el-step__head) {
  font-size: 14px;
}

.status-progress :deep(.el-step__title) {
  font-size: 13px;
  font-weight: 500;
}

.status-progress :deep(.el-step__icon) {
  width: 32px;
  height: 32px;
}

.status-progress :deep(.el-step.is-process .el-step__icon) {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-color: #409eff;
}

.status-progress :deep(.el-step.is-finish .el-step__icon) {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  border-color: #67c23a;
}

.status-progress :deep(.el-step__line) {
  background-color: #e8e8e8;
}

.status-progress :deep(.el-step.is-finish .el-step__line) {
  background-color: #67c23a;
}

.goods-images-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
  background: #fafbfc;
  flex: 1;
}

.goods-image-box {
  flex: 1;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  min-height: 260px;
}

.goods-image-box:hover {
  transform: scale(1.02);
}

.image-label {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 6px 12px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: white;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 10px;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.25);
}

.image-label.my-goods {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  box-shadow: 0 2px 8px rgba(103, 194, 58, 0.25);
}

.image-wrapper {
  width: 100%;
  height: 160px;
  border-radius: 12px;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.goods-image-box:hover .image-wrapper {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.goods-image {
  width: 100%;
  height: 100%;
}

.no-image {
  display: flex;
  align-items: center;
  justify-content: center;
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

.goods-title {
  margin-top: 10px;
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: center;
  height: 22px;
  line-height: 22px;
}

.goods-price {
  margin-top: 6px;
  font-size: 22px;
  font-weight: bold;
  color: #f56c6c;
  text-align: center;
  height: 28px;
  line-height: 28px;
}

.exchange-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  color: #409eff;
  flex-shrink: 0;
  width: 40px;
}

.user-section {
  padding: 16px 20px;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  min-height: 70px;
  display: flex;
  align-items: center;
}

.user-info-box {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-detail {
  flex: 1;
}

.user-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 2px;
}

.user-name {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.message-section {
  padding: 12px 20px;
  border-bottom: 1px solid #f0f0f0;
  min-height: 60px;
}

.message-box {
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  border-left: 3px solid #409eff;
}

.message-box.reject {
  background: #fef0f0;
  border-left-color: #f56c6c;
}

.message-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 6px;
}

.message-box p {
  margin: 0;
  font-size: 14px;
  color: #303133;
  line-height: 1.6;
}

.actions-section {
  padding: 16px 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  min-height: 80px;
  align-content: center;
  border-top: 1px solid #f0f0f0;
}

.action-btn {
  flex: 1;
  min-width: calc(50% - 5px);
  height: 48px;
  border-radius: 10px;
  font-weight: 500;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);
}

.action-btn.accept {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  border: none;
}

.action-btn.reject {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  border: none;
}

.action-btn.complete {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border: none;
}

.action-btn.cancel {
  background: linear-gradient(135deg, #909399 0%, #b1b3b8 100%);
  border: none;
}

.waiting-tip {
  flex: 1;
  min-width: calc(50% - 5px);
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #fff7e6 0%, #ffe8cc 100%);
  border-radius: 10px;
  color: #e6a23c;
  font-weight: 500;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(230, 162, 60, 0.15);
}

.waiting-tip .el-icon {
  font-size: 16px;
}

.action-tip {
  flex: 1;
  min-width: calc(50% - 5px);
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #e6f7ff 0%, #cce8ff 100%);
  border-radius: 10px;
  color: #409eff;
  font-weight: 500;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
}

.action-tip .el-icon {
  font-size: 16px;
}

.status-finished {
  flex: 1;
  min-width: calc(50% - 5px);
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
  border-radius: 10px;
  color: #606266;
  font-weight: 500;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.status-finished .el-icon {
  font-size: 18px;
}

.status-finished .status-rejected {
  color: #f56c6c;
}

.status-finished .status-cancelled {
  color: #909399;
}

.status-completed {
  flex: 1;
  min-width: calc(50% - 5px);
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f0f9ff 0%, #dcfce7 100%);
  border-radius: 10px;
  color: #67c23a;
  font-weight: 600;
  font-size: 15px;
  box-shadow: 0 2px 8px rgba(103, 194, 58, 0.2);
}

.status-completed .status-success {
  font-size: 20px;
}

.time-section {
  padding: 12px 20px;
  background: #fafbfc;
  border-top: 1px solid #f0f0f0;
  font-size: 13px;
  color: #909399;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 30px 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .exchange-grid {
    grid-template-columns: 1fr;
    gap: 20px;
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
  }

  .filter-title {
    font-size: 18px;
  }

  .type-radio-group :deep(.el-radio-button__inner) {
    padding: 10px 12px;
    font-size: 14px;
  }

  .goods-images-section {
    flex-direction: column;
  }

  .exchange-icon {
    transform: rotate(90deg);
    width: 100%;
    margin: 8px 0;
  }

  .image-wrapper {
    height: 160px;
  }

  .action-btn {
    min-width: 100%;
  }
}
</style>
