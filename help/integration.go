package help

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	errx "errors"
	"os"
	"reflect"
	"regexp"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/binding/go_playground"
	"github.com/hertz-contrib/requestid"
)

// Ptr returns a pointer to the given value.
// Useful for creating pointers to literals.
func Ptr[T any](i T) *T {
	return &i
}

// IsEmpty checks if a value is considered empty.
// Returns true for nil, empty strings, zero values, empty slices/maps, etc.
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

// Sha256hex computes SHA256 hash and returns hex-encoded string.
func Sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

// HmacSha256 computes HMAC-SHA256 and returns raw bytes as string.
// For hex output, use hex.EncodeToString on the result.
func HmacSha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

// Validator creates a configured go-playground validator for Hertz.
// It sets "vd" as the validation tag and registers custom validators:
//   - snake: validates snake_case format (e.g., "user_name")
//   - sort: validates sort format (e.g., "created_at:1" or "name:-1")
func Validator() *go_playground.Validator {
	vd := go_playground.NewValidator()
	vd.SetValidateTag("vd")
	vdx := vd.Engine().(*validator.Validate)
	vdx.RegisterValidation("snake", func(fl validator.FieldLevel) bool {
		matched, err := regexp.MatchString("^[a-z_]+$", fl.Field().Interface().(string))
		if err != nil {
			return false
		}
		return matched
	})
	vdx.RegisterValidation("sort", func(fl validator.FieldLevel) bool {
		matched, err := regexp.MatchString("^[a-z_]+:(-1|1)$", fl.Field().Interface().(string))
		if err != nil {
			return false
		}
		return matched
	})
	return vd
}

// R is a standard API response structure.
type R struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Ok returns a success response with code 0 and message "ok".
func Ok() R {
	return R{
		Code:    0,
		Message: "ok",
	}
}

// Fail returns an error response with the given code and message.
func Fail(code int64, msg string) R {
	return R{
		Code:    code,
		Message: msg,
	}
}

// ErrorMeta contains error metadata for Hertz errors.
type ErrorMeta struct {
	Code int64
}

// E creates a public Hertz error with a code.
// Use this for business logic errors that should be shown to users.
func E(code int64, msg string) *errors.Error {
	return errors.NewPublic(msg).SetMeta(&ErrorMeta{Code: code})
}

// ErrorTypePublic is the type for public errors in Hertz framework.
var ErrorTypePublic = errors.ErrorTypePublic

// ErrorHandler returns a Hertz middleware that handles errors.
// It processes different error types:
//   - Public errors: Returns 400 with code and message
//   - Validation errors: Returns 400 with field details
//   - Other errors: Returns 500 (with details in dev mode)
//
// Set MODE=release environment variable for production mode.
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
