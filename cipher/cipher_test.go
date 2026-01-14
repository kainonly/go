package cipher_test

import (
	"testing"

	"github.com/kainonly/go/cipher"
	"github.com/stretchr/testify/assert"
)

var x1 *cipher.Cipher
var x2 *cipher.Cipher

func TestUseCipher(t *testing.T) {
	var err error
	_, err = cipher.New("123456")
	assert.Error(t, err)
	x1, err = cipher.New("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK")
	assert.NoError(t, err)
	x2, err = cipher.New("74rILbVooYLirHrQJcslHEAvKZI7PKF9")
	assert.NoError(t, err)
}

var text = "Gophers, gophers, gophers everywhere!"
var encryptedText string

func TestCipher_Encode(t *testing.T) {
	var err error
	encryptedText, err = x1.Encode([]byte(text))
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedText)

	// Test empty data
	emptyCt, err := x1.Encode([]byte{})
	assert.NoError(t, err)
	assert.NotEmpty(t, emptyCt)
}

func TestCipher_Decode(t *testing.T) {
	decryptedText, err := x1.Decode(encryptedText)
	assert.NoError(t, err)
	assert.Equal(t, text, string(decryptedText))

	// Wrong key should fail
	_, err = x2.Decode(encryptedText)
	assert.Error(t, err)

	// Invalid base64 should fail
	_, err = x1.Decode("@@@notbase64@@@")
	assert.Error(t, err)

	// Ciphertext too short should fail
	_, err = x1.Decode("YWJj") // "abc" in base64
	assert.ErrorIs(t, err, cipher.ErrCiphertextTooShort)
}
