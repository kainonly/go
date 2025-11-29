package csrf_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/kainonly/go/csrf"
	"github.com/stretchr/testify/assert"
)

func TestSetToken(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/csrf", func(ctx context.Context, c *app.RequestContext) {
		x.SetToken(c)
		c.JSON(http.StatusOK, nil)
	})
	w := ut.PerformRequest(router, "GET", "/csrf", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
