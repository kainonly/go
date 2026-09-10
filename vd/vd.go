// Package vd 提供对 go-playground/validator 的可配置封装，
// 支持 Hertz 框架集成与自定义验证规则。
//
// # Hertz 集成
//
// 将验证器引擎传入 server 选项，即可在 Hertz 服务端中使用：
//
//	import (
//	    "github.com/cloudwego/hertz/pkg/app/server"
//	    "github.com/kainonly/go/vd"
//	)
//
//	func main() {
//	    // Create validator with default rules (snake, sort)
//	    v := vd.Default()
//
//	    // Or create with custom rules
//	    v := vd.New(
//	        vd.SetTag("vd"),  // optional, "vd" is default
//	        vd.SetRules(
//	            vd.Snake(),
//	            vd.Sort(),
//	            vd.Phone(),
//	            vd.IDCard(),
//	            vd.PasswordMedium(),
//	        ),
//	    )
//
//	    // Pass to Hertz server
//	    h := server.Default(
//	        server.WithHostPorts(":8080"),
//	        server.WithCustomValidator(v.Engine()),
//	    )
//
//	    h.POST("/user", func(ctx context.Context, c *app.RequestContext) {
//	        var req struct {
//	            Name  string `json:"name" vd:"required,snake"`
//	            Phone string `json:"phone" vd:"required,phone"`
//	        }
//	        if err := c.BindAndValidate(&req); err != nil {
//	            c.JSON(400, map[string]any{"error": err.Error()})
//	            return
//	        }
//	        c.JSON(200, map[string]any{"message": "ok"})
//	    })
//
//	    h.Spin()
//	}
//
// # 可用规则组
//
//   - All()：所有可用规则
//   - Common()：常用规则（snake、sort、phone、idcard、username、slug、password_medium）
//   - Chinese()：中国本地化规则（phone、idcard、bankcard、license_plate 等）
//   - NamingConvention()：代码命名规则（snake、pascal、camel、kebab、upper_snake、variable）
//
// # 自定义规则
//
// 注册自定义验证规则：
//
//	v := vd.New(vd.SetRules(
//	    vd.Rule{
//	        Tag: "even",
//	        Fn: func(fl vd.FieldLevel) bool {
//	            return fl.Field().Int() % 2 == 0
//	        },
//	    },
//	))
package vd

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/binding/go_playground"
)

// ValidationFunc 是自定义验证函数类型。
type ValidationFunc = validator.Func

// FieldLevel 包含校验字段所需的全部信息和辅助函数。
type FieldLevel = validator.FieldLevel

// Rule 定义一条自定义验证规则。
type Rule struct {
	Tag  string
	Fn   ValidationFunc
	Call bool // CallValidationEvenIfNull
}

// Option 用于配置 Validator。
type Option func(*options)

type options struct {
	tag   string
	rules []Rule
}

// SetTag 设置验证使用的 struct tag 名称。
// 默认为 "vd"。
func SetTag(tag string) Option {
	return func(o *options) {
		o.tag = tag
	}
}

// SetRules 设置自定义验证规则。
func SetRules(rules ...Rule) Option {
	return func(o *options) {
		o.rules = append(o.rules, rules...)
	}
}

// Validator 封装 go-playground validator，支持自定义配置。
type Validator struct {
	engine *go_playground.Validator
	core   *validator.Validate
}

// New 根据给定选项创建新的 Validator。
func New(opts ...Option) *Validator {
	o := &options{
		tag: "vd",
	}
	for _, opt := range opts {
		opt(o)
	}

	vd := go_playground.NewValidator()
	vd.SetValidateTag(o.tag)
	core := vd.Engine().(*validator.Validate)

	// 注册自定义规则
	for _, rule := range o.rules {
		if rule.Call {
			_ = core.RegisterValidation(rule.Tag, rule.Fn, true)
		} else {
			_ = core.RegisterValidation(rule.Tag, rule.Fn)
		}
	}

	return &Validator{
		engine: vd,
		core:   core,
	}
}

// Engine 返回供 Hertz 使用的底层 go-playground validator。
func (v *Validator) Engine() *go_playground.Validator {
	return v.engine
}

// Core 返回底层的 go-playground/validator/v10 实例。
func (v *Validator) Core() *validator.Validate {
	return v.core
}

// Validate 校验一个 struct。
func (v *Validator) Validate(obj any) error {
	return v.core.Struct(obj)
}

// ValidateVar 使用 tag 风格校验单个变量。
func (v *Validator) ValidateVar(field any, tag string) error {
	return v.core.Var(field, tag)
}

// RegisterRule 动态注册一条自定义验证规则。
func (v *Validator) RegisterRule(rule Rule) error {
	if rule.Call {
		return v.core.RegisterValidation(rule.Tag, rule.Fn, true)
	}
	return v.core.RegisterValidation(rule.Tag, rule.Fn)
}

// 预设验证规则

// Snake 校验 snake_case 格式（仅小写字母和下划线）。
// 示例："user_name"、"created_at"
func Snake() Rule {
	return Rule{
		Tag: "snake",
		Fn:  snakeValidation,
	}
}

var snakeRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_]+$`)
})

func snakeValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return snakeRegex().MatchString(s)
}

// Sort 校验排序格式："field_name:1" 或 "field_name:-1"
// 示例："created_at:1"、"name:-1"
func Sort() Rule {
	return Rule{
		Tag: "sort",
		Fn:  sortValidation,
	}
}

var sortRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_]+(\.[a-z_]+)*:(-1|1)$`)
})

func sortValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return sortRegex().MatchString(s)
}

// Phone 校验中国手机号格式。
// 示例："13800138000"
func Phone() Rule {
	return Rule{
		Tag: "phone",
		Fn:  phoneValidation,
	}
}

var phoneRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^1[3-9]\d{9}$`)
})

func phoneValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return phoneRegex().MatchString(s)
}

// IDCard 校验中国身份证号（18 位）。
// 示例："110101199003077758"
func IDCard() Rule {
	return Rule{
		Tag: "idcard",
		Fn:  idcardValidation,
	}
}

var idcardRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)
})

func idcardValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return idcardRegex().MatchString(s)
}

// Username 校验用户名格式（字母、数字、下划线，3-20 个字符）。
// 示例："john_doe"、"user123"
func Username() Rule {
	return Rule{
		Tag: "username",
		Fn:  usernameValidation,
	}
}

var usernameRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,19}$`)
})

func usernameValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return usernameRegex().MatchString(s)
}

// Slug 校验 URL slug 格式（小写字母、数字、连字符）。
// 示例："my-blog-post"、"article-123"
func Slug() Rule {
	return Rule{
		Tag: "slug",
		Fn:  slugValidation,
	}
}

var slugRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
})

func slugValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return slugRegex().MatchString(s)
}

// Default 创建带预设规则（snake、sort）的 Validator。
func Default() *Validator {
	return New(
		SetRules(Snake(), Sort()),
	)
}
