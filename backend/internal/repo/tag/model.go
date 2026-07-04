package tag

import "time"

// Tag 全局标签（所有用户共享）
type Tag struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	Name      string `gorm:"uniqueIndex;type:varchar(50)"`
}

func (Tag) TableName() string {
	return "tags"
}

// LifeLogTag 生活记录-标签关联
type LifeLogTag struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	LifeLogID uint64 `gorm:"uniqueIndex:idx_lifelog_tag;type:bigint unsigned"`
	TagID     uint64 `gorm:"uniqueIndex:idx_lifelog_tag;index:idx_lifelog_tag_tag;type:bigint unsigned"`
}

func (LifeLogTag) TableName() string {
	return "life_log_tags"
}

// SummaryTag 总结-标签关联
type SummaryTag struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	SummaryID uint64 `gorm:"uniqueIndex:idx_summary_tag;type:bigint unsigned"`
	TagID     uint64 `gorm:"uniqueIndex:idx_summary_tag;index:idx_summary_tag_tag;type:bigint unsigned"`
}

func (SummaryTag) TableName() string {
	return "summary_tags"
}
