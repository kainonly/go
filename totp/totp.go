// Package totp 提供 TOTP（基于时间的一次性密码）的生成与校验。
// 封装 github.com/pquerna/otp，符合 RFC 6238 规范。
package totp

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// Key 表示一个 TOTP 密钥及其配置。
type Key struct {
	*otp.Key
}

// GenerateOpts 用于配置密钥生成。
type GenerateOpts struct {
	// Issuer 是签发组织的名称（如 "MyApp"）。
	Issuer string
	// AccountName 是用户标识（如 "user@example.com"）。
	AccountName string
	// SecretSize 是密钥的字节长度，默认为 20。
	SecretSize uint
	// Algorithm 是哈希算法，默认为 SHA1。
	Algorithm otp.Algorithm
	// Digits 是 OTP 的位数，默认为 6。
	Digits otp.Digits
}

// GenerateSecret 创建一个新的 TOTP 密钥。
func GenerateSecret(opts GenerateOpts) (*Key, error) {
	genOpts := totp.GenerateOpts{
		Issuer:      opts.Issuer,
		AccountName: opts.AccountName,
	}
	if opts.SecretSize > 0 {
		genOpts.SecretSize = opts.SecretSize
	}
	if opts.Algorithm != 0 {
		genOpts.Algorithm = opts.Algorithm
	}
	if opts.Digits != 0 {
		genOpts.Digits = opts.Digits
	}
	key, err := totp.Generate(genOpts)
	if err != nil {
		return nil, err
	}
	return &Key{Key: key}, nil
}

// Validate 校验给定密钥的 passcode 是否有效。
// 默认时间窗口为 ±1 个周期（30 秒）。
func Validate(passcode, secret string) bool {
	return totp.Validate(passcode, secret)
}

// ValidateOpts 用于配置校验行为。
type ValidateOpts struct {
	// Skew 是当前时间前后允许校验的周期数，默认为 1。
	Skew uint
	// Digits 是期望的位数，默认为 6。
	Digits otp.Digits
	// Algorithm 是哈希算法，默认为 SHA1。
	Algorithm otp.Algorithm
}

// ValidateWithOpts 使用自定义选项校验 passcode 是否有效。
func ValidateWithOpts(passcode, secret string, opts ValidateOpts) (bool, error) {
	valOpts := totp.ValidateOpts{
		Skew:      opts.Skew,
		Digits:    otp.DigitsSix,     // 默认 6 位
		Algorithm: otp.AlgorithmSHA1, // 默认 SHA1
	}
	if opts.Digits != 0 {
		valOpts.Digits = opts.Digits
	}
	if opts.Algorithm != 0 {
		valOpts.Algorithm = opts.Algorithm
	}
	return totp.ValidateCustom(passcode, secret, time.Now(), valOpts)
}

// Generate 使用给定的密钥生成当前时间的 TOTP 验证码。
func Generate(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
