package help_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/smx509"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	v1 := help.Random(16)
	assert.Len(t, v1, 16)
	t.Log(v1)
	v2 := help.Random(32)
	assert.Len(t, v2, 32)
	t.Log(v2)
	v3 := help.Random(8)
	assert.Len(t, v3, 8)
	t.Log(v3)
	v4 := help.Random(32, "0123456789abcdef")
	assert.Len(t, v4, 32)
	t.Log(v4)
	v5 := help.Random(64, "0123456789abcdef")
	assert.Len(t, v5, 64)
	t.Log(v5)
}

func TestRandomNumber(t *testing.T) {
	v := help.RandomNumber(6)
	assert.Len(t, v, 6)
	t.Log(v)
}

func TestRandomAlphabet(t *testing.T) {
	v := help.RandomAlphabet(16)
	assert.Len(t, v, 16)
	t.Log(v)
}

func TestRandomUppercase(t *testing.T) {
	v := help.RandomUppercase(8)
	assert.Len(t, v, 8)
	t.Log(v)
}

func TestRandomLowercase(t *testing.T) {
	v := help.RandomLowercase(8)
	assert.Len(t, v, 8)
	t.Log(v)
}

func TestReverse(t *testing.T) {
	v := []string{"a", "b", "c"}
	help.Reverse(v)
	assert.Equal(t, []string{"c", "b", "a"}, v)
	t.Log(v)
}

func TestShuffle(t *testing.T) {
	v := []int{1, 2, 3, 4, 5, 6, 7}
	help.Shuffle(v)
	t.Log(v)
}

func TestReverseString(t *testing.T) {
	v := help.ReverseString("abcdefg")
	assert.Equal(t, "gfedcba", v)
	t.Log(v)
}

func TestShuffleString(t *testing.T) {
	v := help.ShuffleString("abcdefg")
	t.Log(v)
}

func TestMapToSignText(t *testing.T) {
	type Sign struct {
		input    map[string]any
		expected string
	}

	mocks := []Sign{
		{
			input:    map[string]any{},
			expected: "",
		},
		{
			input:    map[string]any{"key1": "value1"},
			expected: "key1=value1",
		},
		{
			input: map[string]any{
				"b": "2",
				"a": 1,
				"c": "3",
			},
			expected: "a=1&b=2&c=3",
		},
		{
			input: map[string]any{
				"key2": true,
				"key1": "value1",
				"key3": "value3",
			},
			expected: "key1=value1&key2=true&key3=value3",
		},
		{
			input: map[string]any{
				"key3": "",
				"key4": "",
				"key1": "123",
				"key2": "value2",
			},
			expected: "key1=123&key2=value2",
		},
	}
	for _, m := range mocks {
		result := help.MapToSignText(m.input)
		assert.Equal(t, m.expected, result)
	}
}

func TestMapToSignText_MoreTypes(t *testing.T) {
	m := map[string]any{
		"a": float32(1.5),
		"b": nil,
		"c": []int{1, 2},
		"d": struct{ X int }{X: 2},
		"e": uint(3),
	}
	r := help.MapToSignText(m)
	assert.Contains(t, r, "a=1.5")
	assert.Contains(t, r, "c=[1 2]")
	assert.Contains(t, r, "d={2}")
	assert.Contains(t, r, "e=3")
}

func TestPtr(t *testing.T) {
	assert.Equal(t, "hello", *help.Ptr[string]("hello"))
	assert.Equal(t, int64(123), *help.Ptr[int64](123))
	assert.Equal(t, false, *help.Ptr[bool](false))
}

func TestIsEmpty(t *testing.T) {
	// string
	assert.True(t, help.IsEmpty(""))
	assert.False(t, help.IsEmpty("a"))
	// array
	var arr0 [0]int
	var arr1 = [1]int{1}
	assert.True(t, help.IsEmpty(arr0))
	assert.False(t, help.IsEmpty(arr1))
	// map
	var mNil map[string]int
	m0 := map[string]int{}
	m1 := map[string]int{"a": 1}
	assert.True(t, help.IsEmpty(mNil))
	assert.True(t, help.IsEmpty(m0))
	assert.False(t, help.IsEmpty(m1))
	// slice
	var sNil []int
	s0 := make([]int, 0)
	s1 := []int{1}
	assert.True(t, help.IsEmpty(sNil))
	assert.True(t, help.IsEmpty(s0))
	assert.False(t, help.IsEmpty(s1))
	// bool
	assert.True(t, help.IsEmpty(false))
	assert.False(t, help.IsEmpty(true))
	// int/uint
	assert.True(t, help.IsEmpty(int(0)))
	assert.False(t, help.IsEmpty(int(1)))
	assert.True(t, help.IsEmpty(uint(0)))
	assert.False(t, help.IsEmpty(uint(1)))
	// float
	assert.True(t, help.IsEmpty(float32(0)))
	assert.False(t, help.IsEmpty(float32(1)))
	assert.True(t, help.IsEmpty(float64(0)))
	assert.False(t, help.IsEmpty(float64(1)))
	// interface/ptr
	var i interface{}
	assert.True(t, help.IsEmpty(i))
	p := help.Ptr(0)
	assert.False(t, help.IsEmpty(p))
	// func
	var fn func()
	assert.True(t, help.IsEmpty(fn))
	fn = func() {}
	assert.False(t, help.IsEmpty(fn))

	// Additional tests
	assert.True(t, help.IsEmpty(nil))
	assert.True(t, help.IsEmpty(0))
	assert.True(t, help.IsEmpty(false))
	assert.False(t, help.IsEmpty(help.Ptr[int64](0)))
	var a *string
	assert.True(t, help.IsEmpty(a))
	var b struct{}
	assert.True(t, help.IsEmpty(b))
}

func TestSha256hex(t *testing.T) {
	expected := help.Sha256hex("hello")
	assert.Len(t, expected, 64)
	assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", expected)
}

func TestHmacSha256(t *testing.T) {
	expected := help.HmacSha256("hello", "key")
	assert.Len(t, expected, 32)
}

func TestOkFail(t *testing.T) {
	r1 := help.Ok()
	assert.Equal(t, int64(0), r1.Code)
	assert.Equal(t, "ok", r1.Message)

	r2 := help.Fail(1234, "oops")
	assert.Equal(t, int64(1234), r2.Code)
	assert.Equal(t, "oops", r2.Message)
}

func TestE(t *testing.T) {
	e := help.E(2002, "bad")
	assert.True(t, e.IsType(help.ErrorTypePublic))
	meta, ok := e.Meta.(*help.ErrorMeta)
	assert.True(t, ok)
	assert.Equal(t, int64(2002), meta.Code)
}

func TestValidator_CustomRules(t *testing.T) {
	v := help.Validator()
	type S struct {
		Name string `vd:"snake"`
		Sort string `vd:"sort"`
	}
	ok := S{Name: "user_name", Sort: "created_at:1"}
	err := v.Validate(&ok)
	assert.NoError(t, err)

	bad1 := S{Name: "UserName", Sort: "created_at:1"}
	err = v.Validate(&bad1)
	assert.Error(t, err)

	bad2 := S{Name: "user_name", Sort: "created_at:2"}
	err = v.Validate(&bad2)
	assert.Error(t, err)
}

func TestErrorHandler_Public(t *testing.T) {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/public", help.ErrorHandler(), func(ctx context.Context, c *app.RequestContext) {
		c.Error(help.E(1001, "bad request"))
	})
	w := ut.PerformRequest(router, "GET", "/public", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, 400, resp.StatusCode())
}

func TestErrorHandler_Validation(t *testing.T) {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/validation", help.ErrorHandler(), func(ctx context.Context, c *app.RequestContext) {
		vd := validator.New()
		var s struct {
			Email string `validate:"email"`
		}
		s.Email = "not-email"
		err := vd.Struct(s)
		c.Error(err)
	})
	w := ut.PerformRequest(router, "GET", "/validation", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, 400, resp.StatusCode())
}

func TestErrorHandler_Internal(t *testing.T) {
	os.Setenv("MODE", "dev")
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/internal", help.ErrorHandler(), func(ctx context.Context, c *app.RequestContext) {
		c.Error(errors.New("something wrong"))
	})
	w := ut.PerformRequest(router, "GET", "/internal", &ut.Body{bytes.NewBuffer(nil), 0})
	resp := w.Result()
	assert.Equal(t, 500, resp.StatusCode())
}

func TestUuid(t *testing.T) {
	v := help.Uuid()
	id, err := uuid.Parse(v)
	assert.NoError(t, err)
	assert.Equal(t, uuid.Version(4), id.Version())
}

func TestUuid7(t *testing.T) {
	v := help.Uuid7()
	id, err := uuid.Parse(v)
	assert.NoError(t, err)
	assert.Equal(t, uuid.Version(7), id.Version())

	// Test time ordering
	v1 := help.Uuid7()
	v2 := help.Uuid7()
	assert.True(t, v1 < v2, "UUIDv7 should be time-ordered")
}

func TestMustUuid7(t *testing.T) {
	v := help.MustUuid7()
	id, err := uuid.Parse(v)
	assert.NoError(t, err)
	assert.Equal(t, uuid.Version(7), id.Version())
}

func TestUuid7Time(t *testing.T) {
	v := help.Uuid7()

	ts, ok := help.Uuid7Time(v)
	assert.True(t, ok)
	assert.Greater(t, ts, int64(0))

	// Timestamp should be close to now (within 1 second)
	now := time.Now().UnixMilli()
	assert.InDelta(t, now, ts, 1000)

	// Test invalid UUID
	_, ok = help.Uuid7Time("invalid")
	assert.False(t, ok)

	// Test UUIDv4 (not v7)
	v4 := help.Uuid()
	_, ok = help.Uuid7Time(v4)
	assert.False(t, ok)
}

func TestSID(t *testing.T) {
	v1 := help.SID()
	v2 := help.SID()
	assert.NotEmpty(t, v1)
	assert.NotEmpty(t, v2)
	assert.NotEqual(t, v1, v2)
}

func TestSIDWithError(t *testing.T) {
	// Test normal operation
	id, err := help.SIDWithError()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Test with custom Sonyflake that returns an error
	// sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	// id, err = sf.NextID()
	// Use reflection or create a mock to test error condition
	// For now, just ensure the function exists and works
	// _, _ = id, err
}

func TestSm2(t *testing.T) {
	priKeyStr := `MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQg1QG/R5oI4OO3mSh7Nss5RP7d8rV571CCyW+7cI1+w5qgCgYIKoEcz1UBgi2hRANCAAShrB20h+g1nL++oRUMpCsqdAb+ALVoUSpnR4jencQj3arGNQJA9rSdmvh6k64eI6gLZNxxk2YXXm5A70a/s1iz`

	priKey, err := help.PrivKeySM2FromBase64(priKeyStr)
	assert.NoError(t, err)

	sig, err := help.Sm2Sign(priKey, `Hello world!`)
	assert.NoError(t, err)
	t.Log(sig)

	pub := priKey.PublicKey
	ecdsaPub := &ecdsa.PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}

	r, err := help.Sm2Verify(ecdsaPub, `Hello world`, sig)
	assert.NoError(t, err)
	assert.False(t, r)

	r, err = help.Sm2Verify(ecdsaPub, `Hello world!`, sig)
	assert.NoError(t, err)
	assert.True(t, r)
}

func TestSm2PublicKey(t *testing.T) {
	priKey, err := sm2.GenerateKey(rand.Reader)
	assert.NoError(t, err)
	b, err := smx509.MarshalPKIXPublicKey(&priKey.PublicKey)
	assert.NoError(t, err)

	pubKeyStr := base64.StdEncoding.EncodeToString(b)
	t.Log(pubKeyStr)
	pubKey, err := help.PubKeySM2FromBase64(pubKeyStr)
	assert.NoError(t, err)

	t.Log(pubKey)
	t.Log(sm2.IsSM2PublicKey(pubKey))
}

func TestSm2ParseAndVerify_More(t *testing.T) {
	// generate key pair
	priKey, err := sm2.GenerateKey(rand.Reader)
	assert.NoError(t, err)
	b, err := smx509.MarshalPKIXPublicKey(&priKey.PublicKey)
	assert.NoError(t, err)
	pubKeyStr := base64.StdEncoding.EncodeToString(b)
	pubKey, err := help.PubKeySM2FromBase64(pubKeyStr)
	assert.NoError(t, err)

	// convert to ecdsa.PublicKey
	ecdsaPub := &ecdsa.PublicKey{Curve: pubKey.Curve, X: pubKey.X, Y: pubKey.Y}
	sig, err := help.Sm2Sign(priKey, "Hello")
	assert.NoError(t, err)
	ok, err := help.Sm2Verify(ecdsaPub, "Hello", sig)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestSM4(t *testing.T) {
	key := `f93c920868b4e5a88dfb27fd44b9f8db`
	plaintext := "hello world"

	t.Log(`明文`, plaintext)
	ciphertext, err := help.SM4Encrypt(key, plaintext)
	assert.NoError(t, err)
	t.Log(`密文`, ciphertext)

	// 解密
	decryptedText, err := help.SM4Decrypt(key, ciphertext)
	assert.NoError(t, err)
	t.Log(`解密结果`, decryptedText)

	// 验证
	valid, err := help.SM4Verify(key, ciphertext, plaintext)
	assert.NoError(t, err)
	if !valid {
		t.Fatal(`验证失败，解密结果与原文不一致`)
	}
	t.Log(`验证成功`)

	// 使用示例密钥进行测试
	testCiphertext := "056df5b3d1b15e2567d0dcd6e6cfbeff"
	testDecryptedText, err := help.SM4Decrypt(key, testCiphertext)
	assert.NoError(t, err)
	assert.Equal(t, testDecryptedText, plaintext)
	t.Log(`测试解密结果`, testDecryptedText)
}

func TestSM4_InvalidInputs(t *testing.T) {
	_, err := help.SM4Encrypt("zzz", "hello")
	assert.Error(t, err)
	// invalid hex ciphertext for decrypt
	_, err = help.SM4Decrypt("f93c920868b4e5a88dfb27fd44b9f8db", "not-hex")
	assert.Error(t, err)
}
