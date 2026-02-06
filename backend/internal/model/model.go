package model

/**
 * Model 数据模型层
 *
 * 本文件定义了电商系统所有核心业务实体的数据结构。
 * 使用 GORM 框架的 ORM 特性，将 Go 结构体与数据库表映射。
 *
 * 设计原则：
 * 1. 每个结构体对应一张数据库表
 * 2. 使用 gorm 标签定义字段映射关系
 * 3. 使用 json 标签定义 API 响应格式
 * 4. 支持软删除（通过 gorm.DeletedAt）
 * 5. 字段命名使用蛇形命名法（数据库）和驼峰命名法（Go）
 *
 * 表命名规则：
 * - users: 用户表
 * - products: 商品表
 * - orders: 订单表
 * - stocks: 库存表
 * - carts: 购物车表
 */

import (
	"time"

	"gorm.io/gorm"
)

/**
 * User 用户模型
 *
 * 存储用户账号信息，支持普通用户和管理员两种角色。
 *
 * 角色说明：
 * - Role = 1: 普通用户，可以浏览商品、下单购买
 * - Role = 2: 管理员，可以管理商品、查看统计数据
 *
 * 密码安全：
 * - Password 字段使用 bcrypt 加密存储
 * - json 标签设为 "-" 确保密码不会泄露到 API 响应中
 *
 * 软删除：
 * - 使用 gorm.DeletedAt 实现软删除
 * - 删除用户时不会物理删除数据，而是设置 deleted_at 字段
 * - 查询时会自动过滤已删除的记录
 */
type User struct {
	// ID 用户唯一标识，自增主键
	ID uint `gorm:"column:id;primarykey" json:"id"`

	// Username 用户名，唯一索引，长度50
	// 用于登录和显示
	Username string `gorm:"column:username;uniqueIndex;size:50" json:"username"`

	// Password 加密后的密码
	// 使用 bcrypt 算法加密，长度255
	// json 标签 "-" 确保不会序列化到 JSON
	Password string `gorm:"column:password;size:255" json:"-"`

	// Email 用户邮箱，唯一索引，长度100
	// 用于接收通知、可选的登录方式
	Email string `gorm:"column:email;uniqueIndex;size:100" json:"email"`

	// Phone 用户手机号，长度20
	// 可选字段，用于接收短信通知
	Phone string `gorm:"column:phone;size:20" json:"phone"`

	// Role 用户角色，默认普通用户(1)
	// 1: 普通用户, 2: 管理员
	Role int `gorm:"column:role;default:1" json:"role"`

	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// DeletedAt 软删除时间
	// GORM 会自动管理这个字段
	// 当调用 Delete 时，设置这个字段而不是真正删除
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

/**
 * 用户角色常量定义
 */
const (
	RoleUser  = 1 // 普通用户，可以进行正常的购物操作
	RoleAdmin = 2 // 管理员，拥有管理商品的权限
)

/**
 * TableName 指定 User 结构体对应的数据库表名
 *
 * GORM 默认会将结构体名转为复数形式（User -> users）
 * 但我们显式定义以确保表名正确
 *
 * 返回值：
 *   "users" - 对应数据库中的 users 表
 */
func (User) TableName() string {
	return "users"
}

/**
 * Product 商品模型
 *
 * 存储商品的基础信息，包括名称、价格、库存、分类等。
 *
 * 状态说明：
 * - Status = 1: 上架状态，用户可以看到并购买
 * - Status = 0: 下架状态，用户无法购买
 *
 * 价格精度：
 * - 使用 DECIMAL(10,2) 类型存储价格
 * - 精度为10位，小数点后2位
 */
type Product struct {
	// ID 商品唯一标识，自增主键
	ID uint `gorm:"column:id;primarykey" json:"id"`

	// Name 商品名称，必填，长度200
	Name string `gorm:"column:name;size:200;not null" json:"name"`

	// Description 商品描述
	// 使用 TEXT 类型，支持长文本
	Description string `gorm:"column:description;type:text" json:"description"`

	// Price 商品价格，必填
	// DECIMAL(10,2) 存储，确保金额精度
	Price float64 `gorm:"column:price;not null;precision:10;scale:2" json:"price"`

	// Stock 商品库存数量，必填，默认0
	// 下单时会扣减库存
	Stock int `gorm:"column:stock;not null;default:0" json:"stock"`

	// Category 商品分类，长度50
	// 用于商品分类浏览和筛选
	Category string `gorm:"column:category;size:50" json:"category"`

	// ImageURL 商品图片URL，长度500
	// 存储商品主图的访问地址
	ImageURL string `gorm:"column:image_url;size:500" json:"image_url"`

	// Status 商品状态，默认上架(1)
	// 1: 上架, 0: 下架
	Status int `gorm:"column:status;default:1" json:"status"`

	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// DeletedAt 软删除时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

/**
 * TableName 指定 Product 结构体对应的数据库表名
 */
func (Product) TableName() string {
	return "products"
}

/**
 * Order 订单模型
 *
 * 存储用户订单信息，包括订单号、商品信息、价格、状态等。
 *
 * 订单状态（Status）：
 * - 0: 处理中 - 订单刚创建，等待系统处理
 * - 1: 待支付 - 订单已创建，等待用户支付
 * - 2: 已支付 - 用户已完成支付
 * - 3: 已发货 - 商家已发货
 * - 4: 已完成 - 订单交易成功
 * - 5: 已取消 - 订单被取消（用户主动或超时）
 *
 * 支付类型（PayType）：
 * - 1: 支付宝
 * - 2: 微信支付
 *
 * 订单号生成规则：
 * - 格式：ORD + 时间戳(YYYYMMDDHHmmss) + 4位随机数
 * - 例如：ORD202401011200001234
 */
type Order struct {
	// ID 订单唯一标识，自增主键
	ID uint `gorm:"column:id;primarykey" json:"id"`

	// OrderNo 订单号，唯一索引，长度64
	// 格式：ORD + 时间戳 + 随机数
	OrderNo string `gorm:"column:order_no;uniqueIndex;size:64" json:"order_no"`

	// UserID 下单用户ID，外键关联 users 表
	UserID uint `gorm:"column:user_id;index;not null" json:"user_id"`

	// ProductID 购买商品ID，外键关联 products 表
	ProductID uint `gorm:"column:product_id;index;not null" json:"product_id"`

	// ProductName 下单时的商品名称（冗余存储）
	// 商品信息变更时不影响历史订单
	ProductName string `gorm:"column:product_name;size:200" json:"product_name"`

	// Quantity 购买数量，默认1
	Quantity int `gorm:"column:quantity;not null;default:1" json:"quantity"`

	// TotalPrice 订单总金额
	// = 商品单价 × 数量
	TotalPrice float64 `gorm:"column:total_price;precision:10;scale:2" json:"total_price"`

	// Status 订单状态
	// 0: 处理中, 1: 待支付, 2: 已支付, 3: 已发货, 4: 已完成, 5: 已取消
	Status int `gorm:"column:status;default:1" json:"status"`

	// PayType 支付类型
	// 1: 支付宝, 2: 微信
	PayType int `gorm:"column:pay_type;default:1" json:"pay_type"`

	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// DeletedAt 软删除时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

/**
 * TableName 指定 Order 结构体对应的数据库表名
 */
func (Order) TableName() string {
	return "orders"
}

/**
 * Stock 库存模型
 *
 * 存储商品的库存信息，与 Product 表关联。
 * 采用分离设计，将库存信息独立存储，方便库存管理。
 *
 * 库存类型说明：
 * - TotalStock: 总库存数量
 * - LockStock: 锁定库存（已下单但未支付）
 * - SoldStock: 已售库存（已支付）
 *
 * 可用库存计算公式：
 * - 可用库存 = TotalStock - LockStock - SoldStock
 */
type Stock struct {
	// ID 库存记录唯一标识，自增主键
	ID uint `gorm:"column:id;primarykey" json:"id"`

	// ProductID 商品ID，唯一索引
	// 每个商品只有一条库存记录
	ProductID uint `gorm:"column:product_id;uniqueIndex;not null" json:"product_id"`

	// TotalStock 总库存数量
	TotalStock int `gorm:"column:total_stock;not null;default:0" json:"total_stock"`

	// LockStock 锁定库存数量
	// 下单但未支付时，库存会被锁定
	LockStock int `gorm:"column:lock_stock;not null;default:0" json:"lock_stock"`

	// SoldStock 已售库存数量
	// 支付完成后计入已售
	SoldStock int `gorm:"column:sold_stock;not null;default:0" json:"sold_stock"`

	// CreatedAt 创建时间
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

/**
 * TableName 指定 Stock 结构体对应的数据库表名
 */
func (Stock) TableName() string {
	return "stocks"
}

/**
 * Cart 购物车模型
 *
 * 存储用户的购物车商品信息。
 *
 * 设计特点：
 * - 每个用户对每个商品只有一条购物车记录
 * - 数量字段记录该商品的购买数量
 * - 冗余存储商品名称和图片，方便购物车列表展示
 *
 * 唯一索引：
 * (user_id, product_id) 组合唯一
 */
type Cart struct {
	// ID 购物车记录唯一标识，自增主键
	ID uint `gorm:"column:id;primarykey" json:"id"`

	// UserID 用户ID，索引
	UserID uint `gorm:"column:user_id;index;not null" json:"user_id"`

	// ProductID 商品ID，索引
	// 与 UserID 组合成联合唯一索引
	ProductID uint `gorm:"column:product_id;index;not null" json:"product_id"`

	// Quantity 商品数量，默认1
	// 如果商品已存在，数量累加
	Quantity int `gorm:"column:quantity;not null;default:1" json:"quantity"`

	// CreatedAt 添加到购物车的时间
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// UpdatedAt 最后更新时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// DeletedAt 软删除时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

/**
 * TableName 指定 Cart 结构体对应的数据库表名
 */
func (Cart) TableName() string {
	return "carts"
}
