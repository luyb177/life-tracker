package expense

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// Category 支出分类
type Category struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;type:bigint unsigned"`

	UserID uint64 `gorm:"index:idx_user;type:bigint unsigned"`
	Name   string `gorm:"type:varchar(50)"`
	Type   uint8  `gorm:"type:tinyint unsigned"` // 1=系统默认, 2=用户自定义
}

func (Category) TableName() string {
	return "expense_categories"
}

// Log 支出记录
type Log struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;type:bigint unsigned"`

	UserID     uint64    `gorm:"index:idx_user_date;type:bigint unsigned"`
	CategoryID uint64    `gorm:"type:bigint unsigned"`
	Amount     float64   `gorm:"type:decimal(10,2)"`
	Note       string    `gorm:"type:varchar(255)"`
	OccurredAt time.Time `gorm:"index:idx_user_date;type:datetime"`
}

func (Log) TableName() string {
	return "expense_logs"
}
