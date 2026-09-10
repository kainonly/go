// Package passport 为 Hertz 提供 JWT（JSON Web Token）认证工具。
//
// 使用 HS256（HMAC-SHA256）签名令牌，支持自定义 claims。
//
// # Hertz 后端配置
//
//	// Initialize passport
//	auth := passport.New(
//		passport.SetKey("your-secret-key-at-least-32-bytes"),
//		passport.SetIssuer("your-app-name"),
//	)
//
//	// Login endpoint - create JWT token
//	h.POST("/auth/login", func(ctx context.Context, c *app.RequestContext) {
//		// ... validate credentials ...
//		claims := passport.NewClaims(userId, 2*time.Hour).
//			SetJTI(tokenId).
//			SetData(map[string]interface{}{"role": "admin"})
//		token, err := auth.Create(claims)
//		if err != nil {
//			c.JSON(500, utils.H{"error": err.Error()})
//			return
//		}
//		c.JSON(200, utils.H{"accessToken": token})
//	})
//
//	// Auth middleware - verify JWT token
//	func AuthMiddleware(auth *passport.Passport) app.HandlerFunc {
//		return func(ctx context.Context, c *app.RequestContext) {
//			token := c.GetHeader("Authorization")
//			if token == nil {
//				c.AbortWithStatusJSON(401, utils.H{"error": "missing token"})
//				return
//			}
//			// Remove "Bearer " prefix
//			tokenStr := strings.TrimPrefix(string(token), "Bearer ")
//			claims, err := auth.Verify(tokenStr)
//			if err != nil {
//				c.AbortWithStatusJSON(401, utils.H{"error": err.Error()})
//				return
//			}
//			c.Set("userId", claims.ActiveId)
//			c.Set("claims", claims)
//			c.Next(ctx)
//		}
//	}
//
//	// Protected routes
//	api := h.Group("/api", AuthMiddleware(auth))
//	api.GET("/profile", profileHandler)
//
// # Angular 前端配置
//
// 1. 登录后保存令牌：
//
//	login(credentials: LoginRequest) {
//	  return this.http.post<{accessToken: string}>('/auth/login', credentials).pipe(
//	    tap(res => localStorage.setItem('token', res.accessToken))
//	  );
//	}
//
// 2. 通过 HTTP 拦截器为请求附加令牌：
//
//	export const authInterceptor: HttpInterceptorFn = (req, next) => {
//	  const token = localStorage.getItem('token');
//	  if (token) {
//	    req = req.clone({
//	      setHeaders: { Authorization: `Bearer ${token}` }
//	    });
//	  }
//	  return next(req);
//	};
//
// 3. 在 app.config.ts 中注册拦截器：
//
//	provideHttpClient(withInterceptors([authInterceptor]))
//
// # 安全注意事项
//
//   - 使用强密钥（至少 32 字节）
//   - 设置合适的令牌过期时间
//   - 在客户端安全存储令牌
//   - 生产环境使用 HTTPS
//   - 长会话场景考虑实现令牌刷新
package passport

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// passport 函数返回的错误。
var (
	ErrInvalidSigningMethod = errors.New("passport: invalid signing method, expected HS256")
	ErrInvalidIssuer        = errors.New("passport: token issuer does not match")
)

// Passport 提供 JWT 令牌的创建与校验。
type Passport struct {
	Issuer string
	Key    string
}

// New 使用给定选项创建 Passport 实例。
// 正常使用需同时提供 SetKey 和 SetIssuer。
func New(options ...Option) *Passport {
	x := new(Passport)
	for _, v := range options {
		v(x)
	}
	return x
}

// Option 是用于配置 Passport 实例的函数。
type Option func(x *Passport)

// SetIssuer 设置令牌签发者（iss claim），
// 应为应用名称或标识。
func SetIssuer(v string) Option {
	return func(x *Passport) {
		x.Issuer = v
	}
}

// SetKey 设置令牌签名密钥，
// 出于安全考虑应至少 32 字节。
func SetKey(v string) Option {
	return func(x *Passport) {
		x.Key = v
	}
}

// Claims 表示带自定义字段的 JWT claims。
type Claims struct {
	// ActiveId 是主标识（通常是用户 ID 或会话 ID）。
	ActiveId string `json:"active_id,omitempty"`
	// Data 存放附加的自定义数据。
	Data map[string]interface{} `json:"data,omitempty"`

	jwt.RegisteredClaims
}

// NewClaims 使用给定的 activeId 和过期时长创建 Claims，
// IssuedAt 和 NotBefore 设为当前时间。
func NewClaims(activeId string, expire time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		ActiveId: activeId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

// SetJTI 设置 JWT ID（jti claim）保证令牌唯一性，
// 适用于令牌吊销与追踪。
func (x *Claims) SetJTI(v string) *Claims {
	x.ID = v
	return x
}

// SetData 设置包含在令牌中的自定义数据。
// 注意：数据应尽量精简，会增加令牌体积。
func (x *Claims) SetData(v map[string]interface{}) *Claims {
	x.Data = v
	return x
}

// Create 根据给定 claims 生成已签名的 JWT 令牌字符串，
// 使用 HS256 算法签名。
func (x *Passport) Create(claims *Claims) (string, error) {
	claims.Issuer = x.Issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(x.Key))
}

// Verify 解析并校验 JWT 令牌字符串。
// 只允许 HS256 签名方式，并校验签发者匹配。
// 有效时返回 claims，否则返回错误。
func (x *Passport) Verify(tokenString string) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// 严格只允许 HS256
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(x.Key), nil
	})
	if err != nil {
		return claims, err
	}
	// 校验签发者匹配
	if claims.Issuer != x.Issuer {
		return claims, ErrInvalidIssuer
	}
	return claims, nil
}
