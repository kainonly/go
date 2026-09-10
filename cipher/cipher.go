// Package cipher 提供基于 XChaCha20-Poly1305 的对称加密。
// 专为敏感数据写入数据库前的加密场景而设计。
package cipher

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	// ErrCiphertextTooShort 在密文长度小于 nonce 长度时返回。
	ErrCiphertextTooShort = errors.New("cipher: ciphertext too short")
)

// Cipher 封装 XChaCha20-Poly1305 AEAD，用于加密和解密。
type Cipher struct {
	AEAD cipher.AEAD
}

// New 使用给定的密钥创建一个新的 Cipher。
// 密钥必须恰好为 32 字节，这是 XChaCha20-Poly1305 的要求。
func New(key string) (*Cipher, error) {
	aead, err := chacha20poly1305.NewX([]byte(key))
	if err != nil {
		return nil, err
	}
	return &Cipher{AEAD: aead}, nil
}

// Encode 加密明文数据，返回 base64 编码的密文。
// 每次加密都会生成随机 nonce，并将其附在密文之前。
func (x *Cipher) Encode(data []byte) (string, error) {
	nonce := make([]byte, x.AEAD.NonceSize(), x.AEAD.NonceSize()+len(data)+x.AEAD.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	encrypted := x.AEAD.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decode 解密 base64 编码的密文，返回原始明文。
// 如果密文无效，返回 ErrCiphertextTooShort。
func (x *Cipher) Decode(ciphertext string) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < x.AEAD.NonceSize() {
		return nil, ErrCiphertextTooShort
	}
	nonce, text := encrypted[:x.AEAD.NonceSize()], encrypted[x.AEAD.NonceSize():]
	return x.AEAD.Open(nil, nonce, text, nil)
}
