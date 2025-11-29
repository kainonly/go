package help_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestSha256hex(t *testing.T) {
	expected := sha256.Sum256([]byte("hello"))
	assert.Equal(t, hex.EncodeToString(expected[:]), help.Sha256hex("hello"))
}

func TestHmacSha256(t *testing.T) {
	mac := hmac.New(sha256.New, []byte("key"))
	mac.Write([]byte("hello"))
	expected := string(mac.Sum(nil))
	assert.Equal(t, expected, help.HmacSha256("hello", "key"))
}
