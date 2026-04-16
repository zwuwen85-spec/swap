import request from './request'

// 获取分类列表
export function getCategories() {
  return request({
    url: '/category/list',
    method: 'get'
  })
}
