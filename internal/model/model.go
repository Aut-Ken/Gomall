package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	// 显式指定每一个字段对应的数据库列名 (column:xxx)
	ID        uint           `gorm:"column:id;primarykey" json:"id"`
	Username  string         `gorm:"column:username;uniqueIndex;size:50" json:"username"`
	Password  string         `gorm:"column:password;size:255" json:"-"`
	Email     string         `gorm:"column:email;uniqueIndex;size:100" json:"email"`
	Phone     string         `gorm:"column:phone;size:20" json:"phone"`
	Role      int            `gorm:"column:role;default:1" json:"role"` // 1: 普通用户, 2: 管理员
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// 用户角色常量
const (
	RoleUser  = 1 // 普通用户
	RoleAdmin = 2 // 管理员
)

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Product 商品模型
type Product struct {
	ID          uint           `gorm:"column:id;primarykey" json:"id"`
	Name        string         `gorm:"column:name;size:200;not null" json:"name"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Price       float64        `gorm:"column:price;not null;precision:10;scale:2" json:"price"`
	Stock       int            `gorm:"column:stock;not null;default:0" json:"stock"`
	Category    string         `gorm:"column:category;size:50" json:"category"`
	ImageURL    string         `gorm:"column:image_url;size:500" json:"image_url"`
	Status      int            `gorm:"column:status;default:1" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// Order 订单模型
type Order struct {
	ID          uint           `gorm:"column:id;primarykey" json:"id"`
	OrderNo     string         `gorm:"column:order_no;uniqueIndex;size:64" json:"order_no"`
	UserID      uint           `gorm:"column:user_id;index;not null" json:"user_id"`
	ProductID   uint           `gorm:"column:product_id;index;not null" json:"product_id"`
	ProductName string         `gorm:"column:product_name;size:200" json:"product_name"`
	Quantity    int            `gorm:"column:quantity;not null;default:1" json:"quantity"`
	TotalPrice  float64        `gorm:"column:total_price;precision:10;scale:2" json:"total_price"`
	Status      int            `gorm:"column:status;default:1" json:"status"`
	PayType     int            `gorm:"column:pay_type;default:1" json:"pay_type"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// Stock 库存模型
type Stock struct {
	ID         uint      `gorm:"column:id;primarykey" json:"id"`
	ProductID  uint      `gorm:"column:product_id;uniqueIndex;not null" json:"product_id"`
	TotalStock int       `gorm:"column:total_stock;not null;default:0" json:"total_stock"`
	LockStock  int       `gorm:"column:lock_stock;not null;default:0" json:"lock_stock"`
	SoldStock  int       `gorm:"column:sold_stock;not null;default:0" json:"sold_stock"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Stock) TableName() string {
	return "stocks"
}
