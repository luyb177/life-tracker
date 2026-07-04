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
	Tags       string    `gorm:"type:varchar(500)"`
	OccurredAt time.Time `gorm:"index:idx_user_date;type:datetime"`
}

func (LifeLog) TableName() string {
	return "life_logs"
}
