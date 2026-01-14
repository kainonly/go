// Package passport provides JWT (JSON Web Token) authentication utilities for Hertz.
//
// It uses HS256 (HMAC-SHA256) for token signing and supports custom claims.
//
// # Hertz Backend Setup
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
// # Angular Frontend Setup
//
// 1. Store token after login:
//
//	login(credentials: LoginRequest) {
//	  return this.http.post<{accessToken: string}>('/auth/login', credentials).pipe(
//	    tap(res => localStorage.setItem('token', res.accessToken))
//	  );
//	}
//
// 2. Add token to requests via HTTP interceptor:
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
// 3. Register interceptor in app.config.ts:
//
//	provideHttpClient(withInterceptors([authInterceptor]))
//
// # Security Notes
//
//   - Use a strong secret key (at least 32 bytes)
//   - Set appropriate token expiration time
//   - Store tokens securely on client side
//   - Use HTTPS in production
//   - Consider implementing token refresh for long sessions
package passport

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Errors returned by passport functions.
var (
	ErrInvalidSigningMethod = errors.New("passport: invalid signing method, expected HS256")
	ErrInvalidIssuer        = errors.New("passport: token issuer does not match")
)

// Passport provides JWT token creation and verification.
type Passport struct {
	Issuer string
	Key    string
}

// New creates a new Passport instance with the given options.
// Both SetKey and SetIssuer should be provided for proper operation.
func New(options ...Option) *Passport {
	x := new(Passport)
	for _, v := range options {
		v(x)
	}
	return x
}

// Option is a function that configures a Passport instance.
type Option func(x *Passport)

// SetIssuer sets the token issuer (iss claim).
// This should be your application name or identifier.
func SetIssuer(v string) Option {
	return func(x *Passport) {
		x.Issuer = v
	}
}

// SetKey sets the secret key for signing tokens.
// The key should be at least 32 bytes for security.
func SetKey(v string) Option {
	return func(x *Passport) {
		x.Key = v
	}
}

// Claims represents the JWT claims with custom fields.
type Claims struct {
	// ActiveId is the primary identifier (usually user ID or session ID).
	ActiveId string `json:"active_id,omitempty"`
	// Data holds additional custom data.
	Data map[string]interface{} `json:"data,omitempty"`

	jwt.RegisteredClaims
}

// NewClaims creates a new Claims with the given activeId and expiration duration.
// It sets IssuedAt and NotBefore to current time.
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

// SetJTI sets the JWT ID (jti claim) for token uniqueness.
// Useful for token revocation and tracking.
func (x *Claims) SetJTI(v string) *Claims {
	x.RegisteredClaims.ID = v
	return x
}

// SetData sets custom data to be included in the token.
// Note: Keep data minimal as it increases token size.
func (x *Claims) SetData(v map[string]interface{}) *Claims {
	x.Data = v
	return x
}

// Create generates a signed JWT token string from the given claims.
// The token is signed using HS256 algorithm.
func (x *Passport) Create(claims *Claims) (string, error) {
	claims.RegisteredClaims.Issuer = x.Issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(x.Key))
}

// Verify parses and validates a JWT token string.
// It checks the signing method (HS256 only) and verifies the issuer matches.
// Returns the claims if valid, or an error if invalid.
func (x *Passport) Verify(tokenString string) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Strict check for HS256 only
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(x.Key), nil
	})
	if err != nil {
		return claims, err
	}
	// Verify issuer matches
	if claims.Issuer != x.Issuer {
		return claims, ErrInvalidIssuer
	}
	return claims, nil
}
