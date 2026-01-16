# Go Utils

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/kainonly/go/testing.yml?style=flat-square)](https://github.com/kainonly/go/actions/workflows/testing.yml)
[![Coveralls github](https://img.shields.io/coveralls/github/kainonly/go.svg?style=flat-square)](https://coveralls.io/github/kainonly/go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainonly/go?style=flat-square)](https://github.com/kainonly/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kainonly/go?style=flat-square)](https://goreportcard.com/report/github.com/kainonly/go)
[![Release](https://img.shields.io/github/v/release/kainonly/go.svg?style=flat-square)](https://github.com/kainonly/go)
[![GitHub license](https://img.shields.io/github/license/kainonly/go?style=flat-square)](https://raw.githubusercontent.com/kainonly/go/main/LICENSE)

English | [简体中文](README.zh-CN.md)

A collection of commonly used Go utilities for web development (with Hertz-friendly integrations).

## Install

```bash
go get github.com/kainonly/go
```

## Modules

| Module | Description |
|------|------|
| [vd](#vd---validator) | Validator wrapper + Hertz integration |
| [passport](#passport---jwt-auth) | JWT auth (HS256) |
| [csrf](#csrf---csrf-protection) | CSRF protection middleware |
| [captcha](#captcha---captcha) | Captcha management (Redis) |
| [locker](#locker---lockout) | Counter & lockout helpers (Redis) |
| [passlib](#passlib---password-hashing) | Password hashing (Argon2id) |
| [totp](#totp---one-time-password) | TOTP |
| [cipher](#cipher---encryption) | Symmetric encryption helpers |
| [help](#help---helpers) | Misc helpers |

---

## Quick Start

### vd (validator)

```go
import "github.com/kainonly/go/vd"

v := vd.Default()

type User struct {
	Name  string `vd:"required,snake"`
	Phone string `vd:"required,phone"`
}
_ = v.Validate(&User{Name: "user_name", Phone: "13800138000"})
```

### passport (JWT HS256)

```go
import "github.com/kainonly/go/passport"

auth := passport.New(
	passport.SetKey("your-secret-key"),
	passport.SetIssuer("your-app"),
)

claims := passport.NewClaims()
claims.ID = "user-123"
claims.SetData(map[string]any{"role": "admin"})

token, _ := auth.Create(claims)
_, _ = auth.Verify(token)
```

### csrf (double-submit cookie)

```go
import "github.com/kainonly/go/csrf"

_ = csrf.New(
	csrf.SetKey("your-secret-key"),
	csrf.SetDomain("example.com"),
)
```

---

## vd - Validator

Wrapper around [go-playground/validator](https://github.com/go-playground/validator) with Hertz integration and a set of commonly used rules.

### Basic usage

```go
import "github.com/kainonly/go/vd"

v := vd.Default()

v := vd.New(
	vd.SetTag("vd"),
	vd.SetRules(
		vd.Snake(),
		vd.Sort(),
		vd.Phone(),
		vd.IDCard(),
		vd.PasswordMedium(),
	),
)

type User struct {
	Name  string `vd:"required,snake"`
	Phone string `vd:"required,phone"`
}
err := v.Validate(&User{Name: "user_name", Phone: "13800138000"})
_ = err
```

### Hertz integration

```go
import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/kainonly/go/vd"
)

func main() {
	v := vd.New(vd.SetRules(vd.All()...))

	h := server.Default(
		server.WithHostPorts(":8080"),
		server.WithCustomValidator(v.Engine()),
	)

	h.Spin()
}
```

### Available rules

#### Rule groups

| Function | Description |
|------|------|
| `All()` | All rules |
| `Common()` | Common rules (snake, sort, phone, idcard, username, slug, password_medium) |
| `Chinese()` | CN localized rules |
| `NamingConvention()` | Naming convention rules |

#### CN localized

| Rule | Tag | Description | Example |
|------|-----|------|------|
| `Phone()` | `phone` | Mobile phone number | `13800138000` |
| `IDCard()` | `idcard` | CN ID card number | `110101199003071234` |
| `BankCard()` | `bankcard` | Bank card number (Luhn) | `6222021234567890` |
| `LicensePlate()` | `license_plate` | License plate | `京A12345` |
| `USCC()` | `uscc` | Unified social credit code | `91310000MA1FL8TQ32` |
| `ChineseWord()` | `chinese` | Chinese-only string | `你好世界` |
| `ChineseName()` | `chinese_name` | Chinese name | `张三` |
| `TelPhone()` | `tel` | Landline phone | `010-12345678` |
| `QQ()` | `qq` | QQ number | `12345678` |
| `WeChat()` | `wechat` | WeChat ID | `wxid_abc123` |
| `ZipCode()` | `zipcode` | Postal code | `518000` |

#### Password strength

| Rule | Tag | Description |
|------|-----|------|
| `PasswordWeak()` | `password_weak` | ≥ 6 characters |
| `PasswordMedium()` | `password_medium` | ≥ 8 characters; contains letters and digits |
| `PasswordStrong()` | `password_strong` | ≥ 8 characters; contains upper+lower+digits+special |

#### Naming conventions

| Rule | Tag | Example |
|------|-----|------|
| `Snake()` | `snake` | `user_name` |
| `PascalCase()` | `pascal` | `UserName` |
| `CamelCase()` | `camel` | `userName` |
| `KebabCase()` | `kebab` | `user-name` |
| `UpperSnake()` | `upper_snake` | `USER_NAME` |
| `Variable()` | `variable` | `_privateVar` |

#### Misc

| Rule | Tag | Description |
|------|-----|------|
| `Sort()` | `sort` | Sort format: `field:1` or `field:-1` |
| `Username()` | `username` | Username (3-20 chars) |
| `Slug()` | `slug` | URL slug |
| `ObjectID()` | `objectid` | MongoDB ObjectId |
| `Snowflake()` | `snowflake` | Snowflake ID |
| `Version()` | `version` | Semver like `1.0.0` |
| `SafeString()` | `safe_string` | Injection-safe string |
| `AlphaNumDash()` | `alphanumdash` | Letters, digits, underscore, dash |
| `Decimal()` | `decimal` | Decimal string |
| `PositiveDecimal()` | `positive_decimal` | Positive decimal string |
| `Domain()` | `domain` | Domain |
| `FileName()` | `filename` | File name |
| `FileExt()` | `file_ext` | File extension |
| `FilePath()` | `file_path` | File path |
| `Color()` | `color` | Color without `#` |
| `NotBlank()` | `notblank` | Non-whitespace string |

#### Custom rule

```go
v := vd.New(vd.SetRules(
	vd.Rule{
		Tag: "even",
		Fn: func(fl vd.FieldLevel) bool {
			return fl.Field().Int()%2 == 0
		},
	},
))
```

---

## passport - JWT Auth

JWT authentication based on HS256.

```go
import "github.com/kainonly/go/passport"

auth := passport.New(
	passport.SetKey("your-secret-key"),
	passport.SetIssuer("your-app"),
)

claims := passport.NewClaims()
claims.ID = "user-123"
claims.SetData(map[string]any{"role": "admin"})
token, err := auth.Create(claims)

claims, err = auth.Verify(token)
_ = claims
_ = err
```

---

## csrf - CSRF Protection

CSRF middleware using the double-submit cookie approach.

```go
import "github.com/kainonly/go/csrf"

c := csrf.New(
	csrf.SetKey("your-secret-key"),
	csrf.SetDomain("example.com"),
)
_ = c
```

---

## captcha - Captcha

Captcha management backed by Redis.

```go
import "github.com/kainonly/go/captcha"

cap := captcha.New(redisClient)

cap.Create(ctx, "login:user123", "123456", 5*time.Minute)

err := cap.Verify(ctx, "login:user123", "123456")
if errors.Is(err, captcha.ErrInvalidCode) {
}
if errors.Is(err, captcha.ErrNotExists) {
}
```

---

## locker - Lockout

Redis-based counter and lockout helper for scenarios like login failures and rate limiting.

```go
import "github.com/kainonly/go/locker"

lock := locker.New(redisClient)

count, err := lock.Increment(ctx, "login:user123", time.Minute)
_ = count
_ = err

err = lock.Check(ctx, "login:user123", 5)
if errors.Is(err, locker.ErrLocked) {
}

lock.Delete(ctx, "login:user123")
```

---

## passlib - Password Hashing

Argon2id-based password hashing.

```go
import "github.com/kainonly/go/passlib"

hash, err := passlib.Hash("password123")
_ = err

err = passlib.Verify("password123", hash)

if passlib.NeedsRehash(hash) {
	newHash, _ := passlib.Hash("password123")
	_ = newHash
}
```

---

## totp - One-time Password

TOTP secret generation and code validation.

```go
import "github.com/kainonly/go/totp"

secret, err := totp.GenerateSecret()
_ = secret
_ = err

valid, err := totp.Validate(code, secret)
_ = valid
_ = err
```

---

## cipher - Encryption

Symmetric encryption helpers.

```go
import "github.com/kainonly/go/cipher"

c := cipher.New(cipher.SetKey("your-32-byte-key"))

encrypted, err := c.Encode(data)
_ = encrypted
_ = err

decrypted, err := c.Decode(encrypted)
_ = decrypted
_ = err
```

---

## help - Helpers

Misc helper functions.

### Random

```go
import "github.com/kainonly/go/help"

help.Random(16)
help.RandomNumber(6)
help.RandomAlphabet(8)
help.RandomUppercase(8)
help.RandomLowercase(8)
```

### ID

```go
help.Uuid7()
help.MustUuid7()
help.Uuid7Time(id)
help.Uuid()
help.SID()
```

> **Why UUIDv7?**
> - Time-ordered IDs; new records naturally come later
> - Better index locality for databases with sequential inserts
> - Extract creation time from the ID

### Crypto helpers

```go
help.Sha256hex(s)
help.HmacSha256(s, key)
```

### SM2/SM4 (GM/T)

```go
sig, _ := help.Sm2Sign(privateKey, data)
valid := help.Sm2Verify(publicKey, data, sig)
_ = valid

encrypted, _ := help.SM4Encrypt(key, plaintext)
decrypted, _ := help.SM4Decrypt(key, encrypted)
_ = decrypted
```

### Misc

```go
help.Ptr(value)
help.IsEmpty(value)
help.Reverse(slice)
help.Shuffle(slice)
```

## License

[BSD-3-Clause License](https://github.com/kainonly/go/blob/main/LICENSE)
