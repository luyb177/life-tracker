package repo

import (
	"github.com/luyb177/life-tracker/backend/internal/repo/expense"
	"github.com/luyb177/life-tracker/backend/internal/repo/lifelog"
	"github.com/luyb177/life-tracker/backend/internal/repo/summary"
	"github.com/luyb177/life-tracker/backend/internal/repo/tag"
	"github.com/luyb177/life-tracker/backend/internal/repo/token"
	"github.com/luyb177/life-tracker/backend/internal/repo/user"
	"github.com/luyb177/life-tracker/backend/internal/repo/verify"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	User    user.Repository
	Verify  verify.Repository
	Token   *token.Repository
	Summary summary.Repository
	Expense expense.Repository
	LifeLog lifelog.Repository
	Tag     tag.Repository
	db      *gorm.DB
}

func NewRepositories(redisClient *redis.Client, db *gorm.DB) *Repositories {
	return &Repositories{
		User:    user.NewRepository(redisClient, db),
		Verify:  verify.NewVerifyRepo(redisClient),
		Token:   token.NewRepository(redisClient),
		Summary: summary.NewRepository(db),
		Expense: expense.NewRepository(db),
		LifeLog: lifelog.NewRepository(db),
		Tag:     tag.NewRepository(db),
		db:      db,
	}
}

func (r *Repositories) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
