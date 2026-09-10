package passlib_test

import (
	"testing"

	"github.com/kainonly/go/passlib"
	"github.com/stretchr/testify/assert"
)

func TestHashAndVerify(t *testing.T) {
	hash, err := passlib.Hash("pass@VAN1234")
	assert.NoError(t, err)
	err = passlib.Verify("pass@VAN1234", hash)
	assert.NoError(t, err)
	err = passlib.Verify("pass@VAN1235", hash)
	assert.ErrorIs(t, err, passlib.ErrNotMatch)

	// 测试超长密码
	longPw := make([]byte, 100)
	for i := range longPw {
		longPw[i] = 'a'
	}
	longHash, err := passlib.Hash(string(longPw))
	assert.NoError(t, err)
	assert.NotEmpty(t, longHash)
}

const PASS1 = `$argon2i$v=19$m=65536,t=4,p=1$NPCjKIcoU2z6rg6p8glOfg$jrbRcvsTq/ITJP414/xhNNwOtVeHYa478hPn8M6uJLA`
const PASS2 = `$argon2id$v=x$m=65536,t=4,p=1$NPCjKIcoU2z6rg6p8glOfg$jrbRcvsTq/ITJP414/xhNNwOtVeHYa478hPn8M6uJLA`
const PASS3 = `$argon2id$v=18$m=65536,t=4,p=1$NPCjKIcoU2z6rg6p8glOfg$jrbRcvsTq/ITJP414/xhNNwOtVeHYa478hPn8M6uJLA`
const PASS4 = `$argon2id$v=19$xcxcsdsdwe$NPCjKIcoU2z6rg6p8glOfg$jrbRcvsTq/ITJP414/xhNNwOtVeHYa478hPn8M6uJLA`
const PASS5 = `$argon2id$v=19$m=65536,t=4,p=1$()$jrbRcvsTq/ITJP414/xhNNwOtVeHYa478hPn8M6uJLA`
const PASS6 = `$argon2id$v=19$m=65536,t=4,p=1$NPCjKIcoU2z6rg6p8glOfg$()`

func TestVerifyErrors(t *testing.T) {
	var err error
	err = passlib.Verify("pass@VAN1234", "asdaqweqwexcxzcqweqw")
	assert.ErrorIs(t, err, passlib.ErrInvalidHash)
	err = passlib.Verify("pass@VAN1234", PASS1)
	assert.ErrorIs(t, err, passlib.ErrIncompatibleVariant)
	err = passlib.Verify("pass@VAN1234", PASS2)
	assert.ErrorIs(t, err, passlib.ErrIncompatibleVersion)
	err = passlib.Verify("pass@VAN1234", PASS3)
	assert.ErrorIs(t, err, passlib.ErrIncompatibleVersion)
	err = passlib.Verify("pass@VAN1234", PASS4)
	assert.ErrorIs(t, err, passlib.ErrInvalidHash)
	err = passlib.Verify("pass@VAN1234", PASS5)
	assert.ErrorIs(t, err, passlib.ErrInvalidHash)
	err = passlib.Verify("pass@VAN1234", PASS6)
	assert.ErrorIs(t, err, passlib.ErrInvalidHash)
}

func TestNeedsRehash(t *testing.T) {
	// 使用当前参数的哈希不需要重新哈希
	hash, err := passlib.Hash("password")
	assert.NoError(t, err)
	assert.False(t, passlib.NeedsRehash(hash))

	// 使用不同参数的哈希需要重新哈希
	oldHash := `$argon2id$v=19$m=32768,t=4,p=1$NPCjKIcoU2z6rg6p8glOfg$jrbRcvsTq/ITJP414/xhNNwOtVeHYa478hPn8M6uJLA`
	assert.True(t, passlib.NeedsRehash(oldHash))

	// 无效的哈希需要重新哈希
	assert.True(t, passlib.NeedsRehash("invalid"))
}
