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

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	status := x.Create(ctx, "dev1", "abcd", time.Second*60)
	assert.Equal(t, "OK", status)
	status = x.Create(ctx, "dev2", "abcd", time.Millisecond)
	assert.Equal(t, "OK", status)
}

func TestVerify(t *testing.T) {
	var err error
	err = x.Verify(context.TODO(), "dev1", "abc")
	assert.ErrorIs(t, err, captcha.ErrCaptchaInconsistent)
	err = x.Verify(context.TODO(), "dev1", "abcd")
	assert.NoError(t, err)
	time.Sleep(time.Second)
	err = x.Verify(context.TODO(), "dev2", "abcd")
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	exists := x.Exists(context.TODO(), "dev1")
	assert.True(t, exists)
	result := x.Delete(context.TODO(), "dev1")
	assert.Equal(t, int64(1), result)
	exists = x.Exists(context.TODO(), "dev1")
	assert.False(t, exists)
}
