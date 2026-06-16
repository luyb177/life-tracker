package user

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        uint64 `gorm:"primarykey;type:bigint unsigned auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;type:bigint unsigned"`

	Avatar   string `gorm:"type:varchar(512)"`
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"uniqueIndex;type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
}

func (User) TableName() string {
	return "users"
}
