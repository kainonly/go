// Package totp provides TOTP (Time-based One-Time Password) generation and validation.
// It wraps github.com/pquerna/otp for RFC 6238 compliant implementation.
package totp

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// Key represents a TOTP key with its configuration.
type Key struct {
	*otp.Key
}

// GenerateOpts configures the key generation.
type GenerateOpts struct {
	// Issuer is the name of the issuing organization (e.g., "MyApp").
	Issuer string
	// AccountName is the user's identifier (e.g., "user@example.com").
	AccountName string
	// SecretSize is the size of the secret in bytes. Default is 20.
	SecretSize uint
	// Algorithm is the hash algorithm. Default is SHA1.
	Algorithm otp.Algorithm
	// Digits is the number of digits in the OTP. Default is 6.
	Digits otp.Digits
}

// GenerateSecret creates a new TOTP secret key.
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

// Validate checks if the passcode is valid for the given secret.
// Uses a default time window of ±1 period (30 seconds).
func Validate(passcode, secret string) bool {
	return totp.Validate(passcode, secret)
}

// ValidateOpts configures the validation behavior.
type ValidateOpts struct {
	// Skew is the number of periods before/after current time to check. Default is 1.
	Skew uint
	// Digits is the expected number of digits. Default is 6.
	Digits otp.Digits
	// Algorithm is the hash algorithm. Default is SHA1.
	Algorithm otp.Algorithm
}

// ValidateWithOpts checks if the passcode is valid with custom options.
func ValidateWithOpts(passcode, secret string, opts ValidateOpts) (bool, error) {
	valOpts := totp.ValidateOpts{
		Skew:      opts.Skew,
		Digits:    otp.DigitsSix,     // Default to 6 digits
		Algorithm: otp.AlgorithmSHA1, // Default to SHA1
	}
	if opts.Digits != 0 {
		valOpts.Digits = opts.Digits
	}
	if opts.Algorithm != 0 {
		valOpts.Algorithm = opts.Algorithm
	}
	return totp.ValidateCustom(passcode, secret, time.Now(), valOpts)
}

// Generate creates a TOTP code for the given secret at the current time.
func Generate(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
