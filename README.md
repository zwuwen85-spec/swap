# 校园闲置物品交换平台

> 基于 Go + Vue 3 的校园二手物品交易平台，支持物品交易、交换和实时聊天

## 📖 项目简介

本项目是一个完整的校园闲置物品交换平台，学生可以在平台上发布闲置物品、进行交易或交换，并通过实时聊天功能沟通。采用前后端分离架构，具备良好的扩展性和实用性。

### ✨ 主要功能

- 👤 **用户系统**：注册登录、个人信息管理、信誉评分
- 📦 **商品管理**：发布商品、商品列表、详情展示、图片上传
- 🔄 **交换系统**：发起交换、状态流转、交换历史
- 💬 **实时聊天**：WebSocket实时通信、离线消息、商品卡片分享
- 🔍 **搜索功能**：关键词搜索、分类筛选、热门推荐
- ⭐ **互动功能**：收藏商品、发表评论、点赞
- 🔐 **后台管理**：用户管理、商品审核、举报处理、数据统计
- 📊 **数据统计**：用户增长、交易统计、热门商品分析

### 🎯 项目亮点

- 🏗️ **清晰的分层架构**：Handler → Logic → Model 三层分离
- 🚀 **高性能**：Redis缓存、读写分离、连接池优化
- 🔒 **安全可靠**：JWT鉴权、密码加密、SQL注入防护、XSS防护
- 💭 **实时通信**：WebSocket实现即时聊天
- 📱 **多端支持**：使用uni-app可编译为H5/小程序
- 📈 **可扩展**：微服务架构，易于扩展功能模块

---

## 🛠 技术栈

### 后端技术

| 技术 | 版本 | 说明 |
|-----|------|------|
| Go | 1.21+ | 高性能编程语言 |
| Gin | - | Web框架 |
| GORM | - | ORM框架 |
| MySQL | 8.0+ | 关系型数据库 |
| Redis | 6.0+ | 缓存数据库 |
| JWT | - | 身份认证 |
| WebSocket | - | 实时通信 |

### 前端技术

| 技术 | 版本 | 说明 |
|-----|------|------|
| Vue | 3.x | 渐进式框架 |
| Vite | 4.x | 构建工具 |
| Element Plus | - | UI组件库 |
| Pinia | - | 状态管理 |
| Vue Router | 4.x | 路由管理 |
| Axios | - | HTTP客户端 |
| uni-app | - | 多端方案 |

---

## 📂 项目结构

```
campus-swap-shop/
├── backend/                  # 后端项目
│   ├── cmd/                 # 应用入口
│   ├── internal/            # 私有代码
│   │   ├── handler/        # HTTP处理器
│   │   ├── logic/          # 业务逻辑
│   │   ├── model/          # 数据模型
│   │   ├── middleware/     # 中间件
│   │   └── cache/          # 缓存层
│   ├── pkg/                # 公共库
│   ├── routes/             # 路由
│   ├── config/             # 配置
│   ├── scripts/            # 脚本
│   └── docs/               # 文档
│
├── frontend/                # 前端项目
│   ├── src/
│   │   ├── api/           # API接口
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── router/        # 路由
│   │   ├── store/         # 状态管理
│   │   └── utils/         # 工具
│   └── public/
│
└── docs/                   # 项目文档
    ├── 01-数据库设计.md
    ├── 02-API设计文档.md
    ├── 03-安全设计.md
    ├── 04-项目架构设计.md
    ├── 05-缓存设计.md
    ├── 06-开发流程与计划.md
    └── 07-前端设计文档.md
```

---

## 🚀 快速开始

> 💡 **新用户推荐**：先查看 [快速开始指南](QUICKSTART.md)，5分钟即可启动项目！

### 环境要求

**必需软件：**
- Go 1.21+
- Node.js 16+
- MySQL 8.0+
- Redis 6.0+
- Git

---

### 📖 一分钟快速启动

```bash
# 1. 初始化数据库
mysql -u root -p < scripts/init.sql

# 2. 启动后端
cd backend
go run cmd/server/main.go

# 3. 启动前端（新终端）
cd frontend
npm install  # 首次运行
npm run dev
```

访问：http://localhost:3000

详细说明请查看：[快速开始指南](QUICKSTART.md)

---

### 后端启动

1. **克隆项目**
```bash
git clone https://github.com/yourusername/campus-swap-shop.git
cd campus-swap-shop/backend
```

2. **安装依赖**
```bash
go mod download
```

3. **配置环境**
```bash
cp config/config.example.yaml config/config.local.yaml
# 编辑 config.local.yaml
```

4. **初始化数据库**
```bash
mysql -u root -p < scripts/init.sql
```

5. **启动服务**
```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

---

### 前端启动

1. **进入前端目录**
```bash
cd frontend
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

前端将在 `http://localhost:3000` 启动

---

## 📖 文档

### 设计文档

- [数据库设计](docs/01-数据库设计.md) - 完整的数据表设计
- [API设计文档](docs/02-API设计文档.md) - RESTful API接口规范
- [安全设计](docs/03-安全设计.md) - 安全防护方案
- [项目架构设计](docs/04-项目架构设计.md) - 系统架构说明
- [缓存设计](docs/05-缓存设计.md) - Redis缓存策略
- [开发流程与计划](docs/06-开发流程与计划.md) - 开发指南
- [前端设计文档](docs/07-前端设计文档.md) - 前端架构说明

### API文档

后端启动后，访问 Swagger 文档：
```
http://localhost:8080/swagger/index.html
```

---

## 🔧 配置说明

### 后端配置

```yaml
# config/config.yaml
server:
  port: 8080
  mode: debug

mysql:
  host: localhost
  port: 3306
  database: campus_swap
  username: root
  password: ""

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  expire: 604800
```

### 前端配置

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws
```

---

## 🧪 测试

### 后端测试

```bash
# 运行所有测试
go test ./...

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试

```bash
# 单元测试
npm run test

# E2E测试
npm run test:e2e
```

---

## 📦 部署

### Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 生产环境部署

详见：[开发流程与计划 - 部署方案](docs/06-开发流程与计划.md#八部署方案)

---

## 📊 系统截图

> 待补充

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 提交规范

```
<type>(<scope>): <subject>

type: feat/fix/docs/style/refactor/test/chore
```

示例：
```
feat(user): 添加用户注册功能

fix(goods): 修复商品列表分页bug

docs(readme): 更新安装说明
```

---

## 📝 开发计划

- [x] 项目设计文档
- [ ] 后端基础框架
- [ ] 用户模块
- [ ] 商品模块
- [ ] 交换模块
- [ ] 聊天模块
- [ ] 前端开发
- [ ] 测试与优化
- [ ] 部署上线

详细计划见：[开发流程与计划](docs/06-开发流程与计划.md)

---

## ❓ 常见问题

**Q: 如何修改端口？**
A: 编辑 `config/config.yaml`，修改 `server.port`

**Q: 如何重置数据库？**
A: 删除数据库后重新执行 `scripts/init.sql`

**Q: Docker 容器无法连接 MySQL？**
A: 检查 `docker-compose.yml` 中的网络配置，确保使用服务名作为主机名

**Q: WebSocket 连接失败？**
A: 检查 Nginx 配置，确保正确代理 WebSocket 连接

---

## 📄 许可证

[MIT License](LICENSE)

---

## 👨‍💻 作者

**你的名字** - 毕业设计项目

---

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架
- [Element Plus](https://element-plus.org/) - Vue 3 UI组件库
- [GORM](https://gorm.io/) - Go ORM库
- [Redis](https://redis.io/) - 内存数据库

---

## 📮 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 提交 Issue
- 发送邮件至：your@email.com

---

<div align="center">

**如果这个项目对你有帮助，请给一个 ⭐️ Star**

</div>
