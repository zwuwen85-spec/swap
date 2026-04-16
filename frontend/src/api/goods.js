import request from './request'

// 获取商品列表
export function getGoodsList(params) {
  return request({
    url: '/goods/list',
    method: 'get',
    params
  })
}

// 获取商品详情
export function getGoodsDetail(id) {
  return request({
    url: '/goods/detail',
    method: 'get',
    params: { id }
  })
}

// 发布商品
export function createGoods(data) {
  return request({
    url: '/goods/create',
    method: 'post',
    data
  })
}

// 更新商品
export function updateGoods(data) {
  return request({
    url: '/goods/update',
    method: 'put',
    data
  })
}

// 删除商品
export function deleteGoods(id) {
  return request({
    url: '/goods/delete',
    method: 'delete',
    data: { id }
  })
}

// 我的发布
export function getMyGoods(params) {
  return request({
    url: '/goods/my',
    method: 'get',
    params
  })
}

// 搜索商品
export function searchGoods(params) {
  return request({
    url: '/goods/search',
    method: 'get',
    params
  })
}

// 上传商品图片
export function uploadGoodsImage(file) {
  const formData = new FormData()
  formData.append('file', file)

  return request({
    url: '/goods/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取分类列表
export function getCategories() {
  return request({
    url: '/category/list',
    method: 'get'
  })
}

// 根据用户ID获取商品列表
export function getGoodsByUserId(userId, params) {
  return request({
    url: `/user/${userId}/goods`,
    method: 'get',
    params
  })
}
