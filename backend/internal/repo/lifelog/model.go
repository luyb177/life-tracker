package lifelog

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// LifeLog 生活记录
type LifeLog struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;type:bigint unsigned"`

	UserID     uint64    `gorm:"index:idx_user_date;type:bigint unsigned"`
	Content    string    `gorm:"type:text"`
	OccurredAt time.Time `gorm:"index:idx_user_date;type:datetime"`

	LastUpdatedBy uint64    `gorm:"type:bigint unsigned;default:0"`
	LastUpdatedAt time.Time `gorm:"type:datetime(3)"`
}

func (LifeLog) TableName() string {
	return "life_logs"
}
