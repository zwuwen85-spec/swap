import request from './request'

// 创建评论
export function createComment(data) {
  return request({
    url: '/comment/create',
    method: 'post',
    data
  })
}

// 获取商品评论列表
export function getCommentList(params) {
  return request({
    url: '/comment/list',
    method: 'get',
    params
  })
}

// 删除评论
export function deleteComment(commentId) {
  return request({
    url: `/comment/${commentId}`,
    method: 'delete'
  })
}

// 获取商品评分
export function getGoodsRating(goodsId) {
  return request({
    url: '/comment/rating',
    method: 'get',
    params: { goods_id: goodsId }
  })
}

// 获取我的评论
export function getMyComments(params) {
  return request({
    url: '/comment/my',
    method: 'get',
    params
  })
}

// 获取收到的评论
export function getReceivedComments(params) {
  return request({
    url: '/comment/received',
    method: 'get',
    params
  })
}
