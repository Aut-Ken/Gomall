package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`                 // 用户ID
	Username  string         `gorm:"uniqueIndex;size:50" json:"username"`   // 用户名
	Password  string         `gorm:"size:255" json:"-"`                     // 密码（不返回给前端）
	Email     string         `gorm:"uniqueIndex;size:100" json:"email"`     // 邮箱
	Phone     string         `gorm:"size:20" json:"phone"`                  // 手机号
	CreatedAt time.Time      `json:"created_at"`                             // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                             // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                         // 删除时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Product 商品模型
type Product struct {
	ID          uint           `gorm:"primarykey" json:"id"`                      // 商品ID
	Name        string         `gorm:"size:200;not null" json:"name"`             // 商品名称
	Description string         `gorm:"type:text" json:"description"`              // 商品描述
	Price       float64        `gorm:"not null;precision:10;scale:2" json:"price"` // 商品价格
	Stock       int            `gorm:"not null;default:0" json:"stock"`           // 库存数量
	Category    string         `gorm:"size:50" json:"category"`                   // 商品分类
	ImageURL    string         `gorm:"size:500" json:"image_url"`                 // 商品图片URL
	Status      int            `gorm:"default:1;comment:1-上架 0-下架" json:"status"` // 商品状态
	CreatedAt   time.Time      `json:"created_at"`                                // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                                // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                            // 删除时间
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// Order 订单模型
type Order struct {
	ID           uint           `gorm:"primarykey" json:"id"`                       // 订单ID
	OrderNo      string         `gorm:"uniqueIndex;size:64" json:"order_no"`        // 订单号
	UserID       uint           `gorm:"index;not null" json:"user_id"`              // 用户ID
	ProductID    uint           `gorm:"index;not null" json:"product_id"`           // 商品ID
	ProductName  string         `gorm:"size:200" json:"product_name"`               // 商品名称（冗余字段）
	Quantity     int            `gorm:"not null;default:1" json:"quantity"`         // 购买数量
	TotalPrice   float64        `gorm:"precision:10;scale:2" json:"total_price"`    // 订单总金额
	Status       int            `gorm:"default:1;comment:1-待支付 2-已支付 3-已发货 4-已完成 5-已取消" json:"status"` // 订单状态
	PayType      int            `gorm:"default:1;comment:1-支付宝 2-微信 3-银行卡" json:"pay_type"`            // 支付方式
	CreatedAt    time.Time      `json:"created_at"`                                 // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                                 // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                             // 删除时间
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// Stock 库存模型 (用于秒杀场景的库存扣减)
type Stock struct {
	ID         uint           `gorm:"primarykey" json:"id"`             // 记录ID
	ProductID  uint           `gorm:"uniqueIndex;not null" json:"product_id"` // 商品ID
	TotalStock int            `gorm:"not null;default:0" json:"total_stock"`  // 总库存
	LockStock  int            `gorm:"not null;default:0" json:"lock_stock"`   // 锁定库存
	SoldStock  int            `gorm:"not null;default:0" json:"sold_stock"`   // 已售库存
	CreatedAt  time.Time      `json:"created_at"`                        // 创建时间
	UpdatedAt  time.Time      `json:"updated_at"`                        // 更新时间
}

// TableName 指定表名
func (Stock) TableName() string {
	return "stocks"
}
