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

func TestLockerUpdate(t *testing.T) {
	ctx := context.TODO()
	n := x.Update(ctx, "dev", time.Second*60)
	assert.Equal(t, int64(1), n)
	ttl := x.RDb.TTL(ctx, x.Key("dev")).Val()
	t.Log(ttl.Seconds())
}

func TestLockerVerify(t *testing.T) {
	ctx := context.TODO()
	err := x.Verify(ctx, "dev", 3)
	assert.NoError(t, err)

	for i := 0; i < 3; i++ {
		n := x.Update(ctx, "dev", time.Second*60)
		assert.Equal(t, int64(i+2), n)
	}

	err = x.Verify(ctx, "dev", 3)
	assert.ErrorIs(t, err, locker.ErrLocked)
}

func TestLockerVerifyNotExists(t *testing.T) {
	ctx := context.TODO()
	err := x.Verify(ctx, "unknow", 3)
	assert.ErrorIs(t, err, locker.ErrLockerNotExists)
}

func TestLockerVerifyBad(t *testing.T) {
	ctx := context.TODO()
	status := x.RDb.Set(ctx, x.Key("notnumber"), "abc", time.Second*10).Val()
	assert.Equal(t, "OK", status)
	err := x.Verify(ctx, "notnumber", 3)
	assert.Error(t, err)
	t.Log(err)
}

func TestLockerDelete(t *testing.T) {
	ctx := context.TODO()
	result := x.Delete(ctx, "dev")
	assert.Equal(t, int64(1), result)
	count := x.RDb.Exists(ctx, x.Key("dev")).Val()
	assert.Equal(t, int64(0), count)
}
