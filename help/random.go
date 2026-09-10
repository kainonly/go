// Package help 提供 Go 应用程序的常用工具函数。
//
// 包含：
//   - 随机字符串生成（密码学安全）
//   - UUID 与 Snowflake ID 生成
//   - 切片/字符串操作（反转、乱序）
//   - Map 转查询字符串
//   - SM2/SM4 加密工具
//   - Hertz 框架集成辅助
package help

import "crypto/rand"

// Random 生成长度为 n 的密码学安全随机字符串。
// 默认使用字母数字字符（a-zA-Z0-9）。
// 可选参数 charset 用于自定义字符集。
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

// RandomNumber 生成长度为 n 的随机数字字符串。
// 仅使用数字 0-9。
func RandomNumber(n int) string {
	return Random(n, "0123456789")
}

// RandomLowercase 生成长度为 n 的随机小写字母字符串。
func RandomLowercase(n int) string {
	return Random(n, "abcdefghijklmnopqrstuvwxyz")
}

// RandomUppercase 生成长度为 n 的随机大写字母字符串。
func RandomUppercase(n int) string {
	return Random(n, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// RandomAlphabet 生成长度为 n 的随机字母字符串。
// 同时使用大写和小写字母。
func RandomAlphabet(n int) string {
	return Random(n, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}
