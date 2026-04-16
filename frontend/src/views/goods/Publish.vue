<template>
  <div class="publish">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1 class="logo" @click="goHome">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <router-link to="/publish" :class="{ active: !isEdit }">{{ isEdit ? '编辑商品' : '发布商品' }}</router-link>
            <router-link v-if="!isEdit" to="/my-goods">我的发布</router-link>
            <router-link to="/profile">个人中心</router-link>
            <a @click="handleLogout">退出</a>
          </div>
        </div>
      </el-header>

      <el-main>
        <el-card class="publish-card">
          <h2>{{ isEdit ? '编辑商品' : '发布闲置物品' }}</h2>
          <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
            <!-- 标题 -->
            <el-form-item label="商品标题" prop="title">
              <el-input
                v-model="form.title"
                placeholder="请输入商品标题（1-100字符）"
                maxlength="100"
                show-word-limit
              />
            </el-form-item>

            <!-- 描述 -->
            <el-form-item label="商品描述">
              <el-input
                v-model="form.description"
                type="textarea"
                :rows="4"
                placeholder="请详细描述商品的品牌、型号、购买时间、使用情况等"
                maxlength="1000"
                show-word-limit
              />
            </el-form-item>

            <!-- 图片上传 -->
            <el-form-item label="商品图片" prop="images">
              <el-upload
                v-model:file-list="fileList"
                action="/api/v1/upload/image"
                :headers="{ Authorization: `Bearer ${userStore.token}` }"
                name="file"
                list-type="picture-card"
                :on-success="handleUploadSuccess"
                :on-remove="handleRemove"
                :before-upload="beforeUpload"
                :limit="9"
                accept="image/jpeg,image/jpg,image/png,image/gif"
              >
                <el-icon><Plus /></el-icon>
                <template #tip>
                  <div class="upload-tip">
                    最多上传9张图片，支持jpg、png、gif格式，每张不超过5MB
                  </div>
                </template>
              </el-upload>
            </el-form-item>

            <!-- 分类 -->
            <el-form-item label="商品分类" prop="category_id">
              <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%">
                <el-option
                  v-for="category in categories"
                  :key="category.id"
                  :label="category.name"
                  :value="category.id"
                />
              </el-select>
            </el-form-item>

            <!-- 交易类型 -->
            <el-form-item label="交易类型" prop="type">
              <el-radio-group v-model="form.type">
                <el-radio :label="1">售卖</el-radio>
                <el-radio :label="2">交换</el-radio>
                <el-radio :label="3">均可</el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 价格 -->
            <el-form-item
              v-if="form.type === 1 || form.type === 3"
              label="售卖价格"
              prop="price"
            >
              <el-input-number
                v-model="form.price"
                :min="0.01"
                :max="999999"
                :precision="2"
                :step="0.01"
                style="width: 200px"
              />
              <span style="margin-left: 10px">元</span>
            </el-form-item>

            <!-- 原价 -->
            <el-form-item v-if="form.type === 1 || form.type === 3" label="原价">
              <el-input-number
                v-model="form.original_price"
                :min="0.01"
                :max="999999"
                :precision="2"
                :step="0.01"
                placeholder="可选填"
                style="width: 200px"
              />
              <span style="margin-left: 10px">元</span>
            </el-form-item>

            <!-- 成色 -->
            <el-form-item label="商品成色" prop="condition">
              <el-radio-group v-model="form.condition">
                <el-radio :label="1">全新</el-radio>
                <el-radio :label="2">九成新</el-radio>
                <el-radio :label="3">八成新</el-radio>
                <el-radio :label="4">七成新</el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 交易地点 -->
            <el-form-item label="交易地点">
              <el-input
                v-model="form.location"
                placeholder="例如：图书馆、食堂、宿舍楼等"
                maxlength="100"
              />
            </el-form-item>

            <!-- 标签 -->
            <el-form-item label="商品标签">
              <el-input
                v-model="form.tags"
                placeholder="多个标签用逗号分隔，例如：苹果,笔记本,电脑"
                maxlength="200"
              />
            </el-form-item>

            <!-- 按钮 -->
            <el-form-item>
              <el-button type="primary" @click="handleSubmit" :loading="submitting">
                {{ isEdit ? '保存修改' : '立即发布' }}
              </el-button>
              <el-button @click="handleReset">重置</el-button>
              <el-button @click="goBack">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { createGoods, getGoodsDetail, updateGoods, getCategories, uploadGoodsImage } from '@/api/goods'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 是否为编辑模式
const isEdit = computed(() => !!route.query.edit)
const editGoodsId = computed(() => route.query.edit)

const formRef = ref(null)
const submitting = ref(false)
const fileList = ref([])
const categories = ref([])

// 表单数据
const form = reactive({
  title: '',
  description: '',
  category_id: null,
  type: 1,
  price: 0,
  original_price: null,
  images: [],
  condition: 1,
  location: '',
  tags: ''
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入商品标题', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  type: [
    { required: true, message: '请选择交易类型', trigger: 'change' }
  ],
  condition: [
    { required: true, message: '请选择商品成色', trigger: 'change' }
  ],
  price: [
    {
      validator: (rule, value, callback) => {
        if ((form.type === 1 || form.type === 3) && (!value || value <= 0)) {
          callback(new Error('请输入售卖价格'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  images: [
    {
      validator: (rule, value, callback) => {
        if (form.images.length === 0) {
          callback(new Error('请至少上传一张商品图片'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
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

// 获取商品详情（编辑模式）
const fetchGoodsDetail = async (id) => {
  try {
    console.log('开始获取商品详情，ID:', id)
    const res = await getGoodsDetail(id)
    console.log('API返回:', res)

    // 关键修复：从 res.data.goods 中获取商品数据
    const goods = res.data.goods
    console.log('商品数据:', goods)

    // 填充表单
    Object.assign(form, {
      title: goods.title || '',
      description: goods.description || '',
      category_id: goods.category_id || null,
      type: goods.type || 1,
      price: goods.price || 0,
      original_price: goods.original_price || null,
      images: [],
      condition: goods.condition || 1,
      location: goods.location || '',
      tags: goods.tags || ''
    })

    console.log('填充后的表单:', form)

    // 处理图片
    if (goods.images) {
      let imageList = []
      if (typeof goods.images === 'string') {
        try {
          imageList = JSON.parse(goods.images)
        } catch {
          imageList = []
        }
      } else {
        imageList = goods.images
      }

      form.images = imageList

      // 设置文件列表用于预览
      fileList.value = imageList.map((url, index) => ({
        name: `image_${index}`,
        url: url,
        uid: Date.now() + index
      }))

      console.log('图片列表:', imageList)
    }

    ElMessage.success('商品信息已加载')
  } catch (error) {
    console.error('获取商品详情失败:', error)
    ElMessage.error('获取商品详情失败: ' + (error.message || '未知错误'))
  }
}

// 上传前验证
const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response, file) => {
  if (response.code === 0) {
    // 单图上传，每次添加一个URL到数组
    const url = response.data.url
    if (url && !form.images.includes(url)) {
      form.images.push(url)
    }
    // 触发表单验证
    formRef.value.validateField('images')
  } else {
    ElMessage.error(response.message || '上传失败')
    // 从列表中移除失败的文件
    const index = fileList.value.indexOf(file)
    if (index > -1) {
      fileList.value.splice(index, 1)
    }
  }
}

// 移除图片
const handleRemove = (file) => {
  // 从响应中获取URL
  if (file.response && file.response.data && file.response.data.url) {
    const url = file.response.data.url
    const index = form.images.indexOf(url)
    if (index > -1) {
      form.images.splice(index, 1)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const formData = {
        title: form.title,
        description: form.description,
        category_id: form.category_id,
        type: form.type,
        price: form.price,
        original_price: form.original_price,
        images: form.images,
        condition: form.condition,
        location: form.location,
        tags: form.tags
      }

      let res
      if (isEdit.value) {
        // 编辑模式 - 将ID转换为数字
        formData.id = parseInt(editGoodsId.value)
        console.log('更新商品，ID:', formData.id)
        console.log('表单数据:', formData)
        res = await updateGoods(formData)
        ElMessage.success(res.message || '更新成功')
      } else {
        // 新建模式
        res = await createGoods(formData)
        ElMessage.success(res.message || '发布成功')
      }

      setTimeout(() => {
        router.push('/my-goods')
      }, 1500)
    } catch (error) {
      console.error(isEdit.value ? '更新失败:' : '发布失败:', error)
    } finally {
      submitting.value = false
    }
  })
}

// 重置表单
const handleReset = () => {
  formRef.value?.resetFields()
  fileList.value = []
  Object.assign(form, {
    title: '',
    description: '',
    category_id: null,
    type: 1,
    price: 0,
    original_price: null,
    images: [],
    condition: 1,
    location: '',
    tags: ''
  })
}

// 返回
const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/')
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(async () => {
  console.log('=== 发布页面加载 ===')
  console.log('完整路由:', route)
  console.log('路由参数:', route.query)
  console.log('是否编辑模式:', isEdit.value)
  console.log('商品ID:', editGoodsId.value, '类型:', typeof editGoodsId.value)

  await fetchCategories()

  // 如果是编辑模式，获取商品详情
  if (isEdit.value && editGoodsId.value) {
    console.log('检测到编辑模式，准备获取商品详情')
    await fetchGoodsDetail(editGoodsId.value)
  } else {
    console.log('新建模式')
  }
})
</script>

<style scoped>
.publish {
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
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.publish-card {
  padding: 30px;
}

.publish-card h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #303133;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 10px;
}

:deep(.el-upload-list__item) {
  width: 148px;
  height: 148px;
}

:deep(.el-upload--picture-card) {
  width: 148px;
  height: 148px;
}
</style>
