import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home/index.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/goods',
    name: 'GoodsList',
    component: () => import('@/views/goods/List.vue'),
    meta: { title: '商品列表' }
  },
  {
    path: '/goods/:id',
    name: 'GoodsDetail',
    component: () => import('@/views/goods/Detail.vue'),
    meta: { title: '商品详情' }
  },
  {
    path: '/my-goods',
    name: 'MyGoods',
    component: () => import('@/views/goods/MyGoods.vue'),
    meta: { title: '我的发布', requiresAuth: true }
  },
  {
    path: '/publish',
    name: 'Publish',
    component: () => import('@/views/goods/Publish.vue'),
    meta: { title: '发布商品', requiresAuth: true }
  },
  {
    path: '/exchange',
    name: 'ExchangeList',
    component: () => import('@/views/exchange/List.vue'),
    meta: { title: '我的交换', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/user/:id',
    name: 'UserDetail',
    component: () => import('@/views/user/UserDetail.vue'),
    meta: { title: '用户详情' }
  },
  {
    path: '/chat',
    name: 'ChatList',
    component: () => import('@/views/chat/List.vue'),
    meta: { title: '消息列表', requiresAuth: true }
  },
  {
    path: '/chat/:userId',
    name: 'ChatRoom',
    component: () => import('@/views/chat/Room.vue'),
    meta: { title: '聊天', requiresAuth: true }
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: () => import('@/views/favorite/Favorites.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/my-comments',
    name: 'MyComments',
    component: () => import('@/views/comment/MyComments.vue'),
    meta: { title: '我的评论', requiresAuth: true }
  },
  {
    path: '/notifications',
    name: 'Notifications',
    component: () => import('@/views/notification/Notifications.vue'),
    meta: { title: '通知中心', requiresAuth: true }
  },
  {
    path: '/my-reports',
    name: 'MyReports',
    component: () => import('@/views/report/MyReports.vue'),
    meta: { title: '我的举报', requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面不存在' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title || '校园闲置物品交换平台'

  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')
    if (!token) {
      next('/login')
      return
    }

    // 如果有token但没有用户信息，获取用户信息
    const userStore = useUserStore()
    if (token && !userStore.userInfo) {
      try {
        await userStore.getUserInfo()
        console.log('路由守卫：已获取用户信息', userStore.userInfo)
      } catch (error) {
        console.error('获取用户信息失败:', error)
        // 获取用户信息失败，可能token过期，跳转到登录页
        localStorage.removeItem('token')
        next('/login')
        return
      }
    }
  }

  next()
})

export default router
