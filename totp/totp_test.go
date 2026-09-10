package totp_test

import (
	"testing"

	"github.com/kainonly/go/totp"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSecret(t *testing.T) {
	key, err := totp.GenerateSecret(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "user@example.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Equal(t, "TestApp", key.Issuer())
	assert.Equal(t, "user@example.com", key.AccountName())
	assert.NotEmpty(t, key.Secret())
	assert.NotEmpty(t, key.URL())
}

func TestGenerateAndValidate(t *testing.T) {
	key, err := totp.GenerateSecret(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "user@example.com",
	})
	assert.NoError(t, err)

	// 生成验证码
	code, err := totp.Generate(key.Secret())
	assert.NoError(t, err)
	assert.Len(t, code, 6)

	// 校验验证码
	valid := totp.Validate(code, key.Secret())
	assert.True(t, valid)
}

func TestValidate_InvalidCode(t *testing.T) {
	key, err := totp.GenerateSecret(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "user@example.com",
	})
	assert.NoError(t, err)

	valid := totp.Validate("000000", key.Secret())
	assert.False(t, valid)
}

func TestValidateWithOpts(t *testing.T) {
	key, err := totp.GenerateSecret(totp.GenerateOpts{
		Issuer:      "TestApp",
		AccountName: "user@example.com",
	})
	assert.NoError(t, err)

	code, err := totp.Generate(key.Secret())
	assert.NoError(t, err)

	// 以更大的时间窗口校验
	valid, err := totp.ValidateWithOpts(code, key.Secret(), totp.ValidateOpts{
		Skew: 2,
	})
	assert.NoError(t, err)
	assert.True(t, valid)

	// 无效的验证码应失败
	valid, err = totp.ValidateWithOpts("000000", key.Secret(), totp.ValidateOpts{
		Skew: 1,
	})
	assert.NoError(t, err)
	assert.False(t, valid)
}
