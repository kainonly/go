package cipher_test

import (
	"testing"

	"github.com/kainonly/go/cipher"
	"github.com/stretchr/testify/assert"
)

func TestDecode_BadBase64(t *testing.T) {
	x, err := cipher.New("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK")
	assert.NoError(t, err)
	_, err = x.Decode("@@@notbase64@@@")
	assert.Error(t, err)
}

func TestDecode_CiphertextTooShort(t *testing.T) {
	x, err := cipher.New("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK")
	assert.NoError(t, err)
	// Valid base64 but too short to contain nonce
	_, err = x.Decode("YWJj") // "abc" in base64
	assert.ErrorIs(t, err, cipher.ErrCiphertextTooShort)
}
