import request from './request'

// 创建举报
export function createReport(data) {
  return request({
    url: '/report/create',
    method: 'post',
    data
  })
}

// 获取我的举报列表
export function getMyReports(params) {
  return request({
    url: '/report/my',
    method: 'get',
    params
  })
}

// 获取举报列表（管理员）
export function getReportList(params) {
  return request({
    url: '/report/list',
    method: 'get',
    params
  })
}

// 获取举报详情（管理员）
export function getReportDetail(reportId) {
  return request({
    url: `/report/detail/${reportId}`,
    method: 'get'
  })
}

// 处理举报（管理员）
export function handleReport(reportId, data) {
  return request({
    url: `/report/handle/${reportId}`,
    method: 'put',
    data
  })
}

// 获取待处理举报数量（管理员）
export function getPendingCount() {
  return request({
    url: '/report/pending-count',
    method: 'get'
  })
}

// 撤销举报
export function cancelReport(reportId) {
  return request({
    url: `/report/${reportId}`,
    method: 'delete'
  })
}
