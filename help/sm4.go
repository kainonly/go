package help

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"errors"

	"github.com/emmansun/gmsm/sm4"
)

// SM4 related errors.
var (
	ErrSM4InvalidKey        = errors.New("sm4: invalid key, must be 32 hex characters (16 bytes)")
	ErrSM4InvalidCiphertext = errors.New("sm4: invalid ciphertext")
	ErrSM4InvalidPadding    = errors.New("sm4: invalid PKCS5 padding")
	ErrSM4EmptyData         = errors.New("sm4: data is empty")
)

// SM4Encrypt encrypts plaintext using SM4-ECB mode with PKCS5 padding.
// The key must be a 32-character hex string (16 bytes).
// Returns hex-encoded ciphertext.
func SM4Encrypt(hexkey string, plaintext string) (string, error) {
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		return "", err
	}

	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := []byte(plaintext)
	b = pkcs5Padding(b, block.BlockSize())

	content := make([]byte, len(b))
	mode := newECBEncrypter(block)
	mode.CryptBlocks(content, b)

	return hex.EncodeToString(content), nil
}

// SM4Decrypt decrypts ciphertext using SM4-ECB mode with PKCS5 padding.
// The key must be a 32-character hex string (16 bytes).
// The ciphertext must be hex-encoded.
// Returns the decrypted plaintext.
func SM4Decrypt(hexkey string, ciphertext string) (string, error) {
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		return "", err
	}

	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}

	content, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(content) == 0 || len(content)%block.BlockSize() != 0 {
		return "", ErrSM4InvalidCiphertext
	}

	b := make([]byte, len(content))
	mode := newECBDecrypter(block)
	mode.CryptBlocks(b, content)

	unpadding, err := pkcs5UnPadding(b)
	if err != nil {
		return "", err
	}

	return string(unpadding), nil
}

// SM4Verify decrypts and compares ciphertext with expected plaintext.
// Returns true if the decrypted text matches the plaintext.
func SM4Verify(key string, ciphertext string, plaintext string) (bool, error) {
	decryptText, err := SM4Decrypt(key, ciphertext)
	if err != nil {
		return false, err
	}
	return decryptText == plaintext, nil
}

// ecb implements ECB (Electronic Codebook) mode for block ciphers.
type ecb struct {
	b         cipher.Block
	blockSize int
}

// newECB creates a new ECB mode instance.
func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("sm4: input length must be a multiple of block size")
	}
	if len(dst) < len(src) {
		panic("sm4: output buffer too small")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("sm4: input length must be a multiple of block size")
	}
	if len(dst) < len(src) {
		panic("sm4: output buffer too small")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// pkcs5Padding adds PKCS5/PKCS7 padding to data.
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs5UnPadding removes PKCS5/PKCS7 padding from data.
func pkcs5UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrSM4EmptyData
	}
	unpadding := int(data[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, ErrSM4InvalidPadding
	}
	// Validate padding bytes
	for i := length - unpadding; i < length; i++ {
		if data[i] != byte(unpadding) {
			return nil, ErrSM4InvalidPadding
		}
	}
	return data[:length-unpadding], nil
}
