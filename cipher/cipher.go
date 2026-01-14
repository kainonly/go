// Package cipher provides symmetric encryption using XChaCha20-Poly1305.
// It is designed for encrypting sensitive data before storing in databases.
package cipher

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	// ErrCiphertextTooShort is returned when the ciphertext is shorter than the nonce size.
	ErrCiphertextTooShort = errors.New("cipher: ciphertext too short")
)

// Cipher wraps XChaCha20-Poly1305 AEAD for encryption and decryption.
type Cipher struct {
	AEAD cipher.AEAD
}

// New creates a new Cipher with the given key.
// The key must be exactly 32 bytes for XChaCha20-Poly1305.
func New(key string) (*Cipher, error) {
	aead, err := chacha20poly1305.NewX([]byte(key))
	if err != nil {
		return nil, err
	}
	return &Cipher{AEAD: aead}, nil
}

// Encode encrypts the plaintext data and returns a base64-encoded ciphertext.
// A random nonce is generated for each encryption and prepended to the ciphertext.
func (x *Cipher) Encode(data []byte) (string, error) {
	nonce := make([]byte, x.AEAD.NonceSize(), x.AEAD.NonceSize()+len(data)+x.AEAD.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	encrypted := x.AEAD.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decode decrypts a base64-encoded ciphertext and returns the original plaintext.
// Returns ErrCiphertextTooShort if the ciphertext is invalid.
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
