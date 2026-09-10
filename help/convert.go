package help

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
)

// Reverse 原地反转切片中元素的顺序。
func Reverse[T any](v []T) {
	for n, m := 0, len(v)-1; n < len(v)/2; n, m = n+1, m-1 {
		v[n], v[m] = v[m], v[n]
	}
}

// Shuffle 原地随机打乱切片中的元素。
// 使用密码学安全的随机数。
func Shuffle[T any](v []T) {
	for n := len(v) - 1; n > 0; n-- {
		m := secureRandInt(n + 1)
		if n != m {
			v[n], v[m] = v[m], v[n]
		}
	}
}

// ReverseString 返回字符顺序反转后的新字符串。
// 正确处理 Unicode 字符。
func ReverseString(v string) string {
	runes := []rune(v)
	for n, m := 0, len(runes)-1; n < len(runes)/2; n, m = n+1, m-1 {
		runes[n], runes[m] = runes[m], runes[n]
	}
	return string(runes)
}

// ShuffleString 返回字符随机打乱后的新字符串。
// 使用密码学安全的随机数。
func ShuffleString(v string) string {
	runes := []rune(v)
	for n := len(runes) - 1; n > 0; n-- {
		m := secureRandInt(n + 1)
		if n != m {
			runes[n], runes[m] = runes[m], runes[n]
		}
	}
	return string(runes)
}

// secureRandInt 返回 [0, max) 范围内密码学安全的随机整数。
func secureRandInt(max int) int {
	var b [8]byte
	rand.Read(b[:])
	return int(binary.BigEndian.Uint64(b[:]) % uint64(max))
}

// MapToSignText 将 map 转换为 URL 编码的查询字符串格式。
// 键按字母顺序排序，nil 和空值会被省略。
// 格式："key1=value1&key2=value2"
func MapToSignText(d map[string]any) string {
	keys := make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	first := true
	for _, k := range keys {
		v := d[k]
		if v == nil {
			continue
		}

		var strVal string
		switch val := v.(type) {
		case string:
			strVal = val
		case int:
			strVal = strconv.Itoa(val)
		case int64:
			strVal = strconv.FormatInt(val, 10)
		case float64:
			strVal = strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			strVal = strconv.FormatBool(val)
		default:
			strVal = fmt.Sprintf("%v", val)
		}

		if strVal == "" {
			continue
		}

		if !first {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(strVal)
		first = false
	}
	return buf.String()
}
