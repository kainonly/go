package captcha_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kainonly/go/captcha"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var x *captcha.Captcha

func TestMain(m *testing.M) {
	url := os.Getenv("DATABASE_REDIS")
	if url == "" {
		os.Exit(0)
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		os.Exit(0)
	}
	x = captcha.New(redis.NewClient(opts))
	os.Exit(m.Run())
}

func TestKey(t *testing.T) {
	assert.Equal(t, "captcha:login:test", x.Key("login:test"))

	// Test custom prefix
	x2 := captcha.New(x.RDb, captcha.SetPrefix("code"))
	assert.Equal(t, "code:sms:123", x2.Key("sms:123"))
}

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	status := x.Create(ctx, "test1", "123456", time.Minute)
	assert.Equal(t, "OK", status)

	// Cleanup
	x.Delete(ctx, "test1")
}

func TestExists(t *testing.T) {
	ctx := context.TODO()

	// Create and check exists
	x.Create(ctx, "test2", "abcd", time.Minute)
	assert.True(t, x.Exists(ctx, "test2"))

	// Delete and check not exists
	x.Delete(ctx, "test2")
	assert.False(t, x.Exists(ctx, "test2"))
}

func TestVerify_Success(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test3", "correct", time.Minute)

	// Verify with correct code
	err := x.Verify(ctx, "test3", "correct")
	assert.NoError(t, err)

	// Code should be deleted after successful verification (one-time use)
	assert.False(t, x.Exists(ctx, "test3"))
}

func TestVerify_InvalidCode(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test4", "secret", time.Minute)

	// Verify with wrong code
	err := x.Verify(ctx, "test4", "wrong")
	assert.ErrorIs(t, err, captcha.ErrInvalidCode)

	// Code should be deleted even on failed verification (GetDel behavior)
	assert.False(t, x.Exists(ctx, "test4"))
}

func TestVerify_NotExists(t *testing.T) {
	ctx := context.TODO()

	// Verify non-existent code
	err := x.Verify(ctx, "nonexistent", "any")
	assert.ErrorIs(t, err, captcha.ErrNotExists)
}

func TestVerify_Expired(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test5", "temp", 50*time.Millisecond)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should return not exists error
	err := x.Verify(ctx, "test5", "temp")
	assert.ErrorIs(t, err, captcha.ErrNotExists)
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test6", "todelete", time.Minute)

	// Delete existing key
	result := x.Delete(ctx, "test6")
	assert.Equal(t, int64(1), result)

	// Delete non-existing key
	result = x.Delete(ctx, "test6")
	assert.Equal(t, int64(0), result)
}
