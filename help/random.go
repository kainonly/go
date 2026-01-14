// Package help provides common utility functions for Go applications.
//
// It includes:
//   - Random string generation (cryptographically secure)
//   - UUID and Snowflake ID generation
//   - Slice/string manipulation (reverse, shuffle)
//   - Map to query string conversion
//   - SM2/SM4 cryptographic utilities
//   - Hertz framework integration helpers
package help

import "crypto/rand"

// Random generates a cryptographically secure random string of length n.
// By default, it uses alphanumeric characters (a-zA-Z0-9).
// An optional charset can be provided to customize the character set.
func Random(n int, charset ...string) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if len(charset) != 0 {
		letters = charset[0]
	}

	b := make([]byte, n)
	rand.Read(b)
	letterLen := len(letters)
	for i := range b {
		b[i] = letters[int(b[i])%letterLen]
	}

	return string(b)
}

// RandomNumber generates a random numeric string of length n.
// Uses only digits 0-9.
func RandomNumber(n int) string {
	return Random(n, "0123456789")
}

// RandomLowercase generates a random lowercase alphabetic string of length n.
func RandomLowercase(n int) string {
	return Random(n, "abcdefghijklmnopqrstuvwxyz")
}

// RandomUppercase generates a random uppercase alphabetic string of length n.
func RandomUppercase(n int) string {
	return Random(n, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// RandomAlphabet generates a random alphabetic string of length n.
// Uses both uppercase and lowercase letters.
func RandomAlphabet(n int) string {
	return Random(n, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}
