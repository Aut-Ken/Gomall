-- GoMall 数据库初始化脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS Gomall DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE Gomall;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '密码',
    `email` varchar(100) NOT NULL COMMENT '邮箱',
    `phone` varchar(20) DEFAULT '' COMMENT '手机号',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 商品表
CREATE TABLE IF NOT EXISTS `products` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `name` varchar(200) NOT NULL COMMENT '商品名称',
    `description` text COMMENT '商品描述',
    `price` decimal(10,2) NOT NULL COMMENT '商品价格',
    `stock` int NOT NULL DEFAULT 0 COMMENT '库存数量',
    `category` varchar(50) DEFAULT '' COMMENT '商品分类',
    `image_url` varchar(500) DEFAULT '' COMMENT '商品图片URL',
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '商品状态：1-上架，0-下架',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`),
    KEY `idx_category` (`category`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 订单表
CREATE TABLE IF NOT EXISTS `orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` varchar(64) NOT NULL COMMENT '订单号',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `product_name` varchar(200) NOT NULL COMMENT '商品名称（冗余）',
    `quantity` int NOT NULL DEFAULT 1 COMMENT '购买数量',
    `total_price` decimal(10,2) NOT NULL COMMENT '订单总金额',
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '订单状态：1-待支付，2-已支付，3-已发货，4-已完成，5-已取消',
    `pay_type` tinyint NOT NULL DEFAULT 1 COMMENT '支付方式：1-支付宝，2-微信，3-银行卡',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 库存表
CREATE TABLE IF NOT EXISTS `stocks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `total_stock` int NOT NULL DEFAULT 0 COMMENT '总库存',
    `lock_stock` int NOT NULL DEFAULT 0 COMMENT '锁定库存',
    `sold_stock` int NOT NULL DEFAULT 0 COMMENT '已售库存',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存表';

-- 插入测试数据
INSERT INTO `users` (`username`, `password`, `email`, `phone`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'admin@Gomall.com', '13800138000'),
('testuser', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'test@Gomall.com', '13900139000');

INSERT INTO `products` (`name`, `description`, `price`, `stock`, `category`, `image_url`, `status`) VALUES
('iPhone 15 Pro', '苹果最新旗舰手机，A17 Pro芯片', 8999.00, 100, '手机', 'https://example.com/iphone15.jpg', 1),
('MacBook Pro 14', 'Apple M3 Pro芯片笔记本', 14999.00, 50, '电脑', 'https://example.com/macbook.jpg', 1),
('AirPods Pro 2', '第二代主动降噪耳机', 1899.00, 200, '耳机', 'https://example.com/airpods.jpg', 1),
('iPad Air', '10.9英寸平板电脑', 4799.00, 80, '平板', 'https://example.com/ipad.jpg', 1);

INSERT INTO `stocks` (`product_id`, `total_stock`, `lock_stock`, `sold_stock`) VALUES
(1, 100, 0, 0),
(2, 50, 0, 0),
(3, 200, 0, 0),
(4, 80, 0, 0);

CREATE TABLE carts (
      id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
      user_id BIGINT UNSIGNED NOT NULL,
      product_id BIGINT UNSIGNED NOT NULL,
      quantity INT NOT NULL DEFAULT 1,
      created_at DATETIME,
      updated_at DATETIME,
      deleted_at DATETIME,
      INDEX idx_user_id (user_id),
      INDEX idx_product_id (product_id)
  );
