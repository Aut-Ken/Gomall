-- ============================================
-- GoMall 数据库初始化脚本
-- 先删除所有旧数据，再重新创建
-- ============================================

-- 删除现有数据库（如果存在）
DROP DATABASE IF EXISTS Gomall;

-- 创建新数据库
CREATE DATABASE Gomall DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE Gomall;

-- ============================================
-- 创建表结构
-- ============================================

-- 用户表
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    email VARCHAR(100) NOT NULL COMMENT '邮箱',
    phone VARCHAR(20) DEFAULT '' COMMENT '手机号',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    UNIQUE KEY idx_username (username),
    UNIQUE KEY idx_email (email),
    KEY idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';

-- 商品表
CREATE TABLE products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '商品ID',
    name VARCHAR(200) NOT NULL COMMENT '商品名称',
    description TEXT COMMENT '商品描述',
    price DECIMAL(10, 2) NOT NULL COMMENT '商品价格',
    stock INT NOT NULL DEFAULT 0 COMMENT '库存数量',
    category VARCHAR(50) DEFAULT '' COMMENT '商品分类',
    image_url VARCHAR(500) DEFAULT '' COMMENT '商品图片URL',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '商品状态：1-上架，0-下架',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    KEY idx_status (status),
    KEY idx_category (category),
    KEY idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '商品表';

-- 订单表
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
    order_no VARCHAR(64) NOT NULL COMMENT '订单号',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    product_name VARCHAR(200) NOT NULL COMMENT '商品名称',
    quantity INT NOT NULL DEFAULT 1 COMMENT '购买数量',
    total_price DECIMAL(10, 2) NOT NULL COMMENT '订单总金额',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '订单状态：1-待支付，2-已支付，3-已发货，4-已完成，5-已取消',
    pay_type TINYINT NOT NULL DEFAULT 1 COMMENT '支付方式：1-支付宝，2-微信，3-银行卡',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',
    UNIQUE KEY idx_order_no (order_no),
    KEY idx_user_id (user_id),
    KEY idx_product_id (product_id),
    KEY idx_status (status),
    KEY idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '订单表';

-- 库存表
CREATE TABLE stocks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    total_stock INT NOT NULL DEFAULT 0 COMMENT '总库存',
    lock_stock INT NOT NULL DEFAULT 0 COMMENT '锁定库存',
    sold_stock INT NOT NULL DEFAULT 0 COMMENT '已售库存',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    UNIQUE KEY idx_product_id (product_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '库存表';

-- 购物车表
CREATE TABLE carts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    product_id BIGINT UNSIGNED NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    KEY idx_user_id (user_id),
    KEY idx_product_id (product_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '购物车表';

-- 优惠券表
CREATE TABLE coupons (
    id BIGINT UNSIGNED AUTO_INCREMENT COMMENT '优惠券ID',
    code VARCHAR(50) NOT NULL COMMENT '优惠券码',
    name VARCHAR(100) NOT NULL COMMENT '优惠券名称',
    discount_type TINYINT NOT NULL DEFAULT 1 COMMENT '优惠类型：1-金额折扣，2-百分比折扣',
    discount_amount DECIMAL(10, 2) NOT NULL COMMENT '折扣金额或比例',
    min_order_amount DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '最低订单金额',
    max_discount_amount DECIMAL(10, 2) DEFAULT NULL COMMENT '最大折扣金额',
    total_count INT NOT NULL DEFAULT 0 COMMENT '总发行量',
    used_count INT NOT NULL DEFAULT 0 COMMENT '已使用数量',
    valid_from DATETIME NOT NULL COMMENT '有效期开始',
    valid_until DATETIME NOT NULL COMMENT '有效期结束',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-有效，0-无效',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY idx_code (code)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '优惠券表';

-- 用户优惠券表
CREATE TABLE user_coupons (
    id BIGINT UNSIGNED AUTO_INCREMENT COMMENT '记录ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    coupon_id BIGINT UNSIGNED NOT NULL COMMENT '优惠券ID',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-未使用，2-已使用，3-已过期',
    used_at DATETIME(3) DEFAULT NULL COMMENT '使用时间',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    KEY idx_coupon_id (coupon_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户优惠券表';

-- 商品评论表
CREATE TABLE product_reviews (
    id BIGINT UNSIGNED AUTO_INCREMENT COMMENT '评论ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    rating TINYINT NOT NULL COMMENT '评分 1-5',
    content TEXT COMMENT '评论内容',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_product_id (product_id),
    KEY idx_user_id (user_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '商品评论表';

-- 地址表
CREATE TABLE addresses (
    id BIGINT UNSIGNED AUTO_INCREMENT COMMENT '地址ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    receiver_name VARCHAR(50) NOT NULL COMMENT '收货人姓名',
    receiver_phone VARCHAR(20) NOT NULL COMMENT '收货人电话',
    province VARCHAR(50) NOT NULL COMMENT '省份',
    city VARCHAR(50) NOT NULL COMMENT '城市',
    district VARCHAR(50) NOT NULL COMMENT '区县',
    detail VARCHAR(200) NOT NULL COMMENT '详细地址',
    is_default TINYINT NOT NULL DEFAULT 0 COMMENT '是否默认：1-是，0-否',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_user_id (user_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '地址表';

-- ============================================
-- 插入测试用户数据（密码都是 123456）
-- ============================================
INSERT INTO
    users (
        username,
        password,
        email,
        phone
    )
VALUES (
        'admin',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi',
        'admin@Gomall.com',
        '13800138000'
    ),
    (
        'testuser',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi',
        'test@Gomall.com',
        '13900139000'
    ),
    (
        'zhangsan',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi',
        'zhangsan@Gomall.com',
        '15000150000'
    ),
    (
        'lisi',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi',
        'lisi@Gomall.com',
        '15000150001'
    ),
    (
        'wangwu',
        '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi',
        'wangwu@Gomall.com',
        '15000150002'
    );

-- ============================================
-- 插入商品数据
-- ============================================
INSERT INTO
    products (
        name,
        description,
        price,
        stock,
        category,
        image_url,
        status
    )
VALUES (
        'iPhone 15 Pro',
        '苹果最新旗舰手机，A17 Pro芯片',
        8999.00,
        100,
        '手机',
        'https://example.com/iphone15.jpg',
        1
    ),
    (
        'iPhone 15',
        '苹果iPhone 15，A16仿生芯片',
        5999.00,
        150,
        '手机',
        'http://localhost:8080/photos/product_IPhone15.png',
        1
    ),
    (
        'MacBook Pro 14',
        'Apple M3 Pro芯片笔记本',
        14999.00,
        50,
        '电脑',
        'http://localhost:8080/photos/product_MacBook_Pro_14.png',
        1
    ),
    (
        'MacBook Air M3',
        'Apple M3芯片，超薄笔记本',
        7999.00,
        60,
        '电脑',
        'http://localhost:8080/photos/product_MacBook_Air_M3.png',
        1
    ),
    (
        'AirPods Pro 2',
        '第二代主动降噪耳机',
        1899.00,
        200,
        '耳机',
        'http://localhost:8080/photos/product_AirPods_Pro_2.png',
        1
    ),
    (
        'AirPods 3',
        '第三代AirPods',
        1399.00,
        300,
        '耳机',
        'http://localhost:8080/photos/product_AirPods_3.png',
        1
    ),
    (
        'iPad Air',
        '10.9英寸平板电脑',
        4799.00,
        80,
        '平板',
        'http://localhost:8080/photos/product_IPad_Air.png',
        1
    ),
    (
        'iPad Pro 11',
        'Apple iPad Pro 11英寸M2芯片',
        6799.00,
        70,
        '平板',
        'http://localhost:8080/photos/product_iPad_Pro_11.png',
        1
    ),
    (
        'Apple Watch Series 9',
        'Apple Watch Series 9',
        2999.00,
        120,
        '手表',
        'https://example.com/watch9.jpg',
        1
    ),
    (
        'Xiaomi 14 Pro',
        '小米14 Pro，徕卡影像',
        4999.00,
        200,
        '手机',
        'https://example.com/xiaomi14pro.jpg',
        1
    );

-- ============================================
-- 插入库存数据
-- ============================================
INSERT INTO
    stocks (
        product_id,
        total_stock,
        lock_stock,
        sold_stock
    )
VALUES (1, 100, 0, 0),
    (2, 150, 0, 0),
    (3, 50, 0, 0),
    (4, 60, 0, 0),
    (5, 200, 0, 0),
    (6, 300, 0, 0),
    (7, 80, 0, 0),
    (8, 70, 0, 0),
    (9, 120, 0, 0),
    (10, 200, 0, 0);

-- ============================================
-- 插入购物车数据
-- ============================================
INSERT INTO
    carts (user_id, product_id, quantity)
VALUES (1, 1, 1),
    (1, 3, 2),
    (2, 2, 1),
    (3, 5, 1);

-- ============================================
-- 插入优惠券数据
-- ============================================
INSERT INTO
    coupons (
        code,
        name,
        discount_type,
        discount_amount,
        min_order_amount,
        max_discount_amount,
        total_count,
        used_count,
        valid_from,
        valid_until,
        status
    )
VALUES (
        'WELCOME100',
        '新人100元券',
        1,
        100.00,
        500.00,
        100.00,
        10000,
        0,
        NOW(),
        DATE_ADD(NOW(), INTERVAL 90 DAY),
        1
    ),
    (
        'SAVE50',
        '满1000减50',
        1,
        50.00,
        1000.00,
        50.00,
        20000,
        0,
        NOW(),
        DATE_ADD(NOW(), INTERVAL 60 DAY),
        1
    ),
    (
        'SAVE20',
        '满500减20',
        1,
        20.00,
        500.00,
        20.00,
        30000,
        0,
        NOW(),
        DATE_ADD(NOW(), INTERVAL 30 DAY),
        1
    );

-- ============================================
-- 插入用户优惠券
-- ============================================
INSERT INTO
    user_coupons (user_id, coupon_id, status)
VALUES (1, 1, 1),
    (2, 1, 1),
    (3, 2, 1);

-- ============================================
-- 插入地址数据
-- ============================================
INSERT INTO
    addresses (
        user_id,
        receiver_name,
        receiver_phone,
        province,
        city,
        district,
        detail,
        is_default
    )
VALUES (
        1,
        '张三',
        '13800138000',
        '北京市',
        '北京市',
        '朝阳区',
        '建国路100号中央商务区C座1201',
        1
    ),
    (
        2,
        '李四',
        '13900139000',
        '上海市',
        '上海市',
        '浦东新区',
        '陆家嘴环路1000号金融大厦3502',
        1
    );

-- ============================================
-- 插入订单数据
-- ============================================
INSERT INTO
    orders (
        order_no,
        user_id,
        product_id,
        product_name,
        quantity,
        total_price,
        status,
        pay_type
    )
VALUES (
        'ORD202401010001',
        1,
        1,
        'iPhone 15 Pro',
        1,
        8999.00,
        4,
        1
    ),
    (
        'ORD202401010002',
        2,
        3,
        'MacBook Pro 14',
        1,
        14999.00,
        4,
        2
    ),
    (
        'ORD202401010003',
        3,
        5,
        'AirPods Pro 2',
        2,
        3798.00,
        4,
        1
    ),
    (
        'ORD202401010004',
        1,
        7,
        'iPad Air',
        1,
        4799.00,
        4,
        3
    );

-- ============================================
-- 插入商品评论
-- ============================================
INSERT INTO
    product_reviews (
        product_id,
        user_id,
        order_id,
        rating,
        content
    )
VALUES (
        1,
        1,
        1,
        5,
        '手机非常流畅，A17 Pro芯片太强大了，拍照效果也很棒！'
    ),
    (
        3,
        2,
        2,
        5,
        'MacBook Pro M3性能强劲，续航超赞，办公利器！'
    ),
    (
        5,
        3,
        3,
        5,
        'AirPods Pro 2降噪效果太好了，通透模式也很自然。'
    );

-- ============================================
-- 完成提示
-- ============================================
SELECT '数据库初始化完成！' AS message;

-- 登录信息：
-- 用户名: admin
-- 密码: 123456
UPDATE users
SET password = '$2a$10$aKBeAkAydW2EQenlXAb3EugcChMHhwv7UJr9.zzXSCHj/6BiKLymq'
WHERE
    username = 'admin';