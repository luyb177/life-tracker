package verify

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	SetCode(ctx context.Context, meta *Meta, code string, expire time.Duration) error
	VerifyCode(ctx context.Context, meta *Meta, code string) (bool, error)
}

type repo struct {
	client *redis.Client
}

func NewVerifyRepo(client *redis.Client) Repository {
	return &repo{
		client: client,
	}
}
