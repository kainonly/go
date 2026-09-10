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

	// 测试自定义前缀
	x2 := captcha.New(x.RDb, captcha.SetPrefix("code"))
	assert.Equal(t, "code:sms:123", x2.Key("sms:123"))
}

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	status := x.Create(ctx, "test1", "123456", time.Minute)
	assert.Equal(t, "OK", status)

	// 清理
	x.Delete(ctx, "test1")
}

func TestExists(t *testing.T) {
	ctx := context.TODO()

	// 创建后检查存在性
	x.Create(ctx, "test2", "abcd", time.Minute)
	assert.True(t, x.Exists(ctx, "test2"))

	// 删除后检查不存在
	x.Delete(ctx, "test2")
	assert.False(t, x.Exists(ctx, "test2"))
}

func TestVerify_Success(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test3", "correct", time.Minute)

	// 使用正确的验证码校验
	err := x.Verify(ctx, "test3", "correct")
	assert.NoError(t, err)

	// 验证成功后验证码应被删除（一次性使用）
	assert.False(t, x.Exists(ctx, "test3"))
}

func TestVerify_InvalidCode(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test4", "secret", time.Minute)

	// 使用错误的验证码校验
	err := x.Verify(ctx, "test4", "wrong")
	assert.ErrorIs(t, err, captcha.ErrInvalidCode)

	// 验证失败时验证码也应被删除（GetDel 行为）
	assert.False(t, x.Exists(ctx, "test4"))
}

func TestVerify_NotExists(t *testing.T) {
	ctx := context.TODO()

	// 校验不存在的验证码
	err := x.Verify(ctx, "nonexistent", "any")
	assert.ErrorIs(t, err, captcha.ErrNotExists)
}

func TestVerify_Expired(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test5", "temp", 50*time.Millisecond)

	// 等待过期
	time.Sleep(100 * time.Millisecond)

	// 应返回不存在的错误
	err := x.Verify(ctx, "test5", "temp")
	assert.ErrorIs(t, err, captcha.ErrNotExists)
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()
	x.Create(ctx, "test6", "todelete", time.Minute)

	// 删除存在的键
	result := x.Delete(ctx, "test6")
	assert.Equal(t, int64(1), result)

	// 删除不存在的键
	result = x.Delete(ctx, "test6")
	assert.Equal(t, int64(0), result)
}
