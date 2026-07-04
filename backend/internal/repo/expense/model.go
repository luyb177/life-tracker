package expense

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// Category 支出分类
// user_id=0 表示系统默认分类（全局可见）
type Category struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;uniqueIndex:idx_user_name_del;type:bigint unsigned"`

	UserID uint64 `gorm:"uniqueIndex:idx_user_name_del;type:bigint unsigned"` // 0=系统默认
	Name   string `gorm:"uniqueIndex:idx_user_name_del;type:varchar(50)"`
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

	UserID     uint64     `gorm:"index:idx_user_date;type:bigint unsigned"`
	CategoryID uint64     `gorm:"index;type:bigint unsigned"`
	Amount     int64      `gorm:"type:bigint"` // 单位：分
	Note       string     `gorm:"type:varchar(255)"`
	Location   string     `gorm:"type:varchar(255)"`
	OccurredAt time.Time  `gorm:"index:idx_user_date;type:datetime"`
	Status     uint8      `gorm:"type:tinyint unsigned;default:0"` // 0=正常, 1=已退款
	RefundedAt *time.Time `gorm:"type:datetime"`                   // 退款时间
}

func (Log) TableName() string {
	return "expense_logs"
}
