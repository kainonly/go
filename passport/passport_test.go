package passport_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kainonly/go/passport"
	"github.com/stretchr/testify/assert"
)

var x1 *passport.Passport
var x2 *passport.Passport

var key1 = "hZXD^@K9%wydDC3Z@cyDvE%5bz9SP7gy"

func TestMain(m *testing.M) {
	x1 = passport.New(
		passport.SetIssuer("dev"),
		passport.SetKey(key1),
	)
	x2 = passport.New(
		passport.SetIssuer("beta"),
		passport.SetKey("eK4qpn7yCBLo0u5mlAFFRCRsCmf2NQ76"),
	)
	os.Exit(m.Run())
}

var jti1 = "GIlmuxUX1n5N4wAVVF40i"
var userId1 = "FTFD1FnWKwueHAY8h-zXg"
var jti2 = "gxOWtI58ViI2pl3BHxSNs"
var userId2 = "HU3kev7LZEgoaghpIMrGn"
var token string
var otherToken string

func TestCreate(t *testing.T) {
	var err error
	token, err = x1.Create(passport.NewClaims(userId1, time.Hour*2).SetJTI(jti1))
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	otherToken, err = x2.Create(passport.NewClaims(userId2, time.Hour*2).SetJTI(jti2))
	assert.NoError(t, err)
	assert.NotEmpty(t, otherToken)
}

func TestVerify(t *testing.T) {
	var err error
	var claims1 passport.Claims
	claims1, err = x1.Verify(token)
	assert.NoError(t, err)
	assert.Equal(t, jti1, claims1.ID)
	assert.Equal(t, userId1, claims1.ActiveId)
	assert.Equal(t, x1.Issuer, claims1.Issuer)

	var claims2 passport.Claims
	claims2, err = x2.Verify(otherToken)
	assert.NoError(t, err)
	assert.Equal(t, jti2, claims2.ID)
	assert.Equal(t, userId2, claims2.ActiveId)
	assert.Equal(t, x2.Issuer, claims2.Issuer)

	// Cross-verification should fail (different keys)
	_, err = x1.Verify(otherToken)
	assert.Error(t, err)
	_, err = x2.Verify(token)
	assert.Error(t, err)
}

func TestVerify_InvalidIssuer(t *testing.T) {
	// Create a token with x1's key but different issuer
	x3 := passport.New(
		passport.SetIssuer("other"),
		passport.SetKey(key1),
	)
	tokenOther, err := x3.Create(passport.NewClaims(userId1, time.Hour*2))
	assert.NoError(t, err)

	// Verify with x1 should fail due to issuer mismatch
	_, err = x1.Verify(tokenOther)
	assert.ErrorIs(t, err, passport.ErrInvalidIssuer)
}

func TestVerify_InvalidSigningMethod_HS384(t *testing.T) {
	// HS384 should be rejected (only HS256 is allowed)
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, passport.Claims{
		ActiveId: userId1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dev",
			ID:        jti1,
		},
	})
	ts, err := token.SignedString([]byte(key1))
	assert.NoError(t, err)
	_, err = x1.Verify(ts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "passport: invalid signing method")
}

func TestVerify_InvalidSigningMethod_ES256(t *testing.T) {
	// ECDSA should be rejected
	ecPKey := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIAh5qA3rmqQQuu0vbKV/+zouz/y/Iy2pLpIcWUSyImSwoAoGCCqGSM49
AwEHoUQDQgAEYD54V/vp+54P9DXarYqx4MPcm+HKRIQzNasYSoRQHQ/6S6Ps8tpM
cT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
-----END EC PRIVATE KEY-----`

	token := jwt.NewWithClaims(jwt.SigningMethodES256, passport.Claims{
		ActiveId: userId1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dev",
			ID:        jti1,
		},
	})
	ecdsaKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(ecPKey))
	assert.NoError(t, err)
	ts, err := token.SignedString(ecdsaKey)
	assert.NoError(t, err)
	_, err = x1.Verify(ts)
	assert.Error(t, err)
}

func TestSetData(t *testing.T) {
	claims := passport.NewClaims(userId1, time.Hour*2).SetJTI(jti1).SetData(map[string]interface{}{
		"role":   "admin",
		"active": true,
	})
	ts, err := x1.Create(claims)
	assert.NoError(t, err)
	parsed, err := x1.Verify(ts)
	assert.NoError(t, err)
	assert.Equal(t, "admin", parsed.Data["role"])
	assert.Equal(t, true, parsed.Data["active"])
}

func TestNewClaims(t *testing.T) {
	claims := passport.NewClaims("user123", time.Hour)

	assert.Equal(t, "user123", claims.ActiveId)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
	// ExpiresAt should be approximately 1 hour from now
	assert.True(t, claims.ExpiresAt.Time.After(time.Now().Add(59*time.Minute)))
	assert.True(t, claims.ExpiresAt.Time.Before(time.Now().Add(61*time.Minute)))
}
