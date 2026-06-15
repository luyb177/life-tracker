package repo

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	db *gorm.DB
}

// NewRepositories creates Repositories with both Redis and MySQL.
func NewRepositories(redisClient *redis.Client, db *gorm.DB) *Repositories {
	return &Repositories{
		db: db,
	}
}

// Transaction starts a MySQL transaction.
func (r *Repositories) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
