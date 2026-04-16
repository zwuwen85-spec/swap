# 校园闲置物品交换平台 - 前端

基于 Vue 3 + Element Plus + Vite 的校园二手物品交易平台前端项目

## 技术栈

- Vue 3.x (渐进式框架)
- Vite 4.x (构建工具)
- Element Plus (UI组件库)
- Pinia (状态管理)
- Vue Router 4.x (路由管理)
- Axios (HTTP客户端)
- SCSS (CSS预处理器)

## 目录结构

```
frontend/
├── src/
│   ├── api/            # API接口
│   ├── assets/         # 静态资源
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/
│   ├── components/     # 公共组件
│   ├── views/          # 页面
│   ├── router/         # 路由配置
│   ├── store/          # Pinia状态管理
│   ├── utils/          # 工具函数
│   ├── composables/    # 组合式函数
│   ├── directives/     # 自定义指令
│   ├── App.vue         # 根组件
│   └── main.js         # 入口文件
├── public/             # 公共资源
├── index.html          # HTML模板
├── vite.config.js      # Vite配置
├── package.json        # 项目配置
└── .env.*              # 环境变量
```

## 快速开始

### 环境要求

- Node.js 16+
- npm 或 pnpm

### 安装步骤

1. **进入前端目录**
```bash
cd campus-swap-shop/frontend
```

2. **安装依赖**
```bash
npm install
# 或
pnpm install
```

3. **启动开发服务器**
```bash
npm run dev
```

4. **访问应用**
```
打开浏览器访问: http://localhost:3000
```

### 构建生产版本

```bash
npm run build
```

构建产物位于 `dist/` 目录

### 预览生产版本

```bash
npm run preview
```

## 配置说明

### 环境变量

**开发环境** (`.env.development`)：
```bash
VITE_API_BASE_URL=/api/v1
VITE_WS_URL=ws://localhost:8080/ws
```

**生产环境** (`.env.production`)：
```bash
VITE_API_BASE_URL=https://api.example.com/api/v1
VITE_WS_URL=wss://api.example.com/ws
```

### Vite配置

主要配置项（`vite.config.js`）：
- 端口：3000
- 代理：`/api` 代理到 `http://localhost:8080`
- 别名：`@` 指向 `src` 目录

## 项目结构说明

### API层 (`src/api/`)

- `request.js` - Axios实例配置和拦截器
- `user.js` - 用户相关接口
- `goods.js` - 商品相关接口
- `exchange.js` - 交换相关接口
- `message.js` - 消息相关接口

### 页面层 (`src/views/`)

- `home/` - 首页
- `user/` - 用户相关（登录、注册、个人中心）
- `goods/` - 商品相关（列表、详情、发布）
- `exchange/` - 交换相关
- `chat/` - 聊天相关
- `admin/` - 后台管理

### 组件层 (`src/components/`)

- `common/` - 公共组件（Header、Footer等）
- `goods/` - 商品组件（卡片、列表等）
- `upload/` - 上传组件

### 状态管理 (`src/store/`)

- `user.js` - 用户状态
- `goods.js` - 商品状态
- `app.js` - 应用状态

## 开发指南

### 代码规范

```bash
# 格式化代码
npm run lint

# 自动修复
npm run lint -- --fix
```

### 路由使用

```javascript
// 编程式导航
import { useRouter } from 'vue-router'

const router = useRouter()
router.push('/goods/123')

// 声明式导航
<router-link to="/goods/123">商品详情</router-link>
```

### 状态管理

```javascript
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
userStore.login({ username, password })
```

### API调用

```javascript
import { getGoodsList } from '@/api/goods'

const res = await getGoodsList({ page: 1, page_size: 20 })
```

## 组件使用示例

### 表单

```vue
<template>
  <el-form :model="form" :rules="rules" ref="formRef">
    <el-form-item label="用户名" prop="username">
      <el-input v-model="form.username" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleSubmit">提交</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref } from 'vue'

const form = ref({ username: '' })
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ]
}
</script>
```

### 表格

```vue
<template>
  <el-table :data="tableData">
    <el-table-column prop="title" label="标题" />
    <el-table-column prop="price" label="价格" />
  </el-table>
</template>

<script setup>
import { ref } from 'vue'

const tableData = ref([])
</script>
```

## 部署

### Nginx部署

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
    }
}
```

### Docker部署

```dockerfile
FROM node:16-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 开发计划

- [ ] 完善所有页面组件
- [ ] 实现商品发布功能
- [ ] 实现图片上传功能
- [ ] 实现WebSocket聊天
- [ ] 实现交换功能
- [ ] 实现评论功能
- [ ] 实现收藏功能
- [ ] 实现后台管理页面
- [ ] 响应式优化
- [ ] 性能优化

## 常见问题

**Q: 开发环境跨域问题？**
A: 已在 `vite.config.js` 中配置代理，无需额外处理

**Q: 如何修改端口？**
A: 修改 `vite.config.js` 中的 `server.port`

**Q: 如何添加新的页面？**
A: 在 `src/views/` 中创建组件，然后在 `src/router/index.js` 中添加路由

**Q: 如何添加新的API接口？**
A: 在 `src/api/` 中创建对应的文件，使用 `request` 发起请求

## 注意事项

1. **环境变量**：生产环境记得修改 `.env.production` 中的配置
2. **路由模式**：使用 HTML5 History 模式，需要服务器配置支持
3. **组件按需引入**：Element Plus 已配置按需引入
4. **图标使用**：使用 `@element-plus/icons-vue` 中的图标

## 问题反馈

如有问题，请提交 Issue

## 作者

毕业设计项目
