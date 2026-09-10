package help

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	errx "errors"
	"os"
	"reflect"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/requestid"
	"github.com/kainonly/go/vd"
)

// Ptr 返回指向给定值的指针。
// 适合用于创建字面量的指针。
//
// Deprecated: Go 1.26 起 new 内置函数支持传值表达式
// （如 new(42)），请改用它。
func Ptr[T any](i T) *T {
	return &i
}

// IsEmpty 判断一个值是否被视为空值。
// 对 nil、空字符串、零值、空切片/map 等返回 true。
func IsEmpty(i any) bool {
	if i == nil || i == "" {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

// Sha256hex 计算 SHA256 哈希并返回十六进制编码的字符串。
func Sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

// HmacSha256 计算 HMAC-SHA256 并以字符串形式返回原始字节。
// 若需十六进制输出，请对结果使用 hex.EncodeToString。
func HmacSha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

// Validator 创建面向 Hertz 的 go-playground validator 配置。
// 它将 "vd" 设为验证标签，并注册自定义验证器：
//   - snake：验证 snake_case 格式（如 "user_name"）
//   - sort：验证排序格式（如 "created_at:1" 或 "name:-1"）
//
// Deprecated: 请改用 vd.Default()。
func Validator() *vd.Validator {
	return vd.Default()
}

// R 是标准的 API 响应结构。
type R struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Ok 返回 code 为 0、message 为 "ok" 的成功响应。
func Ok() R {
	return R{
		Code:    0,
		Message: "ok",
	}
}

// Fail 返回携带给定 code 和 message 的错误响应。
func Fail(code int64, msg string) R {
	return R{
		Code:    code,
		Message: msg,
	}
}

// ErrorMeta 包含 Hertz 错误的元数据。
type ErrorMeta struct {
	Code int64
}

// E 创建携带 code 的公开 Hertz 错误。
// 用于需要展示给用户的业务逻辑错误。
func E(code int64, msg string) *errors.Error {
	return errors.NewPublic(msg).SetMeta(&ErrorMeta{Code: code})
}

// ErrorTypePublic 是 Hertz 框架中公开错误的类型。
var ErrorTypePublic = errors.ErrorTypePublic

// ErrorHandler 返回处理错误的 Hertz 中间件。
// 它处理不同的错误类型：
//   - 公开错误：返回 400 及 code 和 message
//   - 验证错误：返回 400 及字段详情
//   - 其他错误：返回 500（开发模式下附带详情）
//
// 生产模式请设置 MODE=release 环境变量。
func ErrorHandler() app.HandlerFunc {
	release := os.Getenv("MODE") == "release"
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		e := c.Errors.Last()
		if e == nil {
			return
		}

		if e.IsType(errors.ErrorTypePublic) {
			r := R{Code: 0, Message: e.Error()}
			if meta, ok := e.Meta.(*ErrorMeta); ok {
				r.Code = meta.Code
			}
			c.JSON(400, r)
			return
		}

		var ves validator.ValidationErrors
		if errx.As(e.Err, &ves) {
			message := make([]interface{}, len(ves))
			for i, v := range ves {
				message[i] = utils.H{
					"namespace": v.Namespace(),
					"field":     v.Field(),
					"tag":       v.Tag(),
				}
			}
			c.JSON(400, utils.H{
				"code":    0,
				"message": message,
			})
			return
		}

		if !release {
			c.JSON(500, e.JSON())
			return
		}

		logger.Error(requestid.Get(c), e)
		c.Status(500)
	}
}
