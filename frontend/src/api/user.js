import request from './request'

// 用户注册
export function register(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

// 用户登录
export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

// 更新用户信息
export function updateUserInfo(data) {
  return request({
    url: '/user/info',
    method: 'put',
    data
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/user/change-password',
    method: 'post',
    data
  })
}

// 获取指定用户信息
export function getUserInfoById(userId) {
  return request({
    url: `/user/${userId}`,
    method: 'get'
  })
}

// 获取指定用户的商品列表
export function getGoodsByUserId(userId, params) {
  return request({
    url: `/user/${userId}/goods`,
    method: 'get',
    params
  })
}
