-- =============================================
-- 校园闲置物品交换平台 - 数据库初始化脚本
-- 数据库：MySQL 8.0+
-- 字符集：utf8mb4
-- =============================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS campus_swap DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE campus_swap;

-- =============================================
-- 1. 用户表
-- =============================================
CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(100) NOT NULL COMMENT '密码（bcrypt加密）',
    `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `gender` TINYINT DEFAULT NULL COMMENT '性别（0未知 1男 2女）',
    `student_id` VARCHAR(50) DEFAULT NULL COMMENT '学号',
    `school` VARCHAR(100) DEFAULT NULL COMMENT '学校名称',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（0禁用 1正常 2冻结）',
    `credit_score` INT NOT NULL DEFAULT 100 COMMENT '信誉分（0-100）',
    `create_time` BIGINT NOT NULL COMMENT '创建时间（时间戳）',
    `update_time` BIGINT DEFAULT NULL COMMENT '更新时间（时间戳）',
    `last_login_time` BIGINT DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(50) DEFAULT NULL COMMENT '最后登录IP',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    UNIQUE KEY `uk_phone` (`phone`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =============================================
-- 2. 分类表
-- =============================================
CREATE TABLE IF NOT EXISTS `category` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类ID（0表示顶级）',
    `level` TINYINT NOT NULL DEFAULT 1 COMMENT '层级（1/2/3）',
    `icon` VARCHAR(255) DEFAULT NULL COMMENT '图标URL',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序号',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（0禁用 1启用）',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- =============================================
-- 3. 商品表
-- =============================================
CREATE TABLE IF NOT EXISTS `goods` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `title` VARCHAR(100) NOT NULL COMMENT '商品标题',
    `description` TEXT DEFAULT NULL COMMENT '商品描述',
    `category_id` BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '发布者ID',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型（1售卖 2交换 3均可）',
    `price` DECIMAL(10,2) DEFAULT NULL COMMENT '售卖价格（元）',
    `original_price` DECIMAL(10,2) DEFAULT NULL COMMENT '原价（元）',
    `images` JSON DEFAULT NULL COMMENT '商品图片数组',
    `condition` TINYINT NOT NULL DEFAULT 1 COMMENT '成色（1全新 2九成新 3八成新 4七成新）',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态（0下架 1在售 2已售 3已交换）',
    `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览次数',
    `favorite_count` INT NOT NULL DEFAULT 0 COMMENT '收藏次数',
    `exchange_count` INT NOT NULL DEFAULT 0 COMMENT '交换次数',
    `tags` VARCHAR(200) DEFAULT NULL COMMENT '标签（逗号分隔）',
    `location` VARCHAR(100) DEFAULT NULL COMMENT '交易地点',
    `latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
    `longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    `update_time` BIGINT DEFAULT NULL COMMENT '更新时间',
    `sold_time` BIGINT DEFAULT NULL COMMENT '售出/交换时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`),
    KEY `idx_price` (`price`),
    FULLTEXT KEY `idx_title_desc` (`title`, `description`),
    CONSTRAINT `fk_goods_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_goods_category` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- =============================================
-- 4. 交换表
-- =============================================
CREATE TABLE IF NOT EXISTS `exchange` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '交换ID',
    `initiator_id` BIGINT UNSIGNED NOT NULL COMMENT '发起人ID',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标用户ID',
    `goods_id` BIGINT UNSIGNED NOT NULL COMMENT '目标商品ID',
    `my_goods_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '我的商品ID（用于物物交换）',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型（1购买 2交换）',
    `message` TEXT DEFAULT NULL COMMENT '附加留言',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态（0待处理 1已接受 2已拒绝 3已取消 4已完成）',
    `reject_reason` VARCHAR(200) DEFAULT NULL COMMENT '拒绝原因',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    `update_time` BIGINT DEFAULT NULL COMMENT '更新时间',
    `complete_time` BIGINT DEFAULT NULL COMMENT '完成时间',
    PRIMARY KEY (`id`),
    KEY `idx_initiator_id` (`initiator_id`),
    KEY `idx_target_id` (`target_id`),
    KEY `idx_goods_id` (`goods_id`),
    KEY `idx_status` (`status`),
    KEY `idx_create_time` (`create_time`),
    CONSTRAINT `fk_exchange_initiator` FOREIGN KEY (`initiator_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_exchange_target` FOREIGN KEY (`target_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_exchange_goods` FOREIGN KEY (`goods_id`) REFERENCES `goods` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交换表';

-- =============================================
-- 5. 聊天消息表
-- =============================================
CREATE TABLE IF NOT EXISTS `message` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `sender_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    `receiver_id` BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型（1文本 2图片 3商品卡片）',
    `goods_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联商品ID（类型3时）',
    `is_read` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读（0未读 1已读）',
    `create_time` BIGINT NOT NULL COMMENT '发送时间',
    PRIMARY KEY (`id`),
    KEY `idx_sender_receiver` (`sender_id`, `receiver_id`),
    KEY `idx_create_time` (`create_time`),
    KEY `idx_is_read` (`is_read`),
    CONSTRAINT `fk_message_sender` FOREIGN KEY (`sender_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_message_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';

-- =============================================
-- 6. 收藏表
-- =============================================
CREATE TABLE IF NOT EXISTS `favorite` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `goods_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_goods` (`user_id`, `goods_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_create_time` (`create_time`),
    CONSTRAINT `fk_favorite_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_favorite_goods` FOREIGN KEY (`goods_id`) REFERENCES `goods` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- =============================================
-- 7. 评论表
-- =============================================
CREATE TABLE IF NOT EXISTS `comment` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `goods_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '评论者ID',
    `target_user_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '被评论者ID（卖家ID）',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `rating` TINYINT NOT NULL DEFAULT 5 COMMENT '评分（1-5）',
    `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID（0表示主评论）',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（0隐藏 1显示）',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_goods_id` (`goods_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_create_time` (`create_time`),
    CONSTRAINT `fk_comment_goods` FOREIGN KEY (`goods_id`) REFERENCES `goods` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comment_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- =============================================
-- 8. 举报表
-- =============================================
CREATE TABLE IF NOT EXISTS `report` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '举报ID',
    `reporter_id` BIGINT UNSIGNED NOT NULL COMMENT '举报人ID',
    `target_type` TINYINT NOT NULL COMMENT '目标类型（1商品 2用户 3评论）',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    `reason` VARCHAR(200) NOT NULL COMMENT '举报原因',
    `description` TEXT DEFAULT NULL COMMENT '详细说明',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态（0待处理 1已处理 2已驳回）',
    `handler_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '处理人ID',
    `handle_result` TEXT DEFAULT NULL COMMENT '处理结果',
    `handle_time` BIGINT DEFAULT NULL COMMENT '处理时间',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_reporter_id` (`reporter_id`),
    KEY `idx_target` (`target_type`, `target_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `fk_report_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='举报表';

-- =============================================
-- 9. 系统通知表
-- =============================================
CREATE TABLE IF NOT EXISTS `notification` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '接收用户ID',
    `type` TINYINT NOT NULL COMMENT '类型（1系统 2交换 3评论 4点赞）',
    `title` VARCHAR(100) NOT NULL COMMENT '通知标题',
    `content` TEXT DEFAULT NULL COMMENT '通知内容',
    `link` VARCHAR(255) DEFAULT NULL COMMENT '跳转链接',
    `is_read` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已读',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_is_read` (`is_read`),
    KEY `idx_create_time` (`create_time`),
    CONSTRAINT `fk_notification_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统通知表';

-- =============================================
-- 10. 管理员表
-- =============================================
CREATE TABLE IF NOT EXISTS `admin` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(100) NOT NULL COMMENT '密码（bcrypt加密）',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `role` TINYINT NOT NULL DEFAULT 2 COMMENT '角色（1超级管理员 2普通管理员）',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（0禁用 1正常）',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    `update_time` BIGINT DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- =============================================
-- 11. 操作日志表
-- =============================================
CREATE TABLE IF NOT EXISTS `operation_log` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `admin_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '操作管理员ID',
    `module` VARCHAR(50) NOT NULL COMMENT '操作模块',
    `action` VARCHAR(50) NOT NULL COMMENT '操作动作',
    `method` VARCHAR(100) NOT NULL COMMENT '请求方法',
    `params` TEXT DEFAULT NULL COMMENT '请求参数',
    `ip` VARCHAR(50) DEFAULT NULL COMMENT '操作IP',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（0失败 1成功）',
    `error_msg` TEXT DEFAULT NULL COMMENT '错误信息',
    `create_time` BIGINT NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_admin_id` (`admin_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- =============================================
-- 初始化数据
-- =============================================

-- 插入默认管理员账号（用户名：admin，密码：123456）
-- 注意：生产环境请立即修改密码！
INSERT INTO `admin` (`username`, `password`, `nickname`, `role`, `status`, `create_time`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '超级管理员', 1, 1, UNIX_TIMESTAMP() * 1000);

-- 插入测试用户（密码：123456）
INSERT INTO `user` (`username`, `password`, `nickname`, `email`, `school`, `create_time`) VALUES
('testuser', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '测试用户', 'test@example.com', 'XX大学', UNIX_TIMESTAMP() * 1000);

-- 插入商品分类
INSERT INTO `category` (`name`, `parent_id`, `level`, `sort`, `status`, `create_time`) VALUES
-- 一级分类
('数码产品', 0, 1, 1, 1, UNIX_TIMESTAMP() * 1000),
('书籍资料', 0, 1, 2, 1, UNIX_TIMESTAMP() * 1000),
('服饰鞋包', 0, 1, 3, 1, UNIX_TIMESTAMP() * 1000),
('生活用品', 0, 1, 4, 1, UNIX_TIMESTAMP() * 1000),
('运动户外', 0, 1, 5, 1, UNIX_TIMESTAMP() * 1000),
('其他', 0, 1, 6, 1, UNIX_TIMESTAMP() * 1000);

-- 二级分类（数码产品）
SET @parent_id = LAST_INSERT_ID();
INSERT INTO `category` (`name`, `parent_id`, `level`, `sort`, `status`, `create_time`) VALUES
('手机', 1, 2, 1, 1, UNIX_TIMESTAMP() * 1000),
('电脑', 1, 2, 2, 1, UNIX_TIMESTAMP() * 1000),
('平板', 1, 2, 3, 1, UNIX_TIMESTAMP() * 1000),
('相机', 1, 2, 4, 1, UNIX_TIMESTAMP() * 1000),
('配件', 1, 2, 5, 1, UNIX_TIMESTAMP() * 1000);

-- 二级分类（书籍资料）
INSERT INTO `category` (`name`, `parent_id`, `level`, `sort`, `status`, `create_time`) VALUES
('教材', 2, 2, 1, 1, UNIX_TIMESTAMP() * 1000),
('文学', 2, 2, 2, 1, UNIX_TIMESTAMP() * 1000),
('考试资料', 2, 2, 3, 1, UNIX_TIMESTAMP() * 1000),
('杂志期刊', 2, 2, 4, 1, UNIX_TIMESTAMP() * 1000);

-- =============================================
-- 创建索引优化查询性能
-- =============================================

-- 商品表复合索引（用于列表筛选）
CREATE INDEX idx_goods_category_status_time ON `goods`(`category_id`, `status`, `create_time`);

-- 交换表复合索引（用于我的交换列表）
CREATE INDEX idx_exchange_initiator_status ON `exchange`(`initiator_id`, `status`, `create_time`);
CREATE INDEX idx_exchange_target_status ON `exchange`(`target_id`, `status`, `create_time`);

-- 评论表复合索引（用于商品评论列表）
CREATE INDEX idx_comment_goods_time ON `comment`(`goods_id`, `create_time`);

-- =============================================
-- 完成
-- =============================================

-- 显示创建的表
SHOW TABLES;

-- 显示用户数量
SELECT COUNT(*) as user_count FROM `user`;

-- 显示商品数量
SELECT COUNT(*) as goods_count FROM `goods`;

-- 显示分类数量
SELECT COUNT(*) as category_count FROM `category`;
