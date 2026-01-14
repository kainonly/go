// Package passlib provides password hashing using Argon2id algorithm.
// Argon2id is the recommended password hashing algorithm by OWASP.
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

// Default parameters for Argon2id hashing.
// These values follow OWASP recommendations.
var (
	DefaultMemoryCost uint32 = 65536 // 64 MB
	DefaultTimeCost   uint32 = 4
	DefaultThreads    uint8  = 1
)

// Errors returned by passlib functions.
var (
	ErrInvalidHash         = errors.New("passlib: unable to parse hash value")
	ErrIncompatibleVariant = errors.New("passlib: hash variant is not compatible")
	ErrIncompatibleVersion = errors.New("passlib: hash version is not supported")
	ErrNotMatch            = errors.New("passlib: password does not match hash")
)

// Hash generates an Argon2id hash of the password using default parameters.
// Returns a PHC-formatted string that includes the algorithm, version,
// parameters, salt, and hash.
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

// Verify checks if the password matches the given hash.
// Returns nil if the password matches, or an error otherwise.
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

// NeedsRehash checks if the hash was created with outdated parameters
// and should be rehashed with current default parameters.
func NeedsRehash(hash string) bool {
	memory, time, threads, _, _, err := parseHash(hash)
	if err != nil {
		return true
	}
	return memory != DefaultMemoryCost || time != DefaultTimeCost || threads != DefaultThreads
}

// parseHash extracts parameters from a PHC-formatted Argon2id hash string.
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
