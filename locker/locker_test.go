package locker_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kainonly/go/locker"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var x *locker.Locker

func TestMain(m *testing.M) {
	url := os.Getenv("DATABASE_REDIS")
	if url == "" {
		os.Exit(0)
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		os.Exit(0)
	}
	x = locker.New(redis.NewClient(opts))
	os.Exit(m.Run())
}

func TestKey(t *testing.T) {
	assert.Equal(t, "locker:login:test", x.Key("login:test"))

	// Test custom prefix
	x2 := locker.New(x.RDb, locker.SetPrefix("rate"))
	assert.Equal(t, "rate:api:users", x2.Key("api:users"))
}

func TestIncrement(t *testing.T) {
	ctx := context.TODO()

	// First increment should return 1
	n, err := x.Increment(ctx, "test1", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), n)

	// Second increment should return 2
	n, err = x.Increment(ctx, "test1", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)

	// TTL should be set
	ttl := x.RDb.TTL(ctx, x.Key("test1")).Val()
	assert.True(t, ttl > 0)

	// Cleanup
	x.Delete(ctx, "test1")
}

func TestIncrement_PreservesTTL(t *testing.T) {
	ctx := context.TODO()

	// First increment with 1 second TTL
	_, err := x.Increment(ctx, "test2", time.Second)
	assert.NoError(t, err)

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Second increment should not reset TTL
	_, err = x.Increment(ctx, "test2", time.Minute) // Different TTL, should be ignored
	assert.NoError(t, err)

	// TTL should still be close to original (less than 1 second)
	ttl := x.RDb.TTL(ctx, x.Key("test2")).Val()
	assert.True(t, ttl < time.Second)

	// Cleanup
	x.Delete(ctx, "test2")
}

func TestCheck_NotLocked(t *testing.T) {
	ctx := context.TODO()

	// Non-existent key should not be locked
	err := x.Check(ctx, "nonexistent", 5)
	assert.NoError(t, err)

	// Create a counter below max
	x.Increment(ctx, "test3", time.Minute)
	x.Increment(ctx, "test3", time.Minute)

	// Should not be locked (count=2, max=5)
	err = x.Check(ctx, "test3", 5)
	assert.NoError(t, err)

	// Cleanup
	x.Delete(ctx, "test3")
}

func TestCheck_Locked(t *testing.T) {
	ctx := context.TODO()

	// Increment to reach max
	for i := 0; i < 5; i++ {
		x.Increment(ctx, "test4", time.Minute)
	}

	// Should be locked (count=5, max=5)
	err := x.Check(ctx, "test4", 5)
	assert.ErrorIs(t, err, locker.ErrLocked)

	// Cleanup
	x.Delete(ctx, "test4")
}

func TestGet(t *testing.T) {
	ctx := context.TODO()

	// Non-existent key should return 0
	n, err := x.Get(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)

	// After increments
	x.Increment(ctx, "test5", time.Minute)
	x.Increment(ctx, "test5", time.Minute)
	x.Increment(ctx, "test5", time.Minute)

	n, err = x.Get(ctx, "test5")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), n)

	// Cleanup
	x.Delete(ctx, "test5")
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()

	x.Increment(ctx, "test6", time.Minute)

	// Delete existing key
	result := x.Delete(ctx, "test6")
	assert.Equal(t, int64(1), result)

	// Delete non-existing key
	result = x.Delete(ctx, "test6")
	assert.Equal(t, int64(0), result)
}

func TestExpiration(t *testing.T) {
	ctx := context.TODO()

	// Create with short TTL
	x.Increment(ctx, "test7", 50*time.Millisecond)

	// Should exist
	n, _ := x.Get(ctx, "test7")
	assert.Equal(t, int64(1), n)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be gone
	n, _ = x.Get(ctx, "test7")
	assert.Equal(t, int64(0), n)
}
