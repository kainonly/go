// Package passlib 使用 Argon2id 算法提供密码哈希。
// Argon2id 是 OWASP 推荐的密码哈希算法。
package passlib

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id 哈希的默认参数，遵循 OWASP 建议。
var (
	DefaultMemoryCost uint32 = 65536 // 64 MB
	DefaultTimeCost   uint32 = 4
	DefaultThreads    uint8  = 1
)

// passlib 函数返回的错误。
var (
	ErrInvalidHash         = errors.New("passlib: unable to parse hash value")
	ErrIncompatibleVariant = errors.New("passlib: hash variant is not compatible")
	ErrIncompatibleVersion = errors.New("passlib: hash version is not supported")
	ErrNotMatch            = errors.New("passlib: password does not match hash")
)

// Hash 使用默认参数生成密码的 Argon2id 哈希。
// 返回包含算法、版本、参数、盐值和哈希的 PHC 格式字符串。
func Hash(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt,
		DefaultTimeCost, DefaultMemoryCost, DefaultThreads, 32)

	return fmt.Sprintf(`$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s`,
		argon2.Version, DefaultMemoryCost, DefaultTimeCost, DefaultThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

// Verify 校验密码是否与给定哈希匹配。
// 匹配返回 nil，否则返回错误。
func Verify(password string, hash string) error {
	memory, time, threads, salt, key, err := parseHash(hash)
	if err != nil {
		return err
	}
	otherKey := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(key)))
	if subtle.ConstantTimeCompare(key, otherKey) != 1 {
		return ErrNotMatch
	}
	return nil
}

// NeedsRehash 判断哈希是否由过时的参数生成、
// 需要用当前默认参数重新哈希。
func NeedsRehash(hash string) bool {
	memory, time, threads, _, _, err := parseHash(hash)
	if err != nil {
		return true
	}
	return memory != DefaultMemoryCost || time != DefaultTimeCost || threads != DefaultThreads
}

// parseHash 从 PHC 格式的 Argon2id 哈希字符串中解析参数。
func parseHash(hash string) (memory, time uint32, threads uint8, salt, key []byte, err error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		err = ErrInvalidHash
		return
	}
	if parts[1] != "argon2id" {
		err = ErrIncompatibleVariant
		return
	}
	var version int
	if _, err = fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		err = ErrIncompatibleVersion
		return
	}
	if version != argon2.Version {
		err = ErrIncompatibleVersion
		return
	}
	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		err = ErrInvalidHash
		return
	}
	if salt, err = base64.RawStdEncoding.Strict().DecodeString(parts[4]); err != nil {
		err = ErrInvalidHash
		return
	}
	if key, err = base64.RawStdEncoding.Strict().DecodeString(parts[5]); err != nil {
		err = ErrInvalidHash
		return
	}
	return
}
