// Package captcha provides verification code management with Redis storage.
//
// It supports creating, verifying, and deleting captcha codes with automatic expiration.
// Verification is atomic and one-time (code is deleted after successful verification).
//
// # Hertz Backend Setup
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
// # Angular Frontend Setup
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
// # Security Notes
//
//   - Codes are deleted after successful verification (one-time use)
//   - Use appropriate TTL (e.g., 5 minutes) to limit attack window
//   - Consider rate limiting to prevent brute force attacks
//   - Use Redis key prefix to namespace different captcha types
package captcha

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Errors returned by captcha functions.
var (
	ErrNotExists   = errors.New("captcha: code does not exist or expired")
	ErrInvalidCode = errors.New("captcha: invalid code")
)

// Captcha provides verification code management with Redis storage.
type Captcha struct {
	// RDb is the Redis client for storing captcha codes.
	RDb *redis.Client
	// Prefix is the key prefix for all captcha keys (default: "captcha").
	Prefix string
}

// New creates a new Captcha instance with the given Redis client.
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

// Option is a function that configures a Captcha instance.
type Option func(x *Captcha)

// SetPrefix sets the Redis key prefix for captcha keys.
// Default is "captcha", resulting in keys like "captcha:login:user@example.com".
func SetPrefix(v string) Option {
	return func(x *Captcha) {
		x.Prefix = v
	}
}

// Key generates the full Redis key for a captcha name.
// Format: "{prefix}:{name}"
func (x *Captcha) Key(name string) string {
	return fmt.Sprintf("%s:%s", x.Prefix, name)
}

// Create stores a captcha code with the given name and TTL.
// If a code already exists for this name, it will be overwritten.
// Returns "OK" on success.
func (x *Captcha) Create(ctx context.Context, name string, code string, ttl time.Duration) string {
	return x.RDb.Set(ctx, x.Key(name), code, ttl).Val()
}

// Exists checks if a captcha code exists for the given name.
// Note: This does not consume the code. Use Verify for actual verification.
func (x *Captcha) Exists(ctx context.Context, name string) bool {
	return x.RDb.Exists(ctx, x.Key(name)).Val() != 0
}

// Verify checks if the provided code matches the stored captcha.
// On successful verification, the code is automatically deleted (one-time use).
// Returns ErrNotExists if code doesn't exist or expired.
// Returns ErrInvalidCode if code doesn't match.
func (x *Captcha) Verify(ctx context.Context, name string, code string) error {
	// Use GetDel for atomic get-and-delete operation
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

// Delete removes a captcha code by name.
// Returns the number of keys deleted (0 or 1).
func (x *Captcha) Delete(ctx context.Context, name string) int64 {
	return x.RDb.Del(ctx, x.Key(name)).Val()
}
