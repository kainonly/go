# Go Utils

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/kainonly/go/testing.yml?style=flat-square)](https://github.com/kainonly/go/actions/workflows/testing.yml)
[![Coveralls github](https://img.shields.io/coveralls/github/kainonly/go.svg?style=flat-square)](https://coveralls.io/github/kainonly/go)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainonly/go?style=flat-square)](https://github.com/kainonly/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/kainonly/go?style=flat-square)](https://goreportcard.com/report/github.com/kainonly/go)
[![Release](https://img.shields.io/github/v/release/kainonly/go.svg?style=flat-square)](https://github.com/kainonly/go)
[![GitHub license](https://img.shields.io/github/license/kainonly/go?style=flat-square)](https://raw.githubusercontent.com/kainonly/go/main/LICENSE)

简体中文 | [English](README.md)

Go 常用工具库集合，提供 Web 开发中常用的功能组件。

## 安装

```bash
go get github.com/kainonly/go
```

## 模块

| 模块 | 描述 |
|------|------|
| [vd](#vd---验证器) | 验证器，集成 Hertz 框架 |
| [passport](#passport---jwt-认证) | JWT 认证 (HS256) |
| [csrf](#csrf---csrf-防护) | CSRF 防护中间件 |
| [captcha](#captcha---验证码) | 验证码管理 (Redis) |
| [locker](#locker---失败锁定) | 尝试次数计数与锁定 (Redis) |
| [passlib](#passlib---密码哈希) | 密码哈希 (Argon2id) |
| [totp](#totp---一次性密码) | TOTP 一次性密码 |
| [cipher](#cipher---加密) | 对称加密工具 |
| [help](#help---工具函数) | 公共工具函数 |

---

## vd - 验证器

基于 [go-playground/validator](https://github.com/go-playground/validator) 的验证器封装，支持 Hertz 框架集成和自定义验证规则。

### 基础用法

```go
import "github.com/kainonly/go/vd"

// 使用默认规则 (snake, sort)
v := vd.Default()

// 自定义规则
v := vd.New(
    vd.SetTag("vd"),  // 可选，默认为 "vd"
    vd.SetRules(
        vd.Snake(),
        vd.Sort(),
        vd.Phone(),
        vd.IDCard(),
        vd.PasswordMedium(),
    ),
)

// 验证结构体
type User struct {
    Name  string `vd:"required,snake"`
    Phone string `vd:"required,phone"`
}
err := v.Validate(&User{Name: "user_name", Phone: "13800138000"})
```

### Hertz 集成

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

    h.POST("/user", func(ctx context.Context, c *app.RequestContext) {
        var req struct {
            Name  string `json:"name" vd:"required,snake"`
            Phone string `json:"phone" vd:"required,phone"`
        }
        if err := c.BindAndValidate(&req); err != nil {
            c.JSON(400, map[string]any{"error": err.Error()})
            return
        }
        c.JSON(200, map[string]any{"message": "ok"})
    })

    h.Spin()
}
```

### 可用验证规则

#### 规则分组

| 函数 | 描述 |
|------|------|
| `All()` | 所有规则 |
| `Common()` | 常用规则 (snake, sort, phone, idcard, username, slug, password_medium) |
| `Chinese()` | 中国本地化规则 |
| `NamingConvention()` | 命名规范规则 |

#### 中国本地化

| 规则 | Tag | 描述 | 示例 |
|------|-----|------|------|
| `Phone()` | `phone` | 手机号 | `13800138000` |
| `IDCard()` | `idcard` | 身份证号 | `110101199003071234` |
| `BankCard()` | `bankcard` | 银行卡号 (Luhn) | `6222021234567890` |
| `LicensePlate()` | `license_plate` | 车牌号 | `京A12345` |
| `USCC()` | `uscc` | 统一社会信用代码 | `91310000MA1FL8TQ32` |
| `ChineseWord()` | `chinese` | 纯中文 | `你好世界` |
| `ChineseName()` | `chinese_name` | 中文姓名 | `张三` |
| `TelPhone()` | `tel` | 固定电话 | `010-12345678` |
| `QQ()` | `qq` | QQ号 | `12345678` |
| `WeChat()` | `wechat` | 微信号 | `wxid_abc123` |
| `ZipCode()` | `zipcode` | 邮政编码 | `518000` |

#### 密码强度

| 规则 | Tag | 描述 |
|------|-----|------|
| `PasswordWeak()` | `password_weak` | ≥6 字符 |
| `PasswordMedium()` | `password_medium` | ≥8 字符，含字母和数字 |
| `PasswordStrong()` | `password_strong` | ≥8 字符，含大小写、数字、特殊字符 |

#### 命名规范

| 规则 | Tag | 示例 |
|------|-----|------|
| `Snake()` | `snake` | `user_name` |
| `PascalCase()` | `pascal` | `UserName` |
| `CamelCase()` | `camel` | `userName` |
| `KebabCase()` | `kebab` | `user-name` |
| `UpperSnake()` | `upper_snake` | `USER_NAME` |
| `Variable()` | `variable` | `_privateVar` |

#### 其他规则

| 规则 | Tag | 描述 |
|------|-----|------|
| `Sort()` | `sort` | 排序格式 `field:1` 或 `field:-1` |
| `Username()` | `username` | 用户名 (3-20字符) |
| `Slug()` | `slug` | URL slug |
| `ObjectID()` | `objectid` | MongoDB ObjectId |
| `Snowflake()` | `snowflake` | 雪花ID |
| `Version()` | `version` | 版本号 `1.0.0` |
| `SafeString()` | `safe_string` | 防注入安全字符串 |
| `AlphaNumDash()` | `alphanumdash` | 字母数字下划线横线 |
| `Decimal()` | `decimal` | 十进制数字符串 |
| `PositiveDecimal()` | `positive_decimal` | 正十进制数 |
| `Domain()` | `domain` | 域名 |
| `FileName()` | `filename` | 文件名 |
| `FileExt()` | `file_ext` | 文件扩展名 |
| `FilePath()` | `file_path` | 文件路径 |
| `Color()` | `color` | 颜色值 (无#) |
| `NotBlank()` | `notblank` | 非空白字符串 |

#### 自定义规则

```go
v := vd.New(vd.SetRules(
    vd.Rule{
        Tag: "even",
        Fn: func(fl vd.FieldLevel) bool {
            return fl.Field().Int() % 2 == 0
        },
    },
))
```

---

## passport - JWT 认证

基于 HS256 的 JWT 认证。

```go
import "github.com/kainonly/go/passport"

// 创建认证器
auth := passport.New(
    passport.SetKey("your-secret-key"),
    passport.SetIssuer("your-app"),
)

// 创建 Token
claims := passport.NewClaims()
claims.ID = "user-123"
claims.SetData(map[string]any{"role": "admin"})
token, err := auth.Create(claims)

// 验证 Token
claims, err := auth.Verify(token)
```

---

## csrf - CSRF 防护

双重 Cookie 验证的 CSRF 防护中间件。

```go
import "github.com/kainonly/go/csrf"

// 创建 CSRF 防护
c := csrf.New(
    csrf.SetKey("your-secret-key"),
    csrf.SetDomain("example.com"),
)

// Hertz 中间件
h.Use(c.VerifyToken())

// 设置 Token Cookie
h.GET("/csrf", func(ctx context.Context, c *app.RequestContext) {
    csrf.SetToken(c)
})
```

---

## captcha - 验证码

基于 Redis 的验证码管理。

```go
import "github.com/kainonly/go/captcha"

// 创建验证码管理器
cap := captcha.New(redisClient)

// 创建验证码
cap.Create(ctx, "login:user123", "123456", 5*time.Minute)

// 验证
err := cap.Verify(ctx, "login:user123", "123456")
if errors.Is(err, captcha.ErrInvalidCode) {
    // 验证码错误
}
if errors.Is(err, captcha.ErrNotExists) {
    // 验证码不存在或已过期
}
```

---

## locker - 失败锁定

基于 Redis 的尝试次数计数与锁定，适用于登录失败锁定、验证码错误限制等场景。

```go
import "github.com/kainonly/go/locker"

// 创建锁定器
lock := locker.New(redisClient)

// 登录失败时增加计数
count, err := lock.Increment(ctx, "login:user123", time.Minute)

// 检查是否已锁定（超过 5 次失败）
err = lock.Check(ctx, "login:user123", 5)
if errors.Is(err, locker.ErrLocked) {
    // 账户已锁定，请稍后再试
}

// 登录成功后清除计数
lock.Delete(ctx, "login:user123")
```

---

## passlib - 密码哈希

基于 Argon2id 的密码哈希。

```go
import "github.com/kainonly/go/passlib"

// 哈希密码
hash, err := passlib.Hash("password123")

// 验证密码
err = passlib.Verify("password123", hash)

// 检查是否需要重新哈希
if passlib.NeedsRehash(hash) {
    newHash, _ := passlib.Hash("password123")
}
```

---

## totp - 一次性密码

TOTP 一次性密码生成与验证。

```go
import "github.com/kainonly/go/totp"

// 生成密钥
secret, err := totp.GenerateSecret()

// 验证 OTP
valid, err := totp.Validate(code, secret)
```

---

## cipher - 加密

对称加密工具。

```go
import "github.com/kainonly/go/cipher"

// 创建加密器
c := cipher.New(cipher.SetKey("your-32-byte-key"))

// 加密
encrypted, err := c.Encode(data)

// 解密
decrypted, err := c.Decode(encrypted)
```

---

## help - 工具函数

公共工具函数集合。

### 随机生成

```go
import "github.com/kainonly/go/help"

help.Random(16)           // 随机字符串
help.RandomNumber(6)      // 随机数字
help.RandomAlphabet(8)    // 随机字母
help.RandomUppercase(8)   // 随机大写
help.RandomLowercase(8)   // 随机小写
```

### ID 生成

```go
help.Uuid7()              // UUID v7 (推荐用于主键)
help.MustUuid7()          // UUID v7 (失败时 panic)
help.Uuid7Time(id)        // 从 UUID v7 提取时间戳
help.Uuid()               // UUID v4 (已废弃，建议使用 v7)
help.SID()                // Snowflake ID
```

> **为什么使用 UUIDv7？**
> - 按时间排序，新记录自然排在后面
> - 顺序插入，数据库索引性能更好
> - 可从 ID 中提取创建时间

### 加密工具

```go
help.Sha256hex(s)         // SHA256 哈希
help.HmacSha256(s, key)   // HMAC-SHA256
```

### 国密 SM2/SM4

```go
// SM2 签名验签
sig, _ := help.Sm2Sign(privateKey, data)
valid := help.Sm2Verify(publicKey, data, sig)

// SM4 加解密
encrypted, _ := help.SM4Encrypt(key, plaintext)
decrypted, _ := help.SM4Decrypt(key, encrypted)
```

### 其他工具

```go
help.Ptr(value)           // 获取值的指针
help.IsEmpty(value)       // 检查是否为空
help.Reverse(slice)       // 反转切片
help.Shuffle(slice)       // 打乱切片
```

---

## License

[BSD-3-Clause License](https://github.com/kainonly/go/blob/main/LICENSE)
