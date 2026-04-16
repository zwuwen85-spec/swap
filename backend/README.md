# 校园闲置物品交换平台 - 后端

基于 Go + Gin + MySQL + Redis 的校园二手物品交易平台后端服务

## 技术栈

- Go 1.21+
- Gin (Web框架)
- GORM (ORM)
- MySQL 8.0+ (数据库)
- Redis 6.0+ (缓存)
- JWT (认证)
- WebSocket (实时通信)

## 目录结构

```
backend/
├── cmd/                 # 应用入口
│   └── server/
│       └── main.go     # 主入口
├── internal/           # 私有代码
│   ├── handler/       # HTTP处理器
│   ├── logic/         # 业务逻辑
│   ├── model/         # 数据模型
│   ├── middleware/    # 中间件
│   ├── cache/         # 缓存层
│   ├── dto/           # 数据传输对象
│   └── svc/           # 服务上下文
├── pkg/               # 公共库
│   ├── jwt/          # JWT工具
│   ├── response/     # 统一响应
│   ├── utils/        # 工具函数
│   ├── errors/       # 错误定义
│   └── logger/       # 日志封装
├── routes/            # 路由定义
├── config/            # 配置文件
├── scripts/           # 脚本
└── test/              # 测试
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 安装步骤

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
cp config/config.yaml config/config.local.yaml
# 编辑 config.local.yaml，修改数据库等配置
```

4. **初始化数据库**
```bash
mysql -u root -p < scripts/init.sql
```

5. **运行服务**
```bash
# 方式1：使用 Makefile
make dev

# 方式2：直接运行
go run cmd/server/main.go

# 方式3：先构建再运行
make build
./bin/campus-swap-shop
```

6. **访问服务**
```
API: http://localhost:8080
健康检查: http://localhost:8080/health
```

## 配置说明

配置文件位于 `config/config.yaml`：

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # 运行模式（debug/release）

mysql:
  host: localhost         # MySQL地址
  port: 3306             # MySQL端口
  database: campus_swap  # 数据库名
  username: root         # 用户名
  password: ""           # 密码

redis:
  host: localhost         # Redis地址
  port: 6379             # Redis端口
  password: ""           # 密码
  db: 0                  # 数据库

jwt:
  secret: "your-secret"  # JWT密钥
  expire: 604800        # 过期时间（秒）
```

## 开发指南

### 运行测试

```bash
# 运行所有测试
make test

# 查看测试覆盖率
make test-coverage
```

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make vet
```

## API文档

启动服务后，访问 Swagger 文档：
```
http://localhost:8080/swagger/index.html
```

主要接口：
- `POST /api/v1/user/register` - 用户注册
- `POST /api/v1/user/login` - 用户登录
- `GET /api/v1/user/info` - 获取用户信息（需认证）
- `GET /api/v1/goods/list` - 商品列表
- `GET /api/v1/goods/detail` - 商品详情
- `POST /api/v1/goods/create` - 发布商品（需认证）

详细API文档请查看项目根目录的 `docs/02-API设计文档.md`

## Makefile命令

| 命令 | 说明 |
|-----|------|
| `make dev` | 开发模式运行 |
| `make build` | 构建可执行文件 |
| `make test` | 运行测试 |
| `make test-coverage` | 测试覆盖率 |
| `make clean` | 清理构建文件 |
| `make fmt` | 格式化代码 |
| `make vet` | 代码检查 |
| `make deps` | 下载依赖 |

## 开发计划

- [ ] 用户模块（注册、登录、个人信息）
- [ ] 商品模块（发布、列表、详情、搜索）
- [ ] 交换模块（发起、处理、状态流转）
- [ ] 聊天模块（WebSocket实时通信）
- [ ] 评论模块（发表评论、评分）
- [ ] 收藏模块（收藏、取消收藏）
- [ ] 通知模块（系统通知、消息推送）
- [ ] 后台管理（用户管理、商品管理、举报处理）

## 注意事项

1. **生产环境配置**
   - 修改 `jwt.secret` 为随机字符串
   - 修改 `mysql.password` 为实际密码
   - 设置 `server.mode` 为 `release`

2. **数据库安全**
   - 不要在生产环境使用默认密码
   - 定期备份数据库
   - 使用SSL连接数据库

3. **日志管理**
   - 日志文件位于 `logs/` 目录
   - 定期清理过期日志
   - 生产环境建议使用日志收集系统

## 问题反馈

如有问题，请提交 Issue

## 作者

毕业设计项目
