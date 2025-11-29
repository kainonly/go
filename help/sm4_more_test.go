package help

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSM4_InvalidInputs(t *testing.T) {
	_, err := SM4Encrypt("zzz", "hello")
	assert.Error(t, err)
	// invalid hex ciphertext for decrypt
	_, err = SM4Decrypt(hex.EncodeToString(make([]byte, 16)), "not-hex")
	assert.Error(t, err)
}
