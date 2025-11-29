package help

import (
	"testing"

	"github.com/emmansun/gmsm/sm4"
	"github.com/stretchr/testify/assert"
)

func TestNewECB(t *testing.T) {
	key := make([]byte, 16)
	block, err := sm4.NewCipher(key)
	assert.NoError(t, err)
	modeEnc := newECBEncrypter(block)
	modeDec := newECBDecrypter(block)
	assert.Equal(t, block.BlockSize(), modeEnc.BlockSize())
	assert.Equal(t, block.BlockSize(), modeDec.BlockSize())
}

func TestECB_CryptBlocks_PanicInputs(t *testing.T) {
	key := make([]byte, 16)
	block, err := sm4.NewCipher(key)
	assert.NoError(t, err)
	mode := newECBEncrypter(block)
	assert.Panics(t, func() { mode.CryptBlocks(make([]byte, 15), make([]byte, 15)) })
	assert.Panics(t, func() { mode.CryptBlocks(make([]byte, 15), make([]byte, 16)) })
}

func TestPKCS5PaddingUnpadding(t *testing.T) {
	b := []byte("abc")
	padded := pkcs5Padding(b, 8)
	assert.Equal(t, 8, len(padded))
	unpadded, err := pkcs5UnPadding(padded)
	assert.NoError(t, err)
	assert.Equal(t, string(b), string(unpadded))

	_, err = pkcs5UnPadding([]byte{})
	assert.Error(t, err)
}
