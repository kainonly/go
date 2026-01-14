// Package locker provides rate limiting and attempt counting with Redis storage.
//
// It's useful for limiting login attempts, API rate limiting, or any scenario
// where you need to track and limit the number of operations within a time window.
//
// # Hertz Backend Setup
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
// # How It Works
//
//   - Increment: Atomically increments counter, sets TTL on first call
//   - Check: Returns ErrLocked if counter >= max
//   - Delete: Clears the counter (e.g., after successful login)
//
// # Security Notes
//
//   - Use appropriate TTL based on your security requirements
//   - Consider using IP + username combination for login limiting
//   - The counter auto-expires after TTL, no manual cleanup needed
package locker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Errors returned by locker functions.
var (
	ErrNotExists = errors.New("locker: counter does not exist")
	ErrLocked    = errors.New("locker: limit exceeded")
)

// Locker provides rate limiting and attempt counting with Redis storage.
type Locker struct {
	// RDb is the Redis client for storing counters.
	RDb *redis.Client
	// Prefix is the key prefix for all locker keys (default: "locker").
	Prefix string
}

// New creates a new Locker instance with the given Redis client.
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

// Option is a function that configures a Locker instance.
type Option func(x *Locker)

// SetPrefix sets the Redis key prefix for locker keys.
// Default is "locker", resulting in keys like "locker:login:user@example.com".
func SetPrefix(v string) Option {
	return func(x *Locker) {
		x.Prefix = v
	}
}

// Key generates the full Redis key for a locker name.
// Format: "{prefix}:{name}"
func (x *Locker) Key(name string) string {
	return fmt.Sprintf("%s:%s", x.Prefix, name)
}

// incrWithExpire is a Lua script that atomically increments a counter
// and sets expiration only if the key is new.
var incrWithExpire = redis.NewScript(`
local current = redis.call('INCR', KEYS[1])
if current == 1 then
    redis.call('PEXPIRE', KEYS[1], ARGV[1])
end
return current
`)

// Increment atomically increments the counter for the given name.
// On the first call, it sets the TTL for the counter.
// Subsequent calls within the TTL window only increment without resetting TTL.
// Returns the current count after incrementing.
func (x *Locker) Increment(ctx context.Context, name string, ttl time.Duration) (int64, error) {
	result, err := incrWithExpire.Run(ctx, x.RDb, []string{x.Key(name)}, ttl.Milliseconds()).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Check verifies if the counter has exceeded the maximum allowed value.
// Returns nil if counter < max or counter doesn't exist.
// Returns ErrLocked if counter >= max.
func (x *Locker) Check(ctx context.Context, name string, max int64) error {
	result, err := x.RDb.Get(ctx, x.Key(name)).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // Key doesn't exist, not locked
		}
		return err
	}
	if result >= max {
		return ErrLocked
	}
	return nil
}

// Get returns the current counter value for the given name.
// Returns 0 if the counter doesn't exist.
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

// Delete removes the counter for the given name.
// Returns the number of keys deleted (0 or 1).
// Typically called after successful authentication to reset the counter.
func (x *Locker) Delete(ctx context.Context, name string) int64 {
	return x.RDb.Del(ctx, x.Key(name)).Val()
}
