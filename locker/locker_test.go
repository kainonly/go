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

	// 测试自定义前缀
	x2 := locker.New(x.RDb, locker.SetPrefix("rate"))
	assert.Equal(t, "rate:api:users", x2.Key("api:users"))
}

func TestIncrement(t *testing.T) {
	ctx := context.TODO()

	// 第一次自增应返回 1
	n, err := x.Increment(ctx, "test1", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), n)

	// 第二次自增应返回 2
	n, err = x.Increment(ctx, "test1", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), n)

	// TTL 应已被设置
	ttl := x.RDb.TTL(ctx, x.Key("test1")).Val()
	assert.True(t, ttl > 0)

	// 清理
	x.Delete(ctx, "test1")
}

func TestIncrement_PreservesTTL(t *testing.T) {
	ctx := context.TODO()

	// 第一次以 1 秒 TTL 自增
	_, err := x.Increment(ctx, "test2", time.Second)
	assert.NoError(t, err)

	// 稍作等待
	time.Sleep(100 * time.Millisecond)

	// 第二次自增不应重置 TTL
	_, err = x.Increment(ctx, "test2", time.Minute) // 不同的 TTL，应被忽略
	assert.NoError(t, err)

	// TTL 仍应接近原始值（小于 1 秒）
	ttl := x.RDb.TTL(ctx, x.Key("test2")).Val()
	assert.True(t, ttl > 0)
	assert.True(t, ttl < 2*time.Second)

	// 清理
	x.Delete(ctx, "test2")
}

func TestCheck_NotLocked(t *testing.T) {
	ctx := context.TODO()

	// 不存在的键不应被锁定
	err := x.Check(ctx, "nonexistent", 5)
	assert.NoError(t, err)

	// 创建一个低于上限的计数
	x.Increment(ctx, "test3", time.Minute)
	x.Increment(ctx, "test3", time.Minute)

	// 不应被锁定（count=2, max=5）
	err = x.Check(ctx, "test3", 5)
	assert.NoError(t, err)

	// 清理
	x.Delete(ctx, "test3")
}

func TestCheck_Locked(t *testing.T) {
	ctx := context.TODO()

	// 自增达到上限
	for i := 0; i < 5; i++ {
		x.Increment(ctx, "test4", time.Minute)
	}

	// 应被锁定（count=5, max=5）
	err := x.Check(ctx, "test4", 5)
	assert.ErrorIs(t, err, locker.ErrLocked)

	// 清理
	x.Delete(ctx, "test4")
}

func TestGet(t *testing.T) {
	ctx := context.TODO()

	// 不存在的键应返回 0
	n, err := x.Get(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), n)

	// 自增之后
	x.Increment(ctx, "test5", time.Minute)
	x.Increment(ctx, "test5", time.Minute)
	x.Increment(ctx, "test5", time.Minute)

	n, err = x.Get(ctx, "test5")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), n)

	// 清理
	x.Delete(ctx, "test5")
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()

	x.Increment(ctx, "test6", time.Minute)

	// 删除存在的键
	result := x.Delete(ctx, "test6")
	assert.Equal(t, int64(1), result)

	// 删除不存在的键
	result = x.Delete(ctx, "test6")
	assert.Equal(t, int64(0), result)
}

func TestExpiration(t *testing.T) {
	ctx := context.TODO()

	// 以较短的 TTL 创建
	x.Increment(ctx, "test7", 50*time.Millisecond)

	// 应存在
	n, _ := x.Get(ctx, "test7")
	assert.Equal(t, int64(1), n)

	// 等待过期
	time.Sleep(100 * time.Millisecond)

	// 应已消失
	n, _ = x.Get(ctx, "test7")
	assert.Equal(t, int64(0), n)
}
