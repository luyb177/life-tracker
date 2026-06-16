package summary

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Summary struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;type:bigint unsigned"`

	UserID            uint64 `gorm:"index:idx_user_period;type:bigint unsigned"`
	PeriodType        uint8  `gorm:"index:idx_user_period;type:tinyint unsigned"` // 1=日报, 2=周报, 3=月报, 4=年报
	PeriodStart       string `gorm:"index:idx_user_period;type:varchar(32)"`      // "2006-01-02" / "2006-W01" / "2006-01" / "2006"
	PeriodEnd         string `gorm:"type:varchar(32)"`
	Source            uint8  `gorm:"type:tinyint unsigned"` // 1=AI, 2=用户
	SummaryContent    string `gorm:"type:text"`
	SuggestionContent string `gorm:"type:text"`
	Status            uint8  `gorm:"type:tinyint unsigned;default:1"`
}

func (Summary) TableName() string {
	return "summaries"
}
