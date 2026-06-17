package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "embed"

	"github.com/redis/go-redis/v9"
)

const (
	// refresh:{user_id} -> {jti}
	refreshTokenKey = "refresh:%d"
)

//go:embed lua/rotate.lua
var rotateLua string

var rotateScript = redis.NewScript(rotateLua)

// Repository 管理 refresh token 的 Redis 存储，实现 token 轮换防重放。
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

// Rotate 原子性验证并轮换 JTI：旧 JTI 匹配则替换为新 JTI。
// 返回 true 表示轮换成功，旧 token 立即失效。
// 返回 false 表示旧 JTI 不匹配（可能已被使用过）。
func (r *Repository) Rotate(ctx context.Context, userID uint64, oldJTI, newJTI string, expire time.Duration) (bool, error) {
	key := fmt.Sprintf(refreshTokenKey, userID)

	val, err := rotateScript.Run(ctx, r.client, []string{key}, oldJTI, newJTI, int64(expire.Seconds())).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("rotate token: %w", err)
	}

	result, ok := val.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected rotate result: %T", val)
	}

	return result == 1, nil
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

// Revoke 删除用户的所有 refresh token（登出/改密码时使用）
func (r *Repository) Revoke(ctx context.Context, userID uint64) error {
	key := fmt.Sprintf(refreshTokenKey, userID)
	return r.client.Del(ctx, key).Err()
}
