# 校园闲置物品交换平台 - API设计文档

## 一、接口规范

### 1. 基础规范

- 协议：HTTPS（生产环境）
- 数据格式：JSON
- 字符编码：UTF-8
- API版本：v1

### 2. 统一请求格式

**Headers：**
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**Query参数：**
```
GET /api/v1/goods/list?page=1&page_size=20
```

**Body参数：**
```json
{
  "title": "商品标题",
  "price": 99.99
}
```

### 3. 统一响应格式

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**失败响应：**
```json
{
  "code": 10001,
  "message": "参数错误",
  "data": null
}
```

**分页响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 4. 状态码规范

| code | message | 说明 |
|------|---------|------|
| 0 | success | 成功 |
| 10001 | 参数错误 | 请求参数验证失败 |
| 10002 | 未登录 | 需要登录 |
| 10003 | 无权限 | 权限不足 |
| 10004 | 资源不存在 | 请求的资源不存在 |
| 10005 | 重复操作 | 重复提交/操作 |
| 20001 | 用户不存在 | 用户不存在 |
| 20002 | 密码错误 | 登录密码错误 |
| 20003 | 用户已存在 | 注册时用户已存在 |
| 30001 | 商品不存在 | 商品不存在 |
| 30002 | 商品已下架 | 商品已下架 |
| 30003 | 商品已售出 | 商品已售出 |
| 40001 | 交换请求不存在 | 交换请求不存在 |
| 40002 | 不能处理自己的请求 | 不能处理自己的交换请求 |
| 50001 | 服务器错误 | 服务器内部错误 |

---

## 二、用户模块 API

### 1.1 用户注册

**接口：** `POST /api/v1/user/register`

**请求参数：**
```json
{
  "username": "string",      // 必填，3-20字符
  "password": "string",      // 必填，6-20字符
  "email": "string",         // 可选
  "phone": "string",         // 可选
  "student_id": "string"     // 可选
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "user_id": 123,
    "username": "张三",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

### 1.2 用户登录

**接口：** `POST /api/v1/user/login`

**请求参数：**
```json
{
  "username": "string",      // 必填
  "password": "string"       // 必填
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "user_id": 123,
    "username": "张三",
    "avatar": "https://...",
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expire": 604800        // token过期时间（秒）
  }
}
```

---

### 1.3 获取当前用户信息

**接口：** `GET /api/v1/user/info`

**是否需要登录：** 是

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "username": "zhangsan",
    "nickname": "张三",
    "avatar": "https://...",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "gender": 1,
    "student_id": "202001001",
    "school": "XX大学",
    "status": 1,
    "credit_score": 100,
    "create_time": 1648000000
  }
}
```

---

### 1.4 更新用户信息

**接口：** `PUT /api/v1/user/info`

**是否需要登录：** 是

**请求参数：**
```json
{
  "nickname": "string",
  "avatar": "string",
  "gender": 1,
  "phone": "string",
  "school": "string"
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "更新成功",
  "data": null
}
```

---

### 1.5 修改密码

**接口：** `POST /api/v1/user/change-password`

**是否需要登录：** 是

**请求参数：**
```json
{
  "old_password": "string",
  "new_password": "string"
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "密码修改成功",
  "data": null
}
```

---

### 1.6 上传头像

**接口：** `POST /api/v1/user/avatar`

**是否需要登录：** 是

**请求参数：** `multipart/form-data`
- file: 图片文件（最大5MB，支持jpg/png/gif）

**响应数据：**
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "url": "https://example.com/uploads/avatar/xxx.jpg"
  }
}
```

---

## 三、商品模块 API

### 2.1 发布商品

**接口：** `POST /api/v1/goods/create`

**是否需要登录：** 是

**请求参数：**
```json
{
  "title": "string",           // 必填，1-100字符
  "description": "string",     // 可选
  "category_id": 1,            // 必填
  "type": 1,                   // 必填，1售卖 2交换 3均可
  "price": 99.99,              // 可选，售卖时必填
  "original_price": 199.99,    // 可选
  "images": ["url1", "url2"],  // 必填，至少1张
  "condition": 1,              // 必填，1全新 2九成新 3八成新 4七成新
  "tags": "标签1,标签2",        // 可选
  "location": "图书馆",         // 可选
  "latitude": 39.908823,       // 可选
  "longitude": 116.397470      // 可选
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "发布成功",
  "data": {
    "goods_id": 456
  }
}
```

---

### 2.2 商品列表

**接口：** `GET /api/v1/goods/list`

**请求参数：**
```
page: int         // 页码，默认1
page_size: int    // 每页数量，默认20，最大50
category_id: int  // 分类ID，可选
type: int         // 类型，1售卖 2交换 3均可
condition: int    // 成色
keyword: string   // 搜索关键词
sort: string      // 排序，time_desc/price_asc/price_desc，默认time_desc
```

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 456,
        "title": "出售二手MacBook Pro",
        "description": "2020款，使用良好...",
        "price": 5999.00,
        "original_price": 8999.00,
        "type": 1,
        "images": ["url1", "url2"],
        "condition": 2,
        "view_count": 128,
        "favorite_count": 15,
        "status": 1,
        "location": "图书馆",
        "create_time": 1648000000,
        "user": {
          "id": 123,
          "username": "zhangsan",
          "avatar": "https://...",
          "credit_score": 100
        },
        "category": {
          "id": 1,
          "name": "数码产品"
        }
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 2.3 商品详情

**接口：** `GET /api/v1/goods/detail`

**请求参数：**
```
id: int    // 商品ID，必填
```

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 456,
    "title": "出售二手MacBook Pro",
    "description": "详细描述...",
    "price": 5999.00,
    "original_price": 8999.00,
    "type": 1,
    "images": ["url1", "url2", "url3"],
    "condition": 2,
    "tags": "苹果,笔记本,电脑",
    "view_count": 128,
    "favorite_count": 15,
    "status": 1,
    "location": "图书馆",
    "latitude": 39.908823,
    "longitude": 116.397470,
    "create_time": 1648000000,
    "update_time": 1648100000,
    "is_favorited": false,        // 当前用户是否已收藏
    "user": {
      "id": 123,
      "username": "zhangsan",
      "nickname": "张三",
      "avatar": "https://...",
      "credit_score": 100,
      "school": "XX大学"
    },
    "category": {
      "id": 1,
      "name": "数码产品"
    },
    "comments": [                 // 最新评论
      {
        "id": 1,
        "user": {
          "username": "lisi",
          "avatar": "https://..."
        },
        "content": "很好的商品",
        "rating": 5,
        "create_time": 1648000100
      }
    ]
  }
}
```

---

### 2.4 更新商品

**接口：** `PUT /api/v1/goods/update`

**是否需要登录：** 是

**请求参数：**
```json
{
  "id": 456,                      // 商品ID
  "title": "string",
  "description": "string",
  "price": 99.99,
  "images": ["url1", "url2"],
  "status": 0                     // 0下架 1在售
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "更新成功",
  "data": null
}
```

---

### 2.5 删除商品

**接口：** `DELETE /api/v1/goods/delete`

**是否需要登录：** 是

**请求参数：**
```json
{
  "id": 456    // 商品ID
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

---

### 2.6 我的发布

**接口：** `GET /api/v1/goods/my`

**是否需要登录：** 是

**请求参数：**
```
page: int
page_size: int
status: int    // 状态，0全部 1在售 2已售 3已交换
```

**响应数据：**（同商品列表）

---

### 2.7 搜索商品

**接口：** `GET /api/v1/goods/search`

**请求参数：**
```
keyword: string    // 必填
page: int
page_size: int
category_id: int
type: int
```

**响应数据：**（同商品列表）

---

### 2.8 上传商品图片

**接口：** `POST /api/v1/goods/upload`

**是否需要登录：** 是

**请求参数：** `multipart/form-data`
- file: 图片文件（最大5MB，支持jpg/png/gif）

**响应数据：**
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "url": "https://example.com/uploads/goods/xxx.jpg"
  }
}
```

---

## 四、交换模块 API

### 3.1 发起交换请求

**接口：** `POST /api/v1/exchange/create`

**是否需要登录：** 是

**请求参数：**
```json
{
  "goods_id": 456,              // 目标商品ID
  "type": 1,                    // 1购买 2交换
  "my_goods_id": 789,           // 我的商品ID（交换时必填）
  "message": "可以便宜点吗？"    // 附加留言
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "交换请求已发送",
  "data": {
    "exchange_id": 1
  }
}
```

---

### 3.2 接受交换请求

**接口：** `POST /api/v1/exchange/accept`

**是否需要登录：** 是

**请求参数：**
```json
{
  "exchange_id": 1
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "已接受交换请求",
  "data": null
}
```

---

### 3.3 拒绝交换请求

**接口：** `POST /api/v1/exchange/reject`

**是否需要登录：** 是

**请求参数：**
```json
{
  "exchange_id": 1,
  "reason": "价格不合适"      // 拒绝原因
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "已拒绝交换请求",
  "data": null
}
```

---

### 3.4 取消交换请求

**接口：** `POST /api/v1/exchange/cancel`

**是否需要登录：** 是

**请求参数：**
```json
{
  "exchange_id": 1
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "已取消交换请求",
  "data": null
}
```

---

### 3.5 我的交换列表

**接口：** `GET /api/v1/exchange/list`

**是否需要登录：** 是

**请求参数：**
```
type: string    // incoming（收到的） / outgoing（发出的） / all（全部）
status: int     // 状态筛选，0待处理 1已接受 2已拒绝 3已取消 4已完成
page: int
page_size: int
```

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "initiator": {
          "id": 123,
          "username": "zhangsan",
          "avatar": "https://..."
        },
        "target": {
          "id": 456,
          "username": "lisi"
        },
        "goods": {
          "id": 789,
          "title": "商品标题",
          "images": ["url1"],
          "price": 99.99
        },
        "my_goods": {              // 交换时有此字段
          "id": 999,
          "title": "我的商品"
        },
        "type": 1,
        "message": "留言",
        "status": 0,
        "reject_reason": null,
        "create_time": 1648000000,
        "update_time": 1648100000
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 3.6 完成交换

**接口：** `POST /api/v1/exchange/complete`

**是否需要登录：** 是

**请求参数：**
```json
{
  "exchange_id": 1
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "交换已完成",
  "data": null
}
```

---

## 五、收藏模块 API

### 4.1 添加收藏

**接口：** `POST /api/v1/favorite/add`

**是否需要登录：** 是

**请求参数：**
```json
{
  "goods_id": 456
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "收藏成功",
  "data": null
}
```

---

### 4.2 取消收藏

**接口：** `DELETE /api/v1/favorite/remove`

**是否需要登录：** 是

**请求参数：**
```json
{
  "goods_id": 456
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "已取消收藏",
  "data": null
}
```

---

### 4.3 我的收藏列表

**接口：** `GET /api/v1/favorite/list`

**是否需要登录：** 是

**请求参数：**
```
page: int
page_size: int
```

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "goods": {
          "id": 456,
          "title": "商品标题",
          "price": 99.99,
          "images": ["url1"],
          "status": 1,
          "user": {
            "username": "zhangsan"
          }
        },
        "create_time": 1648000000
      }
    ],
    "total": 30,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 六、评论模块 API

### 5.1 发表评论

**接口：** `POST /api/v1/comment/create`

**是否需要登录：** 是

**请求参数：**
```json
{
  "goods_id": 456,
  "content": "商品很好，卖家人也不错",
  "rating": 5,               // 评分1-5
  "parent_id": 0             // 父评论ID，0表示主评论
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "评论成功",
  "data": {
    "comment_id": 1
  }
}
```

---

### 5.2 商品评论列表

**接口：** `GET /api/v1/comment/list`

**请求参数：**
```
goods_id: int    // 商品ID，必填
page: int
page_size: int
```

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user": {
          "id": 123,
          "username": "zhangsan",
          "avatar": "https://..."
        },
        "content": "商品很好",
        "rating": 5,
        "create_time": 1648000000,
        "replies": [              // 子评论
          {
            "id": 2,
            "user": {
              "username": "seller"
            },
            "content": "谢谢亲的好评",
            "create_time": 1648000100
          }
        ]
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 七、通知模块 API

### 6.1 通知列表

**接口：** `GET /api/v1/notification/list`

**是否需要登录：** 是

**请求参数：**
```
page: int
page_size: int
type: int         // 类型筛选，1系统 2交换 3评论
```

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "type": 2,
        "title": "新的交换请求",
        "content": "张三想要交换你的商品",
        "link": "/exchange/1",
        "is_read": 0,
        "create_time": 1648000000
      }
    ],
    "total": 50,
    "unread_count": 5,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 6.2 标记已读

**接口：** `POST /api/v1/notification/read`

**是否需要登录：** 是

**请求参数：**
```json
{
  "id": 1              // 通知ID，不传则全部标记为已读
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "已标记为已读",
  "data": null
}
```

---

## 八、分类模块 API

### 7.1 获取所有分类

**接口：** `GET /api/v1/category/list`

**响应数据：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "数码产品",
      "parent_id": 0,
      "level": 1,
      "icon": "https://...",
      "children": [
        {
          "id": 11,
          "name": "手机",
          "parent_id": 1,
          "level": 2
        },
        {
          "id": 12,
          "name": "电脑",
          "parent_id": 1,
          "level": 2
        }
      ]
    }
  ]
}
```

---

## 九、WebSocket 实时通信

### 8.1 连接

**接口：** `ws://domain/ws/chat`

**连接参数：**
```
token: string    // JWT token
```

### 8.2 消息格式

**客户端发送：**
```json
{
  "type": 1,               // 1文本 2图片 3商品卡片
  "receiver_id": 456,
  "content": "你好",
  "goods_id": 789          // type=3时必填
}
```

**服务端推送：**
```json
{
  "id": 1,
  "sender_id": 123,
  "receiver_id": 456,
  "content": "你好",
  "type": 1,
  "goods_id": null,
  "is_read": 0,
  "create_time": 1648000000
}
```

### 8.3 心跳保活

客户端每30秒发送一次ping：

```json
{"type": "ping"}
```

服务端响应：

```json
{"type": "pong"}
```

### 8.4 在线状态

**用户上线/离线时，服务端广播：**

```json
{
  "type": "status",
  "user_id": 123,
  "online": true
}
```

---

## 十、举报模块 API

### 9.1 提交举报

**接口：** `POST /api/v1/report/create`

**是否需要登录：** 是

**请求参数：**
```json
{
  "target_type": 1,        // 1商品 2用户 3评论
  "target_id": 456,
  "reason": "虚假商品",
  "description": "商品描述与实际不符"
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "举报已提交",
  "data": {
    "report_id": 1
  }
}
```

---

## 十一、后台管理 API

### 10.1 管理员登录

**接口：** `POST /api/v1/admin/login`

**请求参数：**
```json
{
  "username": "admin",
  "password": "123456"
}
```

**响应数据：**
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "admin_id": 1,
    "username": "admin",
    "token": "eyJhbGci...",
    "role": 1
  }
}
```

---

### 10.2 用户管理

**用户列表：** `GET /api/v1/admin/users`

**封禁用户：** `POST /api/v1/admin/user/ban`

**解封用户：** `POST /api/v1/admin/user/unban`

---

### 10.3 商品管理

**商品列表：** `GET /api/v1/admin/goods`

**下架商品：** `POST /api/v1/admin/goods/remove`

---

### 10.4 举报处理

**举报列表：** `GET /api/v1/admin/reports`

**处理举报：** `POST /api/v1/admin/report/handle`

**请求参数：**
```json
{
  "report_id": 1,
  "status": 1,                // 1已处理 2已驳回
  "handle_result": "违规内容已删除"
}
```

---

## 十二、接口限流规则

为了防止恶意请求，所有接口都需要限流：

| 用户类型 | 限流规则 |
|---------|---------|
| 未登录用户 | 100次/小时 |
| 已登录用户 | 300次/小时 |
| VIP用户 | 1000次/小时 |

超过限流返回：
```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试",
  "data": null
}
```

---

## 十三、接口测试

### Postman Collection

提供完整的Postman测试集合，包含所有接口的示例请求。

### Mock数据

开发环境提供Mock数据支持，方便前端开发调试。
