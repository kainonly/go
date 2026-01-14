package help

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
)

// Reverse reverses the order of elements in a slice in place.
func Reverse[T any](v []T) {
	for n, m := 0, len(v)-1; n < len(v)/2; n, m = n+1, m-1 {
		v[n], v[m] = v[m], v[n]
	}
}

// Shuffle randomly shuffles the elements in a slice in place.
// Uses cryptographically secure random numbers.
func Shuffle[T any](v []T) {
	for n := len(v) - 1; n > 0; n-- {
		m := secureRandInt(n + 1)
		if n != m {
			v[n], v[m] = v[m], v[n]
		}
	}
}

// ReverseString returns a new string with characters in reverse order.
// Properly handles Unicode characters.
func ReverseString(v string) string {
	runes := []rune(v)
	for n, m := 0, len(runes)-1; n < len(runes)/2; n, m = n+1, m-1 {
		runes[n], runes[m] = runes[m], runes[n]
	}
	return string(runes)
}

// ShuffleString returns a new string with characters randomly shuffled.
// Uses cryptographically secure random numbers.
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

// secureRandInt returns a cryptographically secure random int in [0, max).
func secureRandInt(max int) int {
	var b [8]byte
	rand.Read(b[:])
	return int(binary.BigEndian.Uint64(b[:]) % uint64(max))
}

// MapToSignText converts a map to a URL-encoded query string format.
// Keys are sorted alphabetically, nil and empty values are omitted.
// Format: "key1=value1&key2=value2"
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
