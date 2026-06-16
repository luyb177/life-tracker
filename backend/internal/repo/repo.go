package repo

import (
	"github.com/luyb177/life-tracker/backend/internal/repo/token"
	"github.com/luyb177/life-tracker/backend/internal/repo/user"
	"github.com/luyb177/life-tracker/backend/internal/repo/verify"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	User   user.Repository
	Verify verify.Repository
	Token  *token.Repository
	db     *gorm.DB
}

func NewRepositories(redisClient *redis.Client, db *gorm.DB) *Repositories {
	return &Repositories{
		User:   user.NewRepository(redisClient, db),
		Verify: verify.NewVerifyRepo(redisClient),
		Token:  token.NewRepository(redisClient),
		db:     db,
	}
}

func (r *Repositories) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
