package help_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/go-playground/validator/v10"
	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_Public(t *testing.T) {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/public", help.ErrorHandler(), func(ctx context.Context, c *app.RequestContext) {
		c.Error(help.E(1001, "bad request"))
	})
	w := ut.PerformRequest(router, "GET", "/public", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, 400, resp.StatusCode())
}

func TestErrorHandler_Validation(t *testing.T) {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/validation", help.ErrorHandler(), func(ctx context.Context, c *app.RequestContext) {
		vd := validator.New()
		var s struct {
			Email string `validate:"email"`
		}
		s.Email = "not-email"
		err := vd.Struct(s)
		c.Error(err)
	})
	w := ut.PerformRequest(router, "GET", "/validation", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, 400, resp.StatusCode())
}

func TestErrorHandler_Internal(t *testing.T) {
	os.Setenv("MODE", "dev")
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/internal", help.ErrorHandler(), func(ctx context.Context, c *app.RequestContext) {
		c.Error(errors.New("something wrong"))
	})
	w := ut.PerformRequest(router, "GET", "/internal", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, 500, resp.StatusCode())
}
