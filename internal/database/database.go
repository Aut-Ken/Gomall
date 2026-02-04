package database

import (
	"fmt"
	"time"

	"gomall/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB 全局数据库连接对象
var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	dbConfig := config.GetDatabase()

	// 构建 DSN 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.GetString("username"),
		dbConfig.GetString("password"),
		dbConfig.GetString("host"),
		dbConfig.GetString("port"),
		dbConfig.GetString("name"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置日志级别
		Logger: logger.Default.LogMode(logger.Info),
		// 命名策略（表名使用单数）
		NamingStrategy: namingStrategy{
			tablePrefix: "",
		},
	})

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 配置连接池
	maxIdleConns := dbConfig.GetInt("max_idle_conns")
	maxOpenConns := dbConfig.GetInt("max_open_conns")

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// // 自动迁移数据库表结构
	// if err := DB.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{}, &model.Stock{}); err != nil {
	// 	return fmt.Errorf("数据库迁移失败: %w", err)
	// }

	return nil
}

// Ping 检查数据库连接
func Ping() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}
	return sqlDB.Ping()
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// namingStrategy 自定义命名策略（使用单数表名）
type namingStrategy struct {
	tablePrefix string
}

func (n namingStrategy) TableName(table string) string {
	return n.tablePrefix + table
}

func (n namingStrategy) SchemaName(table string) string {
	return table
}

func (n namingStrategy) ColumnName(table, column string) string {
	return column
}

func (n namingStrategy) JoinTableName(joinTable string) string {
	return joinTable
}

func (n namingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return rel.Name + "_id"
}

func (n namingStrategy) CheckerName(table, column string) string {
	return table + "_" + column + "_chk"
}

func (n namingStrategy) IndexName(table, column string) string {
	return "idx_" + table + "_" + column
}

func (n namingStrategy) UniqueName(table, column string) string {
	return "uniq_" + table + "_" + column
}
