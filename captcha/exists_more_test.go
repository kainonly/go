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

func TestExists_AfterDelete(t *testing.T) {
	url := os.Getenv("DATABASE_REDIS")
	if url == "" {
		t.Skip("DATABASE_REDIS not set")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Skip("redis url invalid")
	}
	x := captcha.New(redis.NewClient(opts))
	ctx := context.TODO()
	status := x.Create(ctx, "dev3", "abcd", time.Second*1)
	assert.Equal(t, "OK", status)
	res := x.Exists(ctx, "dev3")
	assert.True(t, res)
	x.Delete(ctx, "dev3")
	res = x.Exists(ctx, "dev3")
	assert.False(t, res)
}
