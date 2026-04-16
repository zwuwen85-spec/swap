<template>
  <div class="comment-list">
    <el-card class="comment-card">
      <template #header>
        <div class="card-header">
          <span class="title">商品评价</span>
          <div class="rating-summary">
            <span class="avg-rating">{{ avgRating.toFixed(1) }}</span>
            <el-rate v-model="displayRating" disabled show-score />
            <span class="count">共 {{ total }} 条评价</span>
          </div>
        </div>
      </template>

      <!-- 评论列表 -->
      <div v-loading="loading">
        <el-empty v-if="comments.length === 0 && !loading" description="暂无评价" />

        <div v-for="comment in comments" :key="comment.id" class="comment-item">
          <!-- 主评论 -->
          <div class="main-comment">
            <el-avatar :size="40" :src="comment.user?.avatar">
              {{ comment.user?.username?.[0] }}
            </el-avatar>
            <div class="comment-content">
              <div class="comment-header">
                <span class="username">{{ comment.user?.nickname || comment.user?.username }}</span>
                <el-rate v-model="comment.rating" disabled size="small" />
                <span class="time">{{ formatTime(comment.create_time) }}</span>
              </div>
              <div class="comment-text">{{ comment.content }}</div>
              <div class="comment-actions">
                <el-button
                  v-if="userStore.isLogin && userStore.userId !== comment.user_id"
                  type="primary"
                  link
                  size="small"
                  @click="showReplyInput(comment.id)"
                >
                  回复
                </el-button>
                <el-button
                  v-if="userStore.userId === comment.user_id"
                  type="danger"
                  link
                  size="small"
                  @click="handleDelete(comment.id)"
                >
                  删除
                </el-button>
              </div>
            </div>
          </div>

          <!-- 回复列表 -->
          <div v-if="comment.replies && comment.replies.length > 0" class="reply-list">
            <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
              <el-avatar :size="30" :src="reply.user?.avatar">
                {{ reply.user?.username?.[0] }}
              </el-avatar>
              <div class="reply-content">
                <div class="reply-header">
                  <span class="username">{{ reply.user?.nickname || reply.user?.username }}</span>
                  <span class="time">{{ formatTime(reply.create_time) }}</span>
                </div>
                <div class="reply-text">{{ reply.content }}</div>
                <div class="reply-actions">
                  <el-button
                    v-if="userStore.userId === reply.user_id"
                    type="danger"
                    link
                    size="small"
                    @click="handleDelete(reply.id)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
            </div>
          </div>

          <!-- 回复输入框 -->
          <div v-if="replyInputVisible[comment.id]" class="reply-input">
            <el-input
              v-model="replyContent[comment.id]"
              type="textarea"
              :rows="3"
              placeholder="输入回复内容..."
              maxlength="500"
              show-word-limit
            />
            <div class="reply-buttons">
              <el-button size="small" @click="cancelReply(comment.id)">取消</el-button>
              <el-button type="primary" size="small" @click="submitReply(comment.id)">发表</el-button>
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

    <!-- 发表评论对话框 -->
    <el-dialog v-model="dialogVisible" title="发表评价" width="500px">
      <el-form :model="commentForm" label-width="80px">
        <el-form-item label="评分">
          <el-rate v-model="commentForm.rating" show-text :texts="['极差', '失望', '一般', '满意', '惊喜']" />
        </el-form-item>
        <el-form-item label="评价内容">
          <el-input
            v-model="commentForm.content"
            type="textarea"
            :rows="5"
            maxlength="500"
            show-word-limit
            placeholder="分享您的使用体验..."
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitComment" :loading="submitting">发表</el-button>
      </template>
    </el-dialog>

    <!-- 发表评论按钮 -->
    <el-button type="primary" class="fab-comment" @click="openDialog">
      <el-icon><ChatDotSquare /></el-icon>
      写评价
    </el-button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { getCommentList, getGoodsRating, createComment, deleteComment } from '@/api/comment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatDotSquare } from '@element-plus/icons-vue'

const props = defineProps({
  goodsId: {
    type: Number,
    required: true
  }
})

const userStore = useUserStore()

const loading = ref(false)
const comments = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const avgRating = ref(0)
const displayRating = ref(0)

const dialogVisible = ref(false)
const submitting = ref(false)
const commentForm = ref({
  rating: 5,
  content: ''
})

const replyInputVisible = ref({})
const replyContent = ref({})

// 获取评论列表
const fetchComments = async () => {
  loading.value = true
  try {
    const res = await getCommentList({
      goods_id: props.goodsId,
      page: currentPage.value,
      page_size: pageSize.value
    })
    comments.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取评论列表失败:', error)
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
  }
}

// 获取评分
const fetchRating = async () => {
  try {
    const res = await getGoodsRating(props.goodsId)
    avgRating.value = res.data.avg_rating || 0
    displayRating.value = Math.round(res.data.avg_rating || 0)
  } catch (error) {
    console.error('获取评分失败:', error)
  }
}

// 打开发表评论对话框
const openDialog = () => {
  if (!userStore.isLogin) {
    ElMessage.warning('请先登录')
    return
  }
  dialogVisible.value = true
}

// 提交评论
const submitComment = async () => {
  if (!commentForm.value.content.trim()) {
    ElMessage.warning('请输入评价内容')
    return
  }

  submitting.value = true
  try {
    await createComment({
      goods_id: props.goodsId,
      content: commentForm.value.content,
      rating: commentForm.value.rating
    })

    ElMessage.success('评价成功')
    dialogVisible.value = false
    commentForm.value = { rating: 5, content: '' }

    // 刷新评论列表
    await fetchComments()
    await fetchRating()
  } catch (error) {
    console.error('发表评论失败:', error)
    ElMessage.error('发表评论失败')
  } finally {
    submitting.value = false
  }
}

// 显示回复输入框
const showReplyInput = (commentId) => {
  replyInputVisible.value[commentId] = true
  replyContent.value[commentId] = ''
}

// 取消回复
const cancelReply = (commentId) => {
  replyInputVisible.value[commentId] = false
  replyContent.value[commentId] = ''
}

// 提交回复
const submitReply = async (commentId) => {
  if (!replyContent.value[commentId].trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  try {
    await createComment({
      goods_id: props.goodsId,
      content: replyContent.value[commentId],
      rating: 5, // 回复不需要评分，默认5
      parent_id: commentId
    })

    ElMessage.success('回复成功')
    replyInputVisible.value[commentId] = false
    replyContent.value[commentId] = ''

    // 刷新评论列表
    await fetchComments()
  } catch (error) {
    console.error('回复失败:', error)
    ElMessage.error('回复失败')
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

    // 刷新评论列表
    await fetchComments()
    await fetchRating()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除评论失败:', error)
      ElMessage.error('删除评论失败')
    }
  }
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
  fetchComments()
  fetchRating()
})
</script>

<style scoped>
.comment-list {
  margin-top: 20px;
  position: relative;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.rating-summary {
  display: flex;
  align-items: center;
  gap: 15px;
}

.avg-rating {
  font-size: 32px;
  font-weight: bold;
  color: #f56c6c;
}

.count {
  color: #909399;
  font-size: 14px;
}

.comment-item {
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.main-comment {
  display: flex;
  gap: 15px;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.username {
  font-weight: 500;
  color: #303133;
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

.reply-list {
  margin-top: 15px;
  margin-left: 55px;
  padding-left: 20px;
  border-left: 2px solid #ebeef5;
}

.reply-item {
  display: flex;
  gap: 10px;
  padding: 10px 0;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 5px;
}

.reply-text {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 5px;
}

.reply-actions {
  display: flex;
  gap: 10px;
}

.reply-input {
  margin-top: 15px;
  margin-left: 55px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.reply-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.fab-comment {
  position: fixed;
  right: 50px;
  bottom: 100px;
  z-index: 100;
}
</style>
