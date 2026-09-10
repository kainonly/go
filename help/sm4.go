package help

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"errors"

	"github.com/emmansun/gmsm/sm4"
)

// SM4 相关错误。
var (
	ErrSM4InvalidKey        = errors.New("sm4: invalid key, must be 32 hex characters (16 bytes)")
	ErrSM4InvalidCiphertext = errors.New("sm4: invalid ciphertext")
	ErrSM4InvalidPadding    = errors.New("sm4: invalid PKCS5 padding")
	ErrSM4EmptyData         = errors.New("sm4: data is empty")
)

// SM4Encrypt 使用 SM4-ECB 模式加 PKCS5 填充加密明文。
// 密钥必须是 32 个字符的十六进制字符串（16 字节）。
// 返回十六进制编码的密文。
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

// SM4Decrypt 使用 SM4-ECB 模式加 PKCS5 填充解密密文。
// 密钥必须是 32 个字符的十六进制字符串（16 字节）。
// 密文必须是十六进制编码。
// 返回解密后的明文。
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

// SM4Verify 解密密文并与预期明文比较。
// 解密后的文本与明文匹配时返回 true。
func SM4Verify(key string, ciphertext string, plaintext string) (bool, error) {
	decryptText, err := SM4Decrypt(key, ciphertext)
	if err != nil {
		return false, err
	}
	return decryptText == plaintext, nil
}

// ecb 为分组密码实现 ECB（Electronic Codebook）模式。
type ecb struct {
	b         cipher.Block
	blockSize int
}

// newECB 创建新的 ECB 模式实例。
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

// pkcs5Padding 为数据添加 PKCS5/PKCS7 填充。
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs5UnPadding 移除数据的 PKCS5/PKCS7 填充。
func pkcs5UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrSM4EmptyData
	}
	unpadding := int(data[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, ErrSM4InvalidPadding
	}
	// 校验填充字节
	for i := length - unpadding; i < length; i++ {
		if data[i] != byte(unpadding) {
			return nil, ErrSM4InvalidPadding
		}
	}
	return data[:length-unpadding], nil
}
