package csrf_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/kainonly/go/csrf"
	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	tok := x.Tokenize("abcd1234")
	assert.Len(t, tok, 64)
}

func TestOptions(t *testing.T) {
	x := csrf.New(
		csrf.SetKey("secret"),
		csrf.SetCookieName("A"),
		csrf.SetSaltName("B"),
		csrf.SetHeaderName("C"),
		csrf.SetIgnoreMethods([]string{"POST"}),
		csrf.SetDomain("example.com"),
	)
	assert.Equal(t, "A", x.CookieName)
	assert.Equal(t, "B", x.SaltName)
	assert.Equal(t, "C", x.HeaderName)
	assert.True(t, x.IgnoreMethods["POST"])
	assert.Equal(t, "example.com", x.Domain)
}

func TestVerifyToken_Skip(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api", x.VerifyToken(true), func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"ok": 1})
	})
	w := ut.PerformRequest(router, "POST", "/api", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestVerifyToken_IgnoreMethods(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/api", x.VerifyToken(false), func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"ok": 1})
	})
	w := ut.PerformRequest(router, "GET", "/api", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestVerifyToken_MissingHeader(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api", x.VerifyToken(false), func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"ok": 1})
	})
	salt := "abcd1234"
	w := ut.PerformRequest(router, "POST", "/api", &ut.Body{bytes.NewBuffer(nil), 0}, ut.Header{"Cookie", "XSRF-SALT=" + salt})
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api", x.VerifyToken(false), func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"ok": 1})
	})
	salt := "abcd1234"
	w := ut.PerformRequest(router, "POST", "/api", &ut.Body{bytes.NewBuffer(nil), 0}, ut.Header{"Cookie", "XSRF-SALT=" + salt}, ut.Header{"X-XSRF-TOKEN", "invalid"})
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
}

func TestVerifyToken_Success(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api", x.VerifyToken(false), func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"ok": 1})
	})
	salt := "abcd1234"
	token := x.Tokenize(salt)
	w := ut.PerformRequest(router, "POST", "/api", &ut.Body{bytes.NewBuffer(nil), 0}, ut.Header{"Cookie", "XSRF-SALT=" + salt}, ut.Header{"X-XSRF-TOKEN", token})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestVerifyToken_CustomIgnore(t *testing.T) {
	x := csrf.New(csrf.SetKey("secret"), csrf.SetIgnoreMethods([]string{"POST"}))
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.POST("/api", x.VerifyToken(false), func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, utils.H{"ok": 1})
	})
	w := ut.PerformRequest(router, "POST", "/api", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}
