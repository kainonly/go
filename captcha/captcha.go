// Package captcha 提供基于 Redis 存储的验证码管理。
//
// 支持创建、验证和删除验证码，验证码会自动过期。
// 验证是原子性且一次性的（验证成功后验证码即被删除）。
//
// # Hertz 后端配置
//
//	// Initialize captcha with Redis client
//	cap := captcha.New(redisClient)
//
//	// Send captcha endpoint
//	h.POST("/captcha/send", func(ctx context.Context, c *app.RequestContext) {
//		email := c.Query("email")
//		code := help.Random(6) // Generate 6-digit code
//		cap.Create(ctx, "login:"+email, code, 5*time.Minute)
//		// Send code via email/SMS...
//		c.JSON(200, utils.H{"message": "sent"})
//	})
//
//	// Verify captcha endpoint
//	h.POST("/captcha/verify", func(ctx context.Context, c *app.RequestContext) {
//		email := c.Query("email")
//		code := c.Query("code")
//		if err := cap.Verify(ctx, "login:"+email, code); err != nil {
//			c.JSON(400, utils.H{"error": err.Error()})
//			return
//		}
//		c.JSON(200, utils.H{"message": "verified"})
//	})
//
// # Angular 前端配置
//
//	// Send captcha request
//	sendCaptcha(email: string) {
//	  return this.http.post('/captcha/send', null, { params: { email } });
//	}
//
//	// Verify captcha
//	verifyCaptcha(email: string, code: string) {
//	  return this.http.post('/captcha/verify', null, { params: { email, code } });
//	}
//
// # 安全注意事项
//
//   - 验证码在验证成功后即被删除（一次性使用）
//   - 使用合适的 TTL（如 5 分钟）缩小攻击窗口
//   - 建议配合限流防止暴力破解
//   - 使用 Redis 键前缀隔离不同类型的验证码
package captcha

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// captcha 函数返回的错误。
var (
	ErrNotExists   = errors.New("captcha: code does not exist or expired")
	ErrInvalidCode = errors.New("captcha: invalid code")
)

// Captcha 提供基于 Redis 存储的验证码管理。
type Captcha struct {
	// RDb 是用于存储验证码的 Redis 客户端。
	RDb *redis.Client
	// Prefix 是所有验证码键的键前缀（默认值："captcha"）。
	Prefix string
}

// New 使用给定的 Redis 客户端创建一个新的 Captcha 实例。
func New(rdb *redis.Client, options ...Option) *Captcha {
	x := &Captcha{
		RDb:    rdb,
		Prefix: "captcha",
	}
	for _, opt := range options {
		opt(x)
	}
	return x
}

// Option 是用于配置 Captcha 实例的函数。
type Option func(x *Captcha)

// SetPrefix 设置验证码键的 Redis 键前缀。
// 默认为 "captcha"，生成的键形如 "captcha:login:user@example.com"。
func SetPrefix(v string) Option {
	return func(x *Captcha) {
		x.Prefix = v
	}
}

// Key 生成验证码名称对应的完整 Redis 键。
// 格式："{prefix}:{name}"
func (x *Captcha) Key(name string) string {
	return fmt.Sprintf("%s:%s", x.Prefix, name)
}

// Create 以给定的名称和 TTL 存储验证码。
// 如果该名称下已存在验证码，将被覆盖。
// 成功时返回 "OK"。
func (x *Captcha) Create(ctx context.Context, name string, code string, ttl time.Duration) string {
	return x.RDb.Set(ctx, x.Key(name), code, ttl).Val()
}

// Exists 检查给定名称的验证码是否存在。
// 注意：此操作不会消耗验证码，实际验证请使用 Verify。
func (x *Captcha) Exists(ctx context.Context, name string) bool {
	return x.RDb.Exists(ctx, x.Key(name)).Val() != 0
}

// Verify 检查提供的验证码与存储的验证码是否匹配。
// 验证成功后，验证码会被自动删除（一次性使用）。
// 如果验证码不存在或已过期，返回 ErrNotExists。
// 如果验证码不匹配，返回 ErrInvalidCode。
func (x *Captcha) Verify(ctx context.Context, name string, code string) error {
	// 使用 GetDel 实现原子性的读取并删除操作
	result, err := x.RDb.GetDel(ctx, x.Key(name)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotExists
		}
		return err
	}
	if result != code {
		return ErrInvalidCode
	}
	return nil
}

// Delete 根据名称删除验证码。
// 返回删除的键数量（0 或 1）。
func (x *Captcha) Delete(ctx context.Context, name string) int64 {
	return x.RDb.Del(ctx, x.Key(name)).Val()
}
