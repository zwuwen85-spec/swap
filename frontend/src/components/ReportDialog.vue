<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="500px"
    @close="handleClose"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
      <el-form-item label="举报类型" prop="reason">
        <el-select v-model="form.reason" placeholder="请选择举报原因">
          <el-option
            v-for="reason in reasons"
            :key="reason"
            :label="reason"
            :value="reason"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="详细说明" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="5"
          maxlength="1000"
          show-word-limit
          placeholder="请详细描述举报原因，提供相关证据..."
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">
        提交举报
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { createReport } from '@/api/report'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  targetType: {
    type: Number,
    required: true
  },
  targetId: {
    type: Number,
    required: true
  },
  targetName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(false)
const loading = ref(false)
const formRef = ref(null)

const form = ref({
  reason: '',
  description: ''
})

const rules = {
  reason: [
    { required: true, message: '请选择举报原因', trigger: 'change' }
  ]
}

// 举报原因选项
const reasons = computed(() => {
  const reasonsMap = {
    1: ['虚假信息', '违规商品', '价格欺诈', '图片不符', '其他'],
    2: ['恶意行为', '欺诈行为', '骚扰用户', '虚假身份', '其他'],
    3: ['违法违规', '恶意攻击', '广告 spam', '不实信息', '其他']
  }
  return reasonsMap[props.targetType] || []
})

// 对话框标题
const title = computed(() => {
  const typeMap = {
    1: '举报商品',
    2: '举报用户',
    3: '举报评论'
  }
  return typeMap[props.targetType] || '举报'
})

// 监听modelValue变化
watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

// 提交举报
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await createReport({
        target_type: props.targetType,
        target_id: props.targetId,
        reason: form.value.reason,
        description: form.value.description
      })

      ElMessage.success('举报成功，我们会尽快处理')
      emit('success')
      handleClose()
    } catch (error) {
      console.error('举报失败:', error)
      ElMessage.error(error.response?.data?.message || '举报失败')
    } finally {
      loading.value = false
    }
  })
}

// 关闭对话框
const handleClose = () => {
  form.value = {
    reason: '',
    description: ''
  }
  formRef.value?.resetFields()
  visible.value = false
}
</script>

<style scoped>
.el-select {
  width: 100%;
}
</style>
