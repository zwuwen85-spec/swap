import request from './request'

// 添加收藏
export function addFavorite(goodsId) {
  return request({
    url: '/favorite/add',
    method: 'post',
    data: { goods_id: Number(goodsId) }
  })
}

// 取消收藏
export function removeFavorite(goodsId) {
  return request({
    url: `/favorite/remove/${goodsId}`,
    method: 'delete'
  })
}

// 检查是否已收藏
export function checkFavorite(goodsId) {
  return request({
    url: '/favorite/check',
    method: 'get',
    params: { goods_id: Number(goodsId) }
  })
}

// 获取收藏列表
export function getFavoriteList(params) {
  return request({
    url: '/favorite/list',
    method: 'get',
    params
  })
}

// 获取收藏数量
export function getFavoriteCount() {
  return request({
    url: '/favorite/count',
    method: 'get'
  })
}
