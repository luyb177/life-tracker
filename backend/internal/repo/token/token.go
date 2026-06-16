package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// refresh:{user_id} -> {jti}
	refreshTokenKey = "refresh:%d"
)

// Repository 管理 refresh token 的 Redis 存储，实现 token 轮换防重放。
// 每次刷新后旧 JTI 被覆盖，旧 refresh token 立即失效。
type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{client: client}
}

// Store 存储用户当前有效的 refresh token JTI，设置过期时间
func (r *Repository) Store(ctx context.Context, userID uint64, jti string, expire time.Duration) error {
	key := fmt.Sprintf(refreshTokenKey, userID)
	return r.client.Set(ctx, key, jti, expire).Err()
}

// Validate 校验 refresh token 的 JTI 是否与 Redis 中存储的一致
func (r *Repository) Validate(ctx context.Context, userID uint64, jti string) (bool, error) {
	key := fmt.Sprintf(refreshTokenKey, userID)
	stored, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == jti, nil
}

// Revoke 删除用户的所有 refresh token（登出时使用）
func (r *Repository) Revoke(ctx context.Context, userID uint64) error {
	key := fmt.Sprintf(refreshTokenKey, userID)
	return r.client.Del(ctx, key).Err()
}
