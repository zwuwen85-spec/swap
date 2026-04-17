import request from './request'

// 获取会话列表
export function getConversations() {
  return request({
    url: '/message/conversations',
    method: 'get'
  })
}

// 获取消息列表
export function getMessageList(userId) {
  return request({
    url: '/message/list',
    method: 'get',
    params: { user_id: userId }
  })
}

// 发送消息
export function sendMessage(data) {
  return request({
    url: '/message/send',
    method: 'post',
    data
  })
}

// 获取未读消息数量
export function getUnreadCount() {
  return request({
    url: '/message/unread-count',
    method: 'get'
  })
}

// 检查对方是否在线
export function checkOnline(userId) {
  return request({
    url: '/message/online-status',
    method: 'get',
    params: { user_id: userId }
  })
}
