import request from './request'

// 发起交换请求
export function createExchange(data) {
  return request({
    url: '/exchange/create',
    method: 'post',
    data
  })
}

// 获取交换列表
export function getExchangeList(params) {
  return request({
    url: '/exchange/list',
    method: 'get',
    params
  })
}

// 获取交换详情
export function getExchangeDetail(id) {
  return request({
    url: '/exchange/detail',
    method: 'get',
    params: { id }
  })
}

// 处理交换请求（接受/拒绝/取消/完成）
export function handleExchange(data) {
  return request({
    url: '/exchange/handle',
    method: 'post',
    data
  })
}

// 获取待处理交换数量
export function getPendingCount() {
  return request({
    url: '/exchange/pending-count',
    method: 'get'
  })
}
