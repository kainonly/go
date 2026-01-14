// Package csrf provides CSRF (Cross-Site Request Forgery) protection middleware for Hertz.
//
// It implements the Double Submit Cookie pattern using HMAC-SHA256.
// Cookies are session-level and automatically cleared when the browser closes.
//
// # Hertz Backend Setup
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
// # Angular Frontend Setup
//
// Angular has built-in XSRF support that works with default cookie/header names.
//
// 1. Configure HttpClient in app.config.ts:
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
// 2. Create CSRF service:
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
// 3. Get CSRF token on app initialization (app.component.ts):
//
//	export class AppComponent implements OnInit {
//	  constructor(private csrfService: CsrfService) {}
//
//	  ngOnInit() {
//	    this.csrfService.initToken().subscribe();
//	  }
//	}
//
// 4. Refresh CSRF token after login (recommended):
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
// 5. Ensure all HTTP requests include credentials:
//
//	this.http.post('/api/submit', data, { withCredentials: true })
//
// # Security Notes
//
//   - Cookies are session-level (cleared when browser closes)
//   - XSRF-TOKEN cookie is readable by JavaScript (HttpOnly=false)
//   - XSRF-SALT cookie is HttpOnly=true for additional security
//   - Both cookies use SameSite=Strict to prevent cross-site requests
//   - Refresh CSRF token after login to prevent pre-auth token theft
//   - Always use HTTPS in production
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

// Default configuration values.
const (
	DefaultCookieName = "XSRF-TOKEN"
	DefaultSaltName   = "XSRF-SALT"
	DefaultHeaderName = "X-XSRF-TOKEN"
	DefaultSaltLength = 16
)

// Errors returned by csrf functions.
var (
	ErrMissingHeader = errors.New("csrf: missing token in header")
	ErrMissingSalt   = errors.New("csrf: missing salt cookie")
	ErrInvalidToken  = errors.New("csrf: invalid token")
	ErrEmptyKey      = errors.New("csrf: secret key cannot be empty")
)

// Csrf provides CSRF protection using Double Submit Cookie pattern.
type Csrf struct {
	Key           string
	CookieName    string
	SaltName      string
	HeaderName    string
	Domain        string
	IgnoreMethods map[string]bool
}

// New creates a new Csrf instance with the given options.
// At minimum, SetKey must be provided with a secret key.
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

// Option is a function that configures a Csrf instance.
type Option func(x *Csrf)

// SetKey sets the secret key for HMAC signing.
// The key should be at least 32 bytes for security.
func SetKey(v string) Option {
	return func(x *Csrf) {
		x.Key = v
	}
}

// SetCookieName sets the name of the token cookie.
func SetCookieName(v string) Option {
	return func(x *Csrf) {
		x.CookieName = v
	}
}

// SetSaltName sets the name of the salt cookie.
func SetSaltName(v string) Option {
	return func(x *Csrf) {
		x.SaltName = v
	}
}

// SetHeaderName sets the expected header name for token verification.
func SetHeaderName(v string) Option {
	return func(x *Csrf) {
		x.HeaderName = v
	}
}

// SetIgnoreMethods sets which HTTP methods should skip CSRF verification.
func SetIgnoreMethods(methods []string) Option {
	return func(x *Csrf) {
		x.IgnoreMethods = map[string]bool{}
		for _, v := range methods {
			x.IgnoreMethods[v] = true
		}
	}
}

// SetDomain sets the cookie domain.
func SetDomain(v string) Option {
	return func(x *Csrf) {
		x.Domain = v
	}
}

// SetToken generates and sets CSRF cookies on the response.
// Cookies are session-level (deleted when browser closes).
// Call this on login or when the frontend needs a fresh token.
func (x *Csrf) SetToken(c *app.RequestContext) {
	salt := help.Random(DefaultSaltLength)
	c.SetCookie(x.SaltName, salt, 0, "/", x.Domain, protocol.CookieSameSiteStrictMode, true, true)
	c.SetCookie(x.CookieName, x.Tokenize(salt), 0, "/", x.Domain, protocol.CookieSameSiteStrictMode, true, false)
}

// Tokenize creates an HMAC-SHA256 token from the given salt.
func (x *Csrf) Tokenize(salt string) string {
	h := hmac.New(sha256.New, []byte(x.Key))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyToken returns a Hertz middleware that validates CSRF tokens.
// Safe methods (GET, HEAD, OPTIONS, TRACE) are skipped by default.
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
