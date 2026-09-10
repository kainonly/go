// Package csrf 为 Hertz 提供 CSRF（跨站请求伪造）防护中间件。
//
// 基于 HMAC-SHA256 实现 Double Submit Cookie 模式。
// Cookie 为会话级，浏览器关闭时自动清除。
//
// # Hertz 后端配置
//
//	// Initialize CSRF protection
//	csrfProtect := csrf.New(
//		csrf.SetKey("your-secret-key-at-least-32-bytes"),
//		csrf.SetDomain("example.com"),
//	)
//
//	// CSRF token endpoint (called on first page load)
//	h.GET("/csrf-token", func(ctx context.Context, c *app.RequestContext) {
//		csrfProtect.SetToken(c)
//		c.JSON(200, utils.H{"message": "ok"})
//	})
//
//	// Login endpoint (refresh CSRF token after login for security)
//	h.POST("/auth/login", func(ctx context.Context, c *app.RequestContext) {
//		// ... validate credentials ...
//		csrfProtect.SetToken(c) // Refresh CSRF token
//		c.JSON(200, utils.H{"accessToken": token})
//	})
//
//	// Protect API routes with CSRF middleware
//	api := h.Group("/api", csrfProtect.VerifyToken())
//	api.POST("/submit", submitHandler)
//	api.PUT("/update", updateHandler)
//	api.DELETE("/remove", removeHandler)
//
// # Angular 前端配置
//
// Angular 内置 XSRF 支持，可直接使用默认的 cookie/header 名称。
//
// 1. 在 app.config.ts 中配置 HttpClient：
//
//	import { provideHttpClient, withXsrfConfiguration } from '@angular/common/http';
//
//	export const appConfig: ApplicationConfig = {
//	  providers: [
//	    provideHttpClient(
//	      withXsrfConfiguration({
//	        cookieName: 'XSRF-TOKEN',  // Must match backend default
//	        headerName: 'X-XSRF-TOKEN' // Must match backend default
//	      })
//	    )
//	  ]
//	};
//
// 2. 创建 CSRF 服务：
//
//	@Injectable({ providedIn: 'root' })
//	export class CsrfService {
//	  constructor(private http: HttpClient) {}
//
//	  initToken(): Observable<void> {
//	    return this.http.get<void>('/csrf-token', { withCredentials: true });
//	  }
//	}
//
// 3. 在应用初始化时获取 CSRF token（app.component.ts）：
//
//	export class AppComponent implements OnInit {
//	  constructor(private csrfService: CsrfService) {}
//
//	  ngOnInit() {
//	    this.csrfService.initToken().subscribe();
//	  }
//	}
//
// 4. 登录后刷新 CSRF token（推荐）：
//
//	login(credentials: LoginRequest) {
//	  return this.http.post<LoginResponse>('/auth/login', credentials, {
//	    withCredentials: true
//	  }).pipe(
//	    tap(res => localStorage.setItem('token', res.accessToken))
//	    // CSRF cookie is refreshed automatically by backend
//	  );
//	}
//
// 5. 确保所有 HTTP 请求都携带凭证：
//
//	this.http.post('/api/submit', data, { withCredentials: true })
//
// # 安全注意事项
//
//   - Cookie 为会话级（浏览器关闭时清除）
//   - XSRF-TOKEN cookie 可被 JavaScript 读取（HttpOnly=false）
//   - XSRF-SALT cookie 设置 HttpOnly=true，提供额外安全防护
//   - 两个 cookie 均使用 SameSite=Strict 防止跨站请求
//   - 登录后刷新 CSRF token，防止认证前的 token 被盗用
//   - 生产环境务必使用 HTTPS
package csrf

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/kainonly/go/help"
)

// 默认配置值。
const (
	DefaultCookieName = "XSRF-TOKEN"
	DefaultSaltName   = "XSRF-SALT"
	DefaultHeaderName = "X-XSRF-TOKEN"
	DefaultSaltLength = 16
)

// csrf 函数返回的错误。
var (
	ErrMissingHeader = errors.New("csrf: missing token in header")
	ErrMissingSalt   = errors.New("csrf: missing salt cookie")
	ErrInvalidToken  = errors.New("csrf: invalid token")
	ErrEmptyKey      = errors.New("csrf: secret key cannot be empty")
)

// Csrf 基于 Double Submit Cookie 模式提供 CSRF 防护。
type Csrf struct {
	Key           string
	CookieName    string
	SaltName      string
	HeaderName    string
	Domain        string
	IgnoreMethods map[string]bool
}

// New 使用给定的选项创建一个新的 Csrf 实例。
// 至少需要通过 SetKey 提供密钥。
func New(options ...Option) *Csrf {
	x := &Csrf{
		CookieName: DefaultCookieName,
		SaltName:   DefaultSaltName,
		HeaderName: DefaultHeaderName,
		Domain:     "",
		IgnoreMethods: map[string]bool{
			"GET":     true,
			"HEAD":    true,
			"OPTIONS": true,
			"TRACE":   true,
		},
	}
	for _, v := range options {
		v(x)
	}
	return x
}

// Option 是用于配置 Csrf 实例的函数。
type Option func(x *Csrf)

// SetKey 设置用于 HMAC 签名的密钥。
// 出于安全考虑，密钥至少应为 32 字节。
func SetKey(v string) Option {
	return func(x *Csrf) {
		x.Key = v
	}
}

// SetCookieName 设置 token cookie 的名称。
func SetCookieName(v string) Option {
	return func(x *Csrf) {
		x.CookieName = v
	}
}

// SetSaltName 设置 salt cookie 的名称。
func SetSaltName(v string) Option {
	return func(x *Csrf) {
		x.SaltName = v
	}
}

// SetHeaderName 设置验证 token 时使用的请求头名称。
func SetHeaderName(v string) Option {
	return func(x *Csrf) {
		x.HeaderName = v
	}
}

// SetIgnoreMethods 设置跳过 CSRF 验证的 HTTP 方法。
func SetIgnoreMethods(methods []string) Option {
	return func(x *Csrf) {
		x.IgnoreMethods = map[string]bool{}
		for _, v := range methods {
			x.IgnoreMethods[v] = true
		}
	}
}

// SetDomain 设置 cookie 的域。
func SetDomain(v string) Option {
	return func(x *Csrf) {
		x.Domain = v
	}
}

// SetToken 生成 CSRF cookie 并写入响应。
// Cookie 为会话级（浏览器关闭时删除）。
// 在登录时或前端需要新 token 时调用。
func (x *Csrf) SetToken(c *app.RequestContext) {
	salt := help.Random(DefaultSaltLength)
	c.SetCookie(x.SaltName, salt, 0, "/", x.Domain, protocol.CookieSameSiteStrictMode, true, true)
	c.SetCookie(x.CookieName, x.Tokenize(salt), 0, "/", x.Domain, protocol.CookieSameSiteStrictMode, true, false)
}

// Tokenize 根据给定的 salt 生成 HMAC-SHA256 token。
func (x *Csrf) Tokenize(salt string) string {
	h := hmac.New(sha256.New, []byte(x.Key))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyToken 返回用于校验 CSRF token 的 Hertz 中间件。
// 安全方法（GET、HEAD、OPTIONS、TRACE）默认跳过校验。
func (x *Csrf) VerifyToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if x.IgnoreMethods[string(c.Method())] {
			c.Next(ctx)
			return
		}

		salt := string(c.Cookie(x.SaltName))
		if salt == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.H{
				"code":    0,
				"message": ErrMissingSalt.Error(),
			})
			return
		}

		token := c.GetHeader(x.HeaderName)
		if token == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.H{
				"code":    0,
				"message": ErrMissingHeader.Error(),
			})
			return
		}

		if !hmac.Equal([]byte(x.Tokenize(salt)), token) {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.H{
				"code":    0,
				"message": ErrInvalidToken.Error(),
			})
			return
		}

		c.Next(ctx)
	}
}
