import { defineStore } from 'pinia'
import { login as loginApi, register as registerApi, getUserInfo } from '@/api/user'
import { getUnreadCount } from '@/api/message'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    isLogin: !!localStorage.getItem('token'),
    unreadMessageCount: 0,
    ws: null, // WebSocket实例
    latestMessage: null, // 最新收到的消息
    onlineStatus: null // 最新收到的在线状态更新
  }),

  getters: {
    userId: (state) => state.userInfo?.id,
    username: (state) => state.userInfo?.username,
    avatar: (state) => state.userInfo?.avatar
  },

  actions: {
    // 初始化 WebSocket 连接
    initWebSocket() {
      if (!this.token || this.ws) return

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      // 假设后端运行在同域名的 8080 端口，可根据实际情况修改
      const host = process.env.NODE_ENV === 'development' ? 'localhost:8080' : window.location.host
      
      this.ws = new WebSocket(`${protocol}//${host}/api/v1/ws/chat?token=${this.token}`)

      this.ws.onopen = () => {
        console.log('WebSocket 连接成功')
        // 发送心跳的机制如果需要的话可以放这里，不过后端应该有 ping/pong
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          // 收到新消息时，如果是 message 类型，全局未读数加一，并触发刷新
          if (data.type === 'message') {
            this.unreadMessageCount += 1
            this.latestMessage = data
            // 你也可以在这里弹出 Notification
          } else if (data.type === 'status') {
            this.onlineStatus = data
          }
        } catch (e) {
          console.error('WebSocket 消息解析失败', e)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket 连接关闭')
        this.ws = null
        // 可以在这里实现断线重连逻辑
        if (this.isLogin) {
          setTimeout(() => {
            this.initWebSocket()
          }, 3000)
        }
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket 错误:', error)
      }
    },

    // 关闭 WebSocket 连接
    closeWebSocket() {
      if (this.ws) {
        this.ws.close()
        this.ws = null
      }
    },

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
      this.initWebSocket()
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
      this.initWebSocket()
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
      this.unreadMessageCount = 0
      localStorage.removeItem('token')
      this.closeWebSocket()
    },

    // 获取未读消息数
    async fetchUnreadCount() {
      if (!this.isLogin) return
      try {
        const res = await getUnreadCount()
        this.unreadMessageCount = res.data.total
      } catch (error) {
        console.error('获取全局未读消息数失败:', error)
      }
    }
  }
})
