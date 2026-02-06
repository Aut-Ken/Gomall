-- GoMall 数据库初始化脚本
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS Gomall DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

SET NAMES utf8mb4;
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

-- 优惠券表
CREATE TABLE IF NOT EXISTS `coupons` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '优惠券ID',
    `code` varchar(50) NOT NULL COMMENT '优惠券码',
    `name` varchar(100) NOT NULL COMMENT '优惠券名称',
    `discount_type` tinyint NOT NULL DEFAULT 1 COMMENT '优惠类型：1-金额折扣，2-百分比折扣',
    `discount_amount` decimal(10,2) NOT NULL COMMENT '折扣金额或比例',
    `min_order_amount` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '最低订单金额',
    `max_discount_amount` decimal(10,2) DEFAULT NULL COMMENT '最大折扣金额',
    `total_count` int NOT NULL DEFAULT 0 COMMENT '总发行量',
    `used_count` int NOT NULL DEFAULT 0 COMMENT '已使用数量',
    `valid_from` datetime NOT NULL COMMENT '有效期开始',
    `valid_until` datetime NOT NULL COMMENT '有效期结束',
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1-有效，0-无效',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券表';

-- 用户优惠券表
CREATE TABLE IF NOT EXISTS `user_coupons` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `coupon_id` bigint unsigned NOT NULL COMMENT '优惠券ID',
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1-未使用，2-已使用，3-已过期',
    `used_at` datetime(3) DEFAULT NULL COMMENT '使用时间',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_coupon_id` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户优惠券表';

-- 商品评论表
CREATE TABLE IF NOT EXISTS `product_reviews` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `rating` tinyint NOT NULL COMMENT '评分 1-5',
    `content` text COMMENT '评论内容',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品评论表';

-- 地址表
CREATE TABLE IF NOT EXISTS `addresses` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '地址ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `receiver_name` varchar(50) NOT NULL COMMENT '收货人姓名',
    `receiver_phone` varchar(20) NOT NULL COMMENT '收货人电话',
    `province` varchar(50) NOT NULL COMMENT '省份',
    `city` varchar(50) NOT NULL COMMENT '城市',
    `district` varchar(50) NOT NULL COMMENT '区县',
    `detail` varchar(200) NOT NULL COMMENT '详细地址',
    `is_default` tinyint NOT NULL DEFAULT 0 COMMENT '是否默认：1-是，0-否',
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地址表';

-- 插入更多测试用户
INSERT INTO `users` (`username`, `password`, `email`, `phone`) VALUES
('zhangsan', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'zhangsan@Gomall.com', '15000150000'),
('lisi', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'lisi@Gomall.com', '15000150001'),
('wangwu', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'wangwu@Gomall.com', '15000150002'),
('zhaoliu', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'zhaoliu@Gomall.com', '15000150003'),
('sunqi', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'sunqi@Gomall.com', '15000150004'),
('zhouba', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'zhouba@Gomall.com', '15000150005'),
('wujiu', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'wujiu@Gomall.com', '15000150006'),
('zhengshi', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'zhengshi@Gomall.com', '15000150007'),
('liushiping', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'liushiping@Gomall.com', '15000150008'),
('songjiu', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'songjiu@Gomall.com', '15000150009');

-- 插入更多商品数据
INSERT INTO `products` (`name`, `description`, `price`, `stock`, `category`, `image_url`, `status`) VALUES
-- 手机类
('iPhone 15', '苹果iPhone 15，A16仿生芯片', 5999.00, 150, '手机', 'https://example.com/iphone15.jpg', 1),
('iPhone 15 Pro Max', '苹果顶级旗舰，钛金属机身', 9999.00, 80, '手机', 'https://example.com/iphone15promax.jpg', 1),
('Samsung Galaxy S24', '三星Galaxy S24，AI手机', 5499.00, 120, '手机', 'https://example.com/s24.jpg', 1),
('Xiaomi 14 Pro', '小米14 Pro，徕卡影像', 4999.00, 200, '手机', 'https://example.com/xiaomi14pro.jpg', 1),
('OPPO Find X7', 'OPPO Find X7，哈苏影像', 3999.00, 100, '手机', 'https://example.com/findx7.jpg', 1),
-- 电脑类
('MacBook Air M3', 'Apple M3芯片，超薄笔记本', 7999.00, 60, '电脑', 'https://example.com/mba_m3.jpg', 1),
('Dell XPS 15', '戴尔XPS 15，性能本', 12999.00, 40, '电脑', 'https://example.com/xps15.jpg', 1),
('ThinkPad X1 Carbon', '联想ThinkPad X1 Carbon商务本', 9999.00, 50, '电脑', 'https://example.com/x1carbon.jpg', 1),
('HP Spectre x360', '惠普Spectre x360翻转本', 8999.00, 30, '电脑', 'https://example.com/spectre.jpg', 1),
-- 耳机类
('AirPods 3', '第三代AirPods', 1399.00, 300, '耳机', 'https://example.com/airpods3.jpg', 1),
('Sony WH-1000XM5', '索尼WH-1000XM5降噪耳机', 2999.00, 150, '耳机', 'https://example.com/sonyxm5.jpg', 1),
('Bose QC Ultra', 'Bose QuietComfort Ultra', 2599.00, 100, '耳机', 'https://example.com/boseqcu.jpg', 1),
('Huawei FreeBuds Pro 3', '华为FreeBuds Pro 3', 1499.00, 200, '耳机', 'https://example.com/freepro3.jpg', 1),
-- 平板类
('iPad Pro 11', 'Apple iPad Pro 11英寸M2芯片', 6799.00, 70, '平板', 'https://example.com/ipadpro11.jpg', 1),
('iPad mini 7', 'Apple iPad mini 7.9英寸', 3999.00, 90, '平板', 'https://example.com/ipadmini.jpg', 1),
('Samsung Tab S9', '三星Tab S9平板', 5999.00, 60, '平板', 'https://example.com/tabs9.jpg', 1),
('Xiaomi Pad 6S Pro', '小米Pad 6S Pro', 3299.00, 100, '平板', 'https://example.com/xiaomipad.jpg', 1),
-- 手表类
('Apple Watch Series 9', 'Apple Watch Series 9', 2999.00, 120, '手表', 'https://example.com/watch9.jpg', 1),
('Apple Watch Ultra 2', 'Apple Watch Ultra 2极限版', 6499.00, 40, '手表', 'https://example.com/watchultra.jpg', 1),
('Samsung Watch 6', '三星Galaxy Watch 6', 1999.00, 80, '手表', 'https://example.com/watch6.jpg', 1),
('Huawei Watch GT4', '华为Watch GT4', 1499.00, 150, '手表', 'https://example.com/gt4.jpg', 1),
-- 智能家居类
('Xiaomi Smart Home Hub', '小米智能家居中枢', 399.00, 500, '智能家居', 'https://example.com/smarthub.jpg', 1),
('Philips Hue Starter', '飞利浦Hue入门套装', 799.00, 200, '智能家居', 'https://example.com/hue.jpg', 1),
('Xiaomi Smart Lock', '小米智能门锁', 1299.00, 300, '智能家居', 'https://example.com/smartlock.jpg', 1),
('Ecovacs Deebot X2', '科沃斯X2扫拖机器人', 4999.00, 80, '智能家居', 'https://example.com/deebot.jpg', 1),
-- 相机类
('Sony A7 IV', '索尼A7 IV全画幅微单', 16999.00, 25, '相机', 'https://example.com/a7iv.jpg', 1),
('Canon R6 Mark II', '佳能R6 Mark II', 17999.00, 20, '相机', 'https://example.com/r6ii.jpg', 1),
('DJI Mini 4 Pro', '大疆Mini 4 Pro无人机', 7388.00, 50, '相机', 'https://example.com/dji.jpg', 1),
-- 游戏类
('Nintendo Switch OLED', '任天堂Switch OLED主机', 2599.00, 100, '游戏', 'https://example.com/switch.jpg', 1),
('PlayStation 5', '索尼PS5光驱版', 4299.00, 60, '游戏', 'https://example.com/ps5.jpg', 1),
('Steam Deck', 'Valve Steam Deck掌机', 3999.00, 40, '游戏', 'https://example.com/steamdeck.jpg', 1),
('Xbox Series X', '微软Xbox Series X', 3899.00, 50, '游戏', 'https://example.com/xbox.jpg', 1),
-- 配件类
('Anker PowerBank 20000', '安克20000mAh充电宝', 299.00, 500, '配件', 'https://example.com/powerbank.jpg', 1),
('Apple USB-C Cable', 'Apple USB-C编织线', 149.00, 1000, '配件', 'https://example.com/usbc.jpg', 1),
('Samsung 990 Pro 1TB', '三星990 Pro 1TB SSD', 999.00, 300, '配件', 'https://example.com/ssd.jpg', 1),
('Logitech MX Master 3S', '罗技MX Master 3S鼠标', 899.00, 200, '配件', 'https://example.com/mxmaster.jpg', 1),
('Keychron K2', 'Keychron K2机械键盘', 648.00, 250, '配件', 'https://example.com/keychron.jpg', 1);

-- 插入新商品的库存数据
INSERT INTO `stocks` (`product_id`, `total_stock`, `lock_stock`, `sold_stock`) VALUES
-- iPhone系列
(5, 150, 5, 20),
(6, 80, 3, 10),
(7, 120, 8, 35),
(8, 200, 10, 45),
(9, 100, 5, 25),
-- 电脑
(10, 60, 2, 8),
(11, 40, 2, 12),
(12, 50, 3, 18),
(13, 30, 1, 7),
-- 耳机
(14, 300, 15, 80),
(15, 150, 8, 42),
(16, 100, 5, 28),
(17, 200, 12, 55),
-- 平板
(18, 70, 3, 18),
(19, 90, 4, 22),
(20, 60, 2, 15),
(21, 100, 5, 30),
-- 手表
(22, 120, 6, 35),
(23, 40, 2, 12),
(24, 80, 4, 22),
(25, 150, 8, 45),
-- 智能家居
(26, 500, 20, 120),
(27, 200, 10, 55),
(28, 300, 15, 88),
(29, 80, 4, 22),
-- 相机
(30, 25, 1, 8),
(31, 20, 1, 6),
(32, 50, 2, 14),
-- 游戏
(33, 100, 5, 32),
(34, 60, 3, 18),
(35, 40, 2, 12),
(36, 50, 2, 15),
-- 配件
(37, 500, 25, 150),
(38, 1000, 50, 320),
(39, 300, 15, 98),
(40, 200, 10, 65),
(41, 250, 12, 78);

-- 插入优惠券数据
INSERT INTO `coupons` (`code`, `name`, `discount_type`, `discount_amount`, `min_order_amount`, `max_discount_amount`, `total_count`, `used_count`, `valid_from`, `valid_until`, `status`) VALUES
('WELCOME100', '新人100元券', 1, 100.00, 500.00, 100.00, 10000, 2341, NOW(), DATE_ADD(NOW(), INTERVAL 90 DAY), 1),
('WELCOME50', '新人50元券', 1, 50.00, 200.00, 50.00, 20000, 5678, NOW(), DATE_ADD(NOW(), INTERVAL 90 DAY), 1),
('SAVE10', '满200减10', 1, 10.00, 200.00, 10.00, 50000, 12345, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 1),
('SAVE20', '满500减20', 1, 20.00, 500.00, 20.00, 30000, 8765, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 1),
('SAVE50', '满1000减50', 1, 50.00, 1000.00, 50.00, 20000, 4532, NOW(), DATE_ADD(NOW(), INTERVAL 60 DAY), 1),
('SAVE100', '满2000减100', 1, 100.00, 2000.00, 100.00, 10000, 2341, NOW(), DATE_ADD(NOW(), INTERVAL 60 DAY), 1),
('VIP95', 'VIP95折', 2, 5.00, 0.00, 500.00, 5000, 1234, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY), 1),
('VIP90', 'VIP9折', 2, 10.00, 1000.00, 1000.00, 2000, 567, NOW(), DATE_ADD(NOW(), INTERVAL 365 DAY), 1),
('SUMMER20', '夏日20元券', 1, 20.00, 300.00, 20.00, 15000, 7890, NOW(), DATE_ADD(NOW(), INTERVAL 15 DAY), 1),
('NEWYEAR100', '新年100元券', 1, 100.00, 800.00, 100.00, 8000, 1234, NOW(), DATE_ADD(NOW(), INTERVAL 45 DAY), 1);

-- 插入用户优惠券
INSERT INTO `user_coupons` (`user_id`, `coupon_id`, `status`, `used_at`) VALUES
(1, 1, 2, NOW() - INTERVAL 5 DAY),
(1, 2, 1, NULL),
(1, 4, 2, NOW() - INTERVAL 3 DAY),
(2, 1, 2, NOW() - INTERVAL 10 DAY),
(2, 3, 2, NOW() - INTERVAL 8 DAY),
(2, 5, 1, NULL),
(3, 1, 1, NULL),
(3, 6, 1, NULL),
(4, 2, 2, NOW() - INTERVAL 2 DAY),
(5, 3, 1, NULL);

-- 插入商品评论
INSERT INTO `product_reviews` (`product_id`, `user_id`, `order_id`, `rating`, `content`) VALUES
(1, 1, 1, 5, '手机非常流畅，A17 Pro芯片太强大了，拍照效果也很棒！'),
(1, 2, 2, 5, 'iPhone 15 Pro 体验很好，系统流畅，生态完善。'),
(1, 3, 3, 4, '总体不错，就是充电速度还可以再快一点。'),
(2, 1, 4, 5, 'MacBook Pro M3性能强劲，续航超赞，办公利器！'),
(2, 4, 5, 5, '屏幕效果惊艳，音质也很好，推荐入手。'),
(3, 2, 6, 5, 'AirPods Pro 2降噪效果太好了，通透模式也很自然。'),
(3, 5, 7, 4, '音质不错，佩戴舒适度很好。'),
(4, 3, 8, 5, 'iPad Air性价比很高，M1芯片性能过剩，用来画画很棒。'),
(5, 1, 9, 5, 'iPhone 15颜色很美，拍照效果很好，喜欢！'),
(8, 2, 10, 5, '小米14 Pro拍照效果超出预期，徕卡色彩很有质感。'),
(15, 3, 11, 5, 'Sony WH-1000XM5降噪效果一流，佩戴也很舒适。'),
(18, 4, 12, 5, 'iPad Pro 11屏幕太棒了，M2芯片运行流畅。'),
(22, 5, 13, 5, 'Apple Watch Series 9功能丰富，健康监测很实用。'),
(26, 1, 14, 5, '小米智能家居中枢很好用，联动很方便。'),
(30, 2, 15, 4, '索尼A7 IV对焦很快，高感表现优秀，视频功能也很强。'),
(33, 3, 16, 5, 'Switch OLED屏幕效果很好，游戏体验很棒。'),
(37, 4, 17, 5, 'Anker充电宝容量大，充电速度快，非常实用。');

-- 插入用户地址
INSERT INTO `addresses` (`user_id`, `receiver_name`, `receiver_phone`, `province`, `city`, `district`, `detail`, `is_default`) VALUES
(1, '张三', '15000150000', '北京市', '北京市', '朝阳区', '建国路100号中央商务区C座1201', 1),
(1, '张三', '15000150000', '上海市', '上海市', '浦东新区', '陆家嘴环路1000号金融大厦3502', 0),
(2, '李四', '15000150001', '广东省', '深圳市', '南山区', '科技园路18号创维大厦801', 1),
(2, '李四', '15000150001', '浙江省', '杭州市', '西湖区', '文三路478号创新产业园A区305', 0),
(3, '王五', '15000150002', '江苏省', '南京市', '鼓楼区', '中山北路256号商贸中心1808', 1),
(4, '赵六', '15000150003', '四川省', '成都市', '高新区', '天府大道中段500号科技大厦2201', 1),
(5, '孙七', '15000150004', '湖北省', '武汉市', '洪山区', '珞珈山路128号创业孵化中心1205', 1),
(6, '周八', '15000150005', '陕西省', '西安市', '雁塔区', '科技路168号电子信息城A座903', 1),
(7, '吴九', '15000150006', '天津市', '天津市', '滨海新区', '经济技术开发区第五大街99号', 1),
(8, '郑十', '15000150007', '重庆市', '重庆市', '渝北区', '财富大道88号金融中心3506', 1);

-- 插入购物车数据
INSERT INTO `carts` (`user_id`, `product_id`, `quantity`, `created_at`, `updated_at`) VALUES
(1, 1, 1, NOW() - INTERVAL 7 DAY, NOW()),
(1, 3, 2, NOW() - INTERVAL 5 DAY, NOW()),
(1, 5, 1, NOW() - INTERVAL 3 DAY, NOW()),
(2, 2, 1, NOW() - INTERVAL 6 DAY, NOW()),
(2, 8, 1, NOW() - INTERVAL 4 DAY, NOW()),
(3, 4, 1, NOW() - INTERVAL 8 DAY, NOW()),
(3, 10, 2, NOW() - INTERVAL 2 DAY, NOW()),
(4, 6, 1, NOW() - INTERVAL 5 DAY, NOW()),
(4, 15, 1, NOW() - INTERVAL 1 DAY, NOW()),
(5, 7, 1, NOW() - INTERVAL 4 DAY, NOW()),
(5, 18, 1, NOW() - INTERVAL 3 DAY, NOW()),
(6, 9, 1, NOW() - INTERVAL 6 DAY, NOW()),
(7, 12, 1, NOW() - INTERVAL 2 DAY, NOW()),
(8, 14, 1, NOW() - INTERVAL 5 DAY, NOW()),
(9, 20, 1, NOW() - INTERVAL 4 DAY, NOW()),
(10, 22, 1, NOW() - INTERVAL 3 DAY, NOW());

-- 插入更多订单数据
INSERT INTO `orders` (`order_no`, `user_id`, `product_id`, `product_name`, `quantity`, `total_price`, `status`, `pay_type`, `created_at`, `updated_at`) VALUES
('ORD202401180001', 1, 1, 'iPhone 15 Pro', 1, 8999.00, 4, 1, NOW() - INTERVAL 30 DAY, NOW() - INTERVAL 25 DAY),
('ORD202401180002', 2, 2, 'MacBook Pro 14', 1, 14999.00, 4, 2, NOW() - INTERVAL 28 DAY, NOW() - INTERVAL 23 DAY),
('ORD202401180003', 3, 3, 'AirPods Pro 2', 2, 3798.00, 4, 1, NOW() - INTERVAL 25 DAY, NOW() - INTERVAL 20 DAY),
('ORD202401180004', 1, 4, 'iPad Air', 1, 4799.00, 4, 3, NOW() - INTERVAL 22 DAY, NOW() - INTERVAL 18 DAY),
('ORD202401180005', 4, 5, 'iPhone 15', 1, 5999.00, 4, 1, NOW() - INTERVAL 20 DAY, NOW() - INTERVAL 15 DAY),
('ORD202401180006', 2, 8, 'Xiaomi 14 Pro', 1, 4999.00, 4, 2, NOW() - INTERVAL 18 DAY, NOW() - INTERVAL 13 DAY),
('ORD202401180007', 5, 10, 'MacBook Air M3', 1, 7999.00, 4, 1, NOW() - INTERVAL 15 DAY, NOW() - INTERVAL 10 DAY),
('ORD202401180008', 3, 15, 'Sony WH-1000XM5', 1, 2999.00, 4, 2, NOW() - INTERVAL 12 DAY, NOW() - INTERVAL 8 DAY),
('ORD202401180009', 6, 18, 'iPad Pro 11', 1, 6799.00, 4, 1, NOW() - INTERVAL 10 DAY, NOW() - INTERVAL 5 DAY),
('ORD202401180010', 7, 22, 'Apple Watch Series 9', 1, 2999.00, 4, 3, NOW() - INTERVAL 8 DAY, NOW() - INTERVAL 3 DAY),
('ORD202401180011', 1, 26, 'Xiaomi Smart Home Hub', 2, 798.00, 4, 1, NOW() - INTERVAL 6 DAY, NOW() - INTERVAL 2 DAY),
('ORD202401180012', 4, 30, 'Sony A7 IV', 1, 16999.00, 4, 2, NOW() - INTERVAL 5 DAY, NOW() - INTERVAL 1 DAY),
('ORD202401180013', 8, 33, 'Nintendo Switch OLED', 1, 2599.00, 3, 1, NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 1 DAY),
('ORD202401180014', 2, 37, 'Anker PowerBank 20000', 2, 598.00, 2, 2, NOW() - INTERVAL 1 DAY, NOW()),
('ORD202401180015', 5, 40, 'DJI Mini 4 Pro', 1, 7388.00, 1, 1, NOW() - INTERVAL 2 HOUR, NOW()),
('ORD202401180016', 9, 41, 'Keychron K2', 1, 648.00, 1, 1, NOW() - INTERVAL 1 HOUR, NOW()),
('ORD202401180017', 10, 39, 'Samsung 990 Pro 1TB', 2, 1998.00, 1, 3, NOW(), NOW());
