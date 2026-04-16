<template>
  <div class="home">
    <el-container>
      <el-header>
        <div class="header-content">
          <h1 class="logo">校园闲置物品交换平台</h1>
          <div class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/goods">商品列表</router-link>
            <template v-if="userStore.isLogin">
              <router-link to="/exchange">我的交换</router-link>
              <router-link to="/profile">个人中心</router-link>
              <a @click="handleLogout">退出</a>
            </template>
            <template v-else>
              <router-link to="/login">登录</router-link>
              <router-link to="/register">注册</router-link>
            </template>
          </div>
        </div>
      </el-header>

      <el-main>
        <div class="banner">
          <h2>欢迎来到校园闲置物品交换平台</h2>
          <p>让闲置物品流转起来，让资源得到更好的利用</p>
        </div>

        <div class="features">
          <el-card>
            <h3>发布商品</h3>
            <p>快速发布你的闲置物品</p>
          </el-card>
          <el-card>
            <h3>在线交易</h3>
            <p>安全便捷的交易方式</p>
          </el-card>
          <el-card>
            <h3>实时沟通</h3>
            <p>与买家实时在线沟通</p>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.home {
  width: 100%;
  min-height: 100vh;
}

.el-header {
  background-color: #409eff;
  color: white;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.logo {
  font-size: 20px;
  margin: 0;
}

.nav {
  display: flex;
  gap: 20px;
}

.nav a {
  color: white;
  text-decoration: none;
  cursor: pointer;
}

.nav a:hover {
  opacity: 0.8;
}

.el-main {
  max-width: 1200px;
  margin: 0 auto;
}

.banner {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  margin-bottom: 40px;
}

.banner h2 {
  font-size: 36px;
  margin-bottom: 10px;
}

.banner p {
  font-size: 18px;
  opacity: 0.9;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.features .el-card {
  text-align: center;
  padding: 20px;
}

.features h3 {
  margin-bottom: 10px;
  color: #409eff;
}
</style>
