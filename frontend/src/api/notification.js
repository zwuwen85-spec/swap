import request from './request'

// 获取通知列表
export function getNotificationList(params) {
  return request({
    url: '/notification/list',
    method: 'get',
    params
  })
}

// 获取未读通知数量
export function getUnreadCount() {
  return request({
    url: '/notification/unread-count',
    method: 'get'
  })
}

// 标记单条通知为已读
export function markAsRead(notificationId) {
  return request({
    url: `/notification/read/${notificationId}`,
    method: 'put'
  })
}

// 标记所有通知为已读
export function markAllAsRead() {
  return request({
    url: '/notification/read-all',
    method: 'put'
  })
}

// 删除通知
export function deleteNotification(notificationId) {
  return request({
    url: `/notification/${notificationId}`,
    method: 'delete'
  })
}

// 清空已读通知
export function clearReadNotifications() {
  return request({
    url: '/notification/clear-read',
    method: 'delete'
  })
}
