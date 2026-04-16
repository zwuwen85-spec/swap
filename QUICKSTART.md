# 校园闲置物品交换平台 - 快速开始

## 前提条件

在开始之前，请确保已安装以下软件：

- **Go** 1.21+ - [下载地址](https://golang.org/dl/)
- **Node.js** 16+ - [下载地址](https://nodejs.org/)
- **MySQL** 8.0+ - [下载地址](https://dev.mysql.com/downloads/mysql/)
- **Redis** 6.0+ - [下载地址](https://redis.io/download)

## 快速启动（5分钟）

### 1️⃣ 初始化数据库（30秒）

```bash
# Windows
mysql -u root -p < D:\毕业设计\CampusSwapShop\scripts\init.sql

# 或使用MySQL客户端工具导入 scripts/init.sql 文件
```

### 2️⃣ 配置后端（30秒）

编辑 `backend/config/config.yaml`：

```yaml
mysql:
  host: localhost
  port: 3306
  database: campus_swap
  username: root
  password: "你的MySQL密码"  # 修改这里

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

### 3️⃣ 启动后端（1分钟）

```bash
cd backend
go mod download
go run cmd/server/main.go
```

看到以下输出表示成功：
```
{"level":"INFO","ts":"...","msg":"MySQL连接成功"}
{"level":"INFO","ts":"...","msg":"Redis连接成功"}
{"level":"INFO","ts":"...","msg":"服务器启动","port":"8080","mode":"debug"}
```

### 4️⃣ 启动前端（2分钟）

```bash
cd frontend
npm install  # 首次运行，约1-2分钟
npm run dev
```

看到以下输出表示成功：
```
VITE v5.0.0  ready in 500 ms

➜  Local:   http://localhost:3000/
```

### 5️⃣ 访问应用

打开浏览器访问：**http://localhost:3000**

## 测试用户模块

### 注册

1. 点击"注册"
2. 输入用户名：`testuser`
3. 输入密码：`123456`
4. 点击"注册"
5. ✅ 显示"注册成功"并跳转到首页

### 登录

1. 点击"登录"
2. 输入用户名：`testuser`
3. 输入密码：`123456`
4. 点击"登录"
5. ✅ 显示"登录成功"并跳转到首页

### 查看用户信息

1. 登录后，点击右上角"个人中心"
2. ✅ 可以看到当前登录的用户信息

## 常见问题

### Q: 后端启动失败，提示MySQL连接错误？

**A:** 检查以下几点：
1. MySQL是否已启动
2. `config.yaml` 中的密码是否正确
3. 数据库 `campus_swap` 是否已创建

### Q: 前端启动失败，提示端口被占用？

**A:** 修改 `frontend/vite.config.js` 中的端口号：
```javascript
server: {
  port: 3001  // 改成其他端口
}
```

### Q: 注册/登录时提示CORS错误？

**A:** 确保后端服务已启动，并检查：
1. 后端端口是否为8080
2. 前端代理配置是否正确

### Q: 忘记MySQL密码？

**A:** 重置MySQL root密码：
```bash
# Windows
1. 停止MySQL服务
2. 以安全模式启动MySQL
3. 重置密码
4. 重启MySQL服务
```

详细步骤请搜索："MySQL重置root密码 Windows"

## 项目结构

```
CampusSwapShop/
├── backend/           # 后端（Go）
│   ├── cmd/          # 入口文件
│   ├── internal/     # 内部代码
│   ├── pkg/          # 公共库
│   └── config/       # 配置文件
│
├── frontend/         # 前端（Vue 3）
│   ├── src/         # 源代码
│   └── public/      # 静态资源
│
├── docs/            # 文档
└── scripts/         # SQL脚本
```

## 下一步

用户模块已完成，可以：

1. **测试完整功能** - 查看 [用户模块测试指南](docs/用户模块测试指南.md)
2. **查看API文档** - 使用Postman测试后端接口
3. **继续开发** - 开始实现商品模块

## 需要帮助？

- 查看 [完整文档](docs/)
- 查看 [常见问题](docs/常见问题.md)
- 提交Issue

---

**祝你使用愉快！** 🎉
