// Package locker 提供基于 Redis 存储的限流与尝试计数功能。
//
// 适用于登录尝试限制、API 限流，或任何需要追踪并限制时间窗口内操作次数的场景。
//
// # Hertz 后端配置
//
//	// Initialize locker with Redis client
//	lock := locker.New(redisClient)
//
//	// Login endpoint with attempt limiting
//	h.POST("/auth/login", func(ctx context.Context, c *app.RequestContext) {
//		username := c.Query("username")
//		lockKey := "login:" + username
//
//		// Check if already locked (max 5 attempts)
//		if err := lock.Check(ctx, lockKey, 5); err != nil {
//			if errors.Is(err, locker.ErrLocked) {
//				c.JSON(429, utils.H{"error": "too many attempts, try again later"})
//				return
//			}
//		}
//
//		// Validate credentials...
//		if !validCredentials {
//			// Increment attempt counter (5 minute window)
//			lock.Increment(ctx, lockKey, 5*time.Minute)
//			c.JSON(401, utils.H{"error": "invalid credentials"})
//			return
//		}
//
//		// Success - clear the counter
//		lock.Delete(ctx, lockKey)
//		c.JSON(200, utils.H{"token": token})
//	})
//
// # 工作原理
//
//   - Increment：原子递增计数器，首次调用时设置 TTL
//   - Check：计数器达到上限时返回 ErrLocked
//   - Delete：清空计数器（如登录成功后）
//
// # 安全注意事项
//
//   - 根据安全需求设置合适的 TTL
//   - 登录限制建议结合 IP 与用户名
//   - 计数器在 TTL 后自动过期，无需手动清理
package locker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// locker 函数返回的错误。
var (
	ErrNotExists = errors.New("locker: counter does not exist")
	ErrLocked    = errors.New("locker: limit exceeded")
)

// Locker 提供基于 Redis 存储的限流与尝试计数。
type Locker struct {
	// RDb 是存储计数器的 Redis 客户端。
	RDb *redis.Client
	// Prefix 是所有 locker 键的前缀（默认："locker"）。
	Prefix string
}

// New 使用给定的 Redis 客户端创建新的 Locker 实例。
func New(rdb *redis.Client, options ...Option) *Locker {
	x := &Locker{
		RDb:    rdb,
		Prefix: "locker",
	}
	for _, opt := range options {
		opt(x)
	}
	return x
}

// Option 是用于配置 Locker 实例的函数。
type Option func(x *Locker)

// SetPrefix 设置 locker 键的 Redis 前缀。
// 默认为 "locker"，生成形如 "locker:login:user@example.com" 的键。
func SetPrefix(v string) Option {
	return func(x *Locker) {
		x.Prefix = v
	}
}

// Key 根据名称生成完整的 Redis 键。
// 格式："{prefix}:{name}"
func (x *Locker) Key(name string) string {
	return fmt.Sprintf("%s:%s", x.Prefix, name)
}

// incrWithExpire 是一段 Lua 脚本，原子地递增计数器，
// 且仅在键为新键时设置过期时间。
var incrWithExpire = redis.NewScript(`
local current = redis.call('INCR', KEYS[1])
if current == 1 then
    redis.call('PEXPIRE', KEYS[1], ARGV[1])
end
return current
`)

// Increment 原子递增指定名称的计数器。
// 首次调用时设置计数器的 TTL，TTL 窗口内的后续调用只递增、不重置 TTL。
// 返回递增后的当前计数。
func (x *Locker) Increment(ctx context.Context, name string, ttl time.Duration) (int64, error) {
	result, err := incrWithExpire.Run(ctx, x.RDb, []string{x.Key(name)}, ttl.Milliseconds()).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Check 校验计数器是否超过允许的最大值。
// 计数器小于 max 或不存在时返回 nil；达到上限返回 ErrLocked。
func (x *Locker) Check(ctx context.Context, name string, max int64) error {
	result, err := x.RDb.Get(ctx, x.Key(name)).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // 键不存在，未锁定
		}
		return err
	}
	if result >= max {
		return ErrLocked
	}
	return nil
}

// Get 返回指定名称的当前计数值。
// 计数器不存在时返回 0。
func (x *Locker) Get(ctx context.Context, name string) (int64, error) {
	result, err := x.RDb.Get(ctx, x.Key(name)).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

// Delete 删除指定名称的计数器，返回删除的键数（0 或 1）。
// 通常在认证成功后调用以重置计数器。
func (x *Locker) Delete(ctx context.Context, name string) int64 {
	return x.RDb.Del(ctx, x.Key(name)).Val()
}
