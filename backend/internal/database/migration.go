package database

import (
	"fmt"
	"time"

	"gomall/backend/internal/model"

	"gorm.io/gorm"
)

// Migration 数据库迁移接口
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

// BaseMigration 基础迁移结构
type BaseMigration struct{}

// TimestampMigration 带时间戳的迁移
type TimestampMigration struct {
	Timestamp time.Time
}

// Version 迁移版本
type Version struct {
	Version    string    `gorm:"primaryKey"`
	AppliedAt  time.Time `gorm:"autoUpdateTime"`
	Name       string
	IsRollback bool
}

// MigrationRunner 迁移运行器
type MigrationRunner struct {
	db     *gorm.DB
	migrations []Migration
}

// NewMigrationRunner 创建迁移运行器
func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{
		db: db,
		migrations: make([]Migration, 0),
	}
}

// Register 注册迁移
func (r *MigrationRunner) Register(m Migration) {
	r.migrations = append(r.migrations, m)
}

// Run 执行所有未应用的迁移
func (r *MigrationRunner) Run() error {
	// 自动迁移（使用GORM的AutoMigrate作为后备）
	if err := r.autoMigrate(); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	// 执行自定义迁移
	for _, m := range r.migrations {
		if err := m.Up(r.db); err != nil {
			return fmt.Errorf("迁移执行失败: %w", err)
		}
	}

	return nil
}

// autoMigrate 自动迁移（保留向后兼容）
func (r *MigrationRunner) autoMigrate() error {
	return r.db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.Order{},
		&model.Cart{},
		&model.Stock{},
	)
}

// Rollback 回滚最后一个迁移
func (r *MigrationRunner) Rollback() error {
	if len(r.migrations) == 0 {
		return nil
	}

	lastMigration := r.migrations[len(r.migrations)-1]
	return lastMigration.Down(r.db)
}

// ==================== 自定义迁移 ====================

// CreateUsersTableMigration 创建用户表
type CreateUsersTableMigration struct{}

func (m *CreateUsersTableMigration) Up(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		phone VARCHAR(20),
		role TINYINT DEFAULT 1 COMMENT '1:普通用户, 2:管理员',
		status TINYINT DEFAULT 1 COMMENT '0:禁用, 1:正常',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return db.Exec(sql).Error
}

func (m *CreateUsersTableMigration) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS users").Error
}

// CreateProductsTableMigration 创建商品表
type CreateProductsTableMigration struct{}

func (m *CreateProductsTableMigration) Up(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS products (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(200) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		stock INT DEFAULT 0,
		category VARCHAR(100),
		image_url VARCHAR(500),
		status TINYINT DEFAULT 1 COMMENT '0:下架, 1:上架',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_category (category),
		INDEX idx_status (status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return db.Exec(sql).Error
}

func (m *CreateProductsTableMigration) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS products").Error
}

// CreateOrdersTableMigration 创建订单表
type CreateOrdersTableMigration struct{}

func (m *CreateOrdersTableMigration) Up(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS orders (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		order_no VARCHAR(32) NOT NULL UNIQUE,
		user_id BIGINT UNSIGNED NOT NULL,
		product_id BIGINT UNSIGNED NOT NULL,
		product_name VARCHAR(200) NOT NULL,
		quantity INT NOT NULL,
		total_price DECIMAL(10, 2) NOT NULL,
		status TINYINT DEFAULT 0 COMMENT '0:处理中, 1:待支付, 2:已支付, 3:已发货, 4:已完成, 5:已取消',
		pay_type TINYINT DEFAULT 0 COMMENT '0:未支付, 1:支付宝, 2:微信',
		pay_time DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_user_id (user_id),
		INDEX idx_order_no (order_no),
		INDEX idx_status (status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return db.Exec(sql).Error
}

func (m *CreateOrdersTableMigration) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS orders").Error
}

// CreateCartsTableMigration 创建购物车表
type CreateCartsTableMigration struct{}

func (m *CreateCartsTableMigration) Up(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS carts (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		user_id BIGINT UNSIGNED NOT NULL,
		product_id BIGINT UNSIGNED NOT NULL,
		quantity INT NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE INDEX idx_user_product (user_id, product_id),
		INDEX idx_user_id (user_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return db.Exec(sql).Error
}

func (m *CreateCartsTableMigration) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS carts").Error
}

// RunMigrations 运行所有迁移
func RunMigrations(db *gorm.DB) error {
	runner := NewMigrationRunner(db)

	// 注册迁移（按依赖顺序）
	runner.Register(&CreateUsersTableMigration{})
	runner.Register(&CreateProductsTableMigration{})
	runner.Register(&CreateOrdersTableMigration{})
	runner.Register(&CreateCartsTableMigration{})

	return runner.Run()
}
