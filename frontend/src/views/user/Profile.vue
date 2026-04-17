<template>
  <div class="profile">
    <el-container>
      <!-- 头部导航 -->
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/chat">
              <el-badge :value="userStore.unreadMessageCount" :hidden="userStore.unreadMessageCount === 0" :max="99" class="nav-badge">
                消息
              </el-badge>
            </router-link>
            <router-link to="/profile" class="active">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <el-main>
        <div class="profile-container">
          <!-- 用户信息卡片 -->
          <el-card class="user-card">
            <div class="user-info">
              <el-avatar :size="100" :src="userInfo.avatar">
                {{ userInfo.username?.[0] }}
              </el-avatar>
              <div class="user-detail">
                <h2>{{ userInfo.nickname || userInfo.username }}</h2>
                <p class="username">@{{ userInfo.username }}</p>
                <div class="user-meta">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <p class="meta-item">
                        <el-icon><School /></el-icon>
                        <span class="meta-label">学校：</span>{{ userInfo.school || '未设置' }}
                      </p>
                      <p class="meta-item">
                        <el-icon><Postcard /></el-icon>
                        <span class="meta-label">学号：</span>{{ userInfo.student_id || '未设置' }}
                      </p>
                      <p class="meta-item">
                        <el-icon><Reading /></el-icon>
                        <span class="meta-label">专业：</span>{{ userInfo.major || '未设置' }}
                      </p>
                      <p class="meta-item">
                        <el-icon><User /></el-icon>
                        <span class="meta-label">性别：</span>{{ getGenderText(userInfo.gender) }}
                      </p>
                      <p class="meta-item">
                        <el-icon><Phone /></el-icon>
                        <span class="meta-label">手机：</span>{{ userInfo.phone || '未设置' }}
                      </p>
                    </el-col>
                    <el-col :span="12">
                      <p class="meta-item">
                        <el-icon><Message /></el-icon>
                        <span class="meta-label">邮箱：</span>{{ userInfo.email || '未设置' }}
                      </p>
                      <p class="meta-item">
                        <el-icon><ChatDotRound /></el-icon>
                        <span class="meta-label">QQ：</span>{{ userInfo.qq || '未设置' }}
                      </p>
                      <p class="meta-item">
                        <el-icon><ChatDotRound /></el-icon>
                        <span class="meta-label">微信：</span>{{ userInfo.wechat || '未设置' }}
                      </p>
                    </el-col>
                  </el-row>
                </div>
                <el-button type="primary" size="small" @click="openEditDialog" class="edit-btn">
                  <el-icon><Edit /></el-icon>
                  编辑资料
                </el-button>
              </div>
            </div>
          </el-card>

          <!-- 统计数据 -->
          <el-row :gutter="20" class="stats-row">
            <el-col :span="6">
              <el-card class="stat-card" shadow="hover" @click="goTo('/my-goods')">
                <div class="stat-item">
                  <el-icon class="stat-icon"><ShoppingBag /></el-icon>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.goodsCount }}</div>
                    <div class="stat-label">我的发布</div>
                  </div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card" shadow="hover" @click="goTo('/favorites')">
                <div class="stat-item">
                  <el-icon class="stat-icon"><Star /></el-icon>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.favoriteCount }}</div>
                    <div class="stat-label">我的收藏</div>
                  </div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card" shadow="hover" @click="goTo('/exchange')">
                <div class="stat-item">
                  <el-icon class="stat-icon"><Switch /></el-icon>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.exchangeCount }}</div>
                    <div class="stat-label">我的交换</div>
                  </div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card class="stat-card" shadow="hover" @click="goTo('/my-comments')">
                <div class="stat-item">
                  <el-icon class="stat-icon"><ChatDotSquare /></el-icon>
                  <div class="stat-info">
                    <div class="stat-value">{{ stats.commentCount }}</div>
                    <div class="stat-label">我的评论</div>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>

          <!-- 快捷操作 -->
          <el-card class="actions-card">
            <template #header>
              <h3>快捷操作</h3>
            </template>
            <el-row :gutter="20">
              <el-col :span="6">
                <el-button type="primary" @click="goTo('/publish')" class="action-btn">
                  <el-icon><Plus /></el-icon>
                  发布商品
                </el-button>
              </el-col>
              <el-col :span="6">
                <el-button @click="goTo('/my-goods')" class="action-btn">
                  <el-icon><ShoppingBag /></el-icon>
                  我的发布
                </el-button>
              </el-col>
              <el-col :span="6">
                <el-button @click="goTo('/favorites')" class="action-btn">
                  <el-icon><Star /></el-icon>
                  我的收藏
                </el-button>
              </el-col>
              <el-col :span="6">
                <el-button @click="goTo('/exchange')" class="action-btn">
                  <el-icon><Switch /></el-icon>
                  我的交换
                </el-button>
              </el-col>
            </el-row>
          </el-card>

          <!-- 待处理事项 -->
          <el-card v-if="pendingCount > 0" class="pending-card" shadow="hover">
            <template #header>
              <div class="pending-header">
                <h3>待处理事项</h3>
                <el-badge :value="pendingCount" class="badge" />
              </div>
            </template>
            <div class="pending-item" @click="goTo('/exchange')">
              <el-icon class="pending-icon"><Warning /></el-icon>
              <span>您有 {{ pendingCount }} 个交换请求待处理</span>
              <el-icon class="arrow"><ArrowRight /></el-icon>
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>

    <!-- 编辑资料对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑资料" width="700px">
      <el-form ref="editFormRef" :model="editForm" :rules="formRules" label-width="90px">
        <el-divider content-position="left">头像</el-divider>
        <el-form-item label="头像">
          <div class="avatar-upload-container">
            <img v-if="editForm.avatar" :src="editForm.avatar" class="avatar-preview" />
            <div v-else class="avatar-placeholder">
              <el-icon class="avatar-placeholder-icon"><User /></el-icon>
              <span>暂无头像</span>
            </div>
            <el-upload
              class="avatar-uploader-btn"
              action="#"
              :show-file-list="false"
              :before-upload="handleAvatarUpload"
              :http-request="customUpload"
              accept="image/jpeg,image/jpg,image/png,image/gif,image/webp,image/bmp"
            >
              <el-button type="primary" size="small">
                <el-icon><Upload /></el-icon>
                选择图片
              </el-button>
            </el-upload>
          </div>
          <div class="avatar-tip">支持 JPG、PNG、GIF、WebP、BMP 等格式，大小不超过 5MB</div>
        </el-form-item>

        <el-divider content-position="left">基本信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="昵称">
              <el-input v-model="editForm.nickname" placeholder="请输入昵称" maxlength="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学校">
              <el-input v-model="editForm.school" placeholder="请输入学校名称" maxlength="100" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">学籍信息</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="学号">
              <el-input v-model="editForm.student_id" placeholder="请输入学号" maxlength="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="专业">
              <el-input v-model="editForm.major" placeholder="请输入专业" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别">
              <el-select v-model="editForm.gender" placeholder="请选择性别" style="width: 100%">
                <el-option label="未设置" :value="0" />
                <el-option label="男" :value="1" />
                <el-option label="女" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left">联系方式</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="手机号">
              <el-input v-model="editForm.phone" placeholder="请输入手机号" maxlength="11">
                <template #prefix>
                  <el-icon><Phone /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱">
              <el-input v-model="editForm.email" placeholder="请输入邮箱" maxlength="50">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="QQ">
              <el-input v-model="editForm.qq" placeholder="请输入QQ号" maxlength="15">
                <template #prefix>
                  <el-icon><ChatDotRound /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="微信">
              <el-input v-model="editForm.wechat" placeholder="请输入微信号" maxlength="50">
                <template #prefix>
                  <el-icon><ChatDotRound /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-alert
          title="提示"
          type="info"
          :closable="false"
          show-icon
          style="margin-top: 10px"
        >
          <template #default>
            <div>• 手机号、邮箱、QQ、微信号仅对已交换成功的用户可见</div>
            <div>• 建议完善联系方式，方便交易沟通</div>
          </template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSaveProfile" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getUserInfo, updateUserInfo } from '@/api/user'
import { getFavoriteCount } from '@/api/favorite'
import { getPendingCount, getExchangeList } from '@/api/exchange'
import { getMyGoods } from '@/api/goods'
import { getMyComments } from '@/api/comment'
import { ElMessage } from 'element-plus'
import {
  Edit,
  ShoppingBag,
  Star,
  Switch,
  ChatDotSquare,
  ChatDotRound,
  Plus,
  Warning,
  ArrowRight,
  School,
  Postcard,
  Reading,
  Phone,
  Message,
  User,
  Upload
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const userInfo = ref({
  username: '',
  nickname: '',
  avatar: '',
  school: '',
  student_id: '',
  major: '',
  gender: 0,
  phone: '',
  email: '',
  qq: '',
  wechat: ''
})

const stats = reactive({
  goodsCount: 0,
  favoriteCount: 0,
  exchangeCount: 0,
  commentCount: 0
})

const pendingCount = ref(0)
const showEditDialog = ref(false)
const saving = ref(false)
const editFormRef = ref(null)

// 邮箱验证器
const validateEmail = (rule, value, callback) => {
  if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
    callback(new Error('请输入正确的邮箱格式'))
  } else {
    callback()
  }
}

// 表单验证规则
const formRules = {
  email: [
    { validator: validateEmail, trigger: 'blur' }
  ]
}

const editForm = reactive({
  avatar: '',
  nickname: '',
  school: '',
  student_id: '',
  major: '',
  gender: 0,
  phone: '',
  email: '',
  qq: '',
  wechat: ''
})

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const res = await getUserInfo()
    userInfo.value = res.data

    // 更新 store 中的用户信息
    userStore.userInfo = res.data
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    // 获取收藏数量
    const favRes = await getFavoriteCount()
    stats.favoriteCount = favRes.data.count

    // 获取待处理交换数量
    const pendingRes = await getPendingCount()
    pendingCount.value = pendingRes.data.count

    // 获取商品数量
    try {
      const goodsRes = await getMyGoods({
        page: 1,
        page_size: 1
      })
      stats.goodsCount = goodsRes.data.total
    } catch (error) {
      console.error('获取商品数量失败:', error)
    }

    // 获取评论数量
    try {
      const commentRes = await getMyComments({
        page: 1,
        page_size: 1
      })
      stats.commentCount = commentRes.data.total
    } catch (error) {
      console.error('获取评论数量失败:', error)
    }

    // 获取交换总数
    try {
      const exchangeRes = await getExchangeList({
        type: 'all',
        page: 1,
        page_size: 1
      })
      stats.exchangeCount = exchangeRes.data.total
    } catch (error) {
      console.error('获取交换数量失败:', error)
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 打开编辑对话框
const openEditDialog = () => {
  editForm.avatar = userInfo.value.avatar || ''
  editForm.nickname = userInfo.value.nickname || ''
  editForm.school = userInfo.value.school || ''
  editForm.student_id = userInfo.value.student_id || ''
  editForm.major = userInfo.value.major || ''
  editForm.gender = userInfo.value.gender || 0
  editForm.phone = userInfo.value.phone || ''
  editForm.email = userInfo.value.email || ''
  editForm.qq = userInfo.value.qq || ''
  editForm.wechat = userInfo.value.wechat || ''
  showEditDialog.value = true
}

// 头像上传验证
const handleAvatarUpload = (file) => {
  // 支持的图片类型
  const allowedTypes = [
    'image/jpeg',
    'image/jpg',
    'image/png',
    'image/gif',
    'image/webp',
    'image/bmp'
  ]

  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('只支持 JPG、PNG、GIF、WebP、BMP 格式的图片!')
    return false
  }

  // 文件大小限制 5MB
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }

  return true
}

// 自定义上传请求
const customUpload = async ({ file }) => {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const token = localStorage.getItem('token')
    const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    const response = await fetch(`${baseURL}/user/avatar`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    })

    const result = await response.json()

    if (result.code === 0) {
      editForm.avatar = result.data.url
      ElMessage.success('头像上传成功')
    } else {
      ElMessage.error(result.message || '头像上传失败')
    }
  } catch (error) {
    console.error('头像上传失败:', error)
    ElMessage.error('头像上传失败，请重试')
  }
}

// 保存资料
const handleSaveProfile = async () => {
  try {
    // 验证表单
    const valid = await editFormRef.value?.validate().catch(() => false)
    if (!valid) {
      return
    }

    saving.value = true

    // 邮箱格式验证正则
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    // 过滤空字符串和无效格式，只发送有值的字段
    const updateData = Object.fromEntries(
      Object.entries({
        avatar: editForm.avatar,
        nickname: editForm.nickname,
        school: editForm.school,
        student_id: editForm.student_id,
        major: editForm.major,
        gender: editForm.gender,
        phone: editForm.phone,
        email: editForm.email,
        qq: editForm.qq,
        wechat: editForm.wechat
      }).filter(([key, value]) => {
        // 过滤空值
        if (value === '' || value === null || value === undefined) return false
        // 对于 email 字段，还要验证格式
        if (key === 'email' && !emailRegex.test(value)) return false
        return true
      })
    )

    await updateUserInfo(updateData)

    // 更新本地数据
    userInfo.value.avatar = editForm.avatar
    userInfo.value.nickname = editForm.nickname
    userInfo.value.school = editForm.school
    userInfo.value.student_id = editForm.student_id
    userInfo.value.major = editForm.major
    userInfo.value.gender = editForm.gender
    userInfo.value.phone = editForm.phone
    userInfo.value.email = editForm.email
    userInfo.value.qq = editForm.qq
    userInfo.value.wechat = editForm.wechat

    // 更新 store 中的用户信息
    if (userStore.userInfo) {
      userStore.userInfo.avatar = editForm.avatar
      userStore.userInfo.nickname = editForm.nickname
      userStore.userInfo.school = editForm.school
      userStore.userInfo.student_id = editForm.student_id
      userStore.userInfo.major = editForm.major
      userStore.userInfo.gender = editForm.gender
      userStore.userInfo.phone = editForm.phone
      userStore.userInfo.email = editForm.email
      userStore.userInfo.qq = editForm.qq
      userStore.userInfo.wechat = editForm.wechat
    }

    ElMessage.success('保存成功')
    showEditDialog.value = false
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 跳转
const goTo = (path) => {
  router.push(path)
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

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

onMounted(() => {
  fetchUserInfo()
  fetchStats()
})
</script>

<style scoped>
.profile {
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
}

.profile-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.user-info {
  display: flex;
  gap: 20px;
  align-items: center;
}

.user-detail {
  flex: 1;
}

.user-detail h2 {
  margin: 0 0 5px 0;
  font-size: 24px;
}

.username {
  margin: 0 0 10px 0;
  opacity: 0.9;
  font-size: 14px;
}

.user-meta {
  margin: 15px 0;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0 0 8px 0;
  font-size: 14px;
  opacity: 0.9;
}

.meta-item .el-icon {
  font-size: 16px;
}

.meta-label {
  font-weight: 500;
  opacity: 0.95;
}

.edit-btn {
  margin-top: 10px;
}

.stats-row {
  margin-bottom: 0;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  font-size: 32px;
  color: #409eff;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.actions-card {
  margin-bottom: 0;
}

.actions-card h3 {
  margin: 0;
  font-size: 16px;
}

.action-btn {
  width: 100%;
  margin-bottom: 10px;
}

.pending-card {
  background-color: #fff7e6;
  border: 1px solid #ffd591;
  cursor: pointer;
}

.pending-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pending-header h3 {
  margin: 0;
  font-size: 16px;
  color: #fa8c16;
}

.pending-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  background-color: #fff;
  border-radius: 4px;
  color: #fa8c16;
}

.pending-icon {
  font-size: 20px;
}

.arrow {
  margin-left: auto;
}

/* 头像上传 */
.avatar-upload-container {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 8px;
  object-fit: cover;
  border: 2px solid #e0e0e0;
}

.avatar-placeholder {
  width: 100px;
  height: 100px;
  border-radius: 8px;
  border: 2px dashed #d9d9d9;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #909399;
}

.avatar-placeholder-icon {
  font-size: 32px;
  margin-bottom: 5px;
}

.avatar-placeholder span {
  font-size: 12px;
}

.avatar-uploader-btn {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.avatar-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}
</style>
