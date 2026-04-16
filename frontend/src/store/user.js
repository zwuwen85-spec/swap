import { defineStore } from 'pinia'
import { login as loginApi, register as registerApi, getUserInfo } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    isLogin: !!localStorage.getItem('token')
  }),

  getters: {
    userId: (state) => state.userInfo?.id,
    username: (state) => state.userInfo?.username,
    avatar: (state) => state.userInfo?.avatar
  },

  actions: {
    // 登录
    async login(loginForm) {
      const res = await loginApi(loginForm)
      this.token = res.data.token
      this.userInfo = {
        id: res.data.user_id,
        username: res.data.username,
        nickname: res.data.nickname,
        avatar: res.data.avatar
      }
      this.isLogin = true
      localStorage.setItem('token', res.data.token)
      return res
    },

    // 注册
    async register(registerForm) {
      const res = await registerApi(registerForm)
      this.token = res.data.token
      this.userInfo = {
        id: res.data.user_id,
        username: res.data.username,
        nickname: res.data.nickname,
        avatar: res.data.avatar
      }
      this.isLogin = true
      localStorage.setItem('token', res.data.token)
      return res
    },

    // 获取用户信息
    async getUserInfo() {
      const res = await getUserInfo()
      this.userInfo = res.data
      return res
    },

    // 退出登录
    logout() {
      this.token = ''
      this.userInfo = null
      this.isLogin = false
      localStorage.removeItem('token')
    }
  }
})
