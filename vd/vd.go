// Package vd provides a configurable validator wrapper for go-playground/validator
// with Hertz framework integration and custom validation rules support.
package vd

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/hertz-contrib/binding/go_playground"
)

// ValidationFunc is a custom validation function type.
type ValidationFunc = validator.Func

// FieldLevel contains all the information and helper functions
// to validate a field.
type FieldLevel = validator.FieldLevel

// Rule defines a custom validation rule.
type Rule struct {
	Tag  string
	Fn   ValidationFunc
	Call bool // CallValidationEvenIfNull
}

// Option configures the Validator.
type Option func(*options)

type options struct {
	tag   string
	rules []Rule
}

// SetTag sets the validation struct tag name.
// Default is "vd".
func SetTag(tag string) Option {
	return func(o *options) {
		o.tag = tag
	}
}

// SetRules sets the custom validation rules.
func SetRules(rules ...Rule) Option {
	return func(o *options) {
		o.rules = append(o.rules, rules...)
	}
}

// Validator wraps go-playground validator with custom configuration.
type Validator struct {
	engine *go_playground.Validator
	core   *validator.Validate
}

// New creates a new Validator with the given options.
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

	// Register custom rules
	for _, rule := range o.rules {
		if rule.Call {
			core.RegisterValidation(rule.Tag, rule.Fn, true)
		} else {
			core.RegisterValidation(rule.Tag, rule.Fn)
		}
	}

	return &Validator{
		engine: vd,
		core:   core,
	}
}

// Engine returns the underlying go-playground validator for Hertz.
func (v *Validator) Engine() *go_playground.Validator {
	return v.engine
}

// Core returns the underlying go-playground/validator/v10 instance.
func (v *Validator) Core() *validator.Validate {
	return v.core
}

// Validate validates a struct.
func (v *Validator) Validate(obj any) error {
	return v.core.Struct(obj)
}

// ValidateVar validates a single variable using tag style validation.
func (v *Validator) ValidateVar(field any, tag string) error {
	return v.core.Var(field, tag)
}

// RegisterRule registers a custom validation rule dynamically.
func (v *Validator) RegisterRule(rule Rule) error {
	if rule.Call {
		return v.core.RegisterValidation(rule.Tag, rule.Fn, true)
	}
	return v.core.RegisterValidation(rule.Tag, rule.Fn)
}

// Preset validation rules

// Snake validates snake_case format (lowercase letters and underscores only).
// Example: "user_name", "created_at"
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

// Sort validates sort format: "field_name:1" or "field_name:-1"
// Example: "created_at:1", "name:-1"
func Sort() Rule {
	return Rule{
		Tag: "sort",
		Fn:  sortValidation,
	}
}

var sortRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z_]+:(-1|1)$`)
})

func sortValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return sortRegex().MatchString(s)
}

// Phone validates Chinese mobile phone number format.
// Example: "13800138000"
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

// IDCard validates Chinese ID card number (18 digits).
// Example: "110101199003077758"
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

// Username validates username format (alphanumeric, underscore, 3-20 chars).
// Example: "john_doe", "user123"
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

// Slug validates URL slug format (lowercase, numbers, hyphens).
// Example: "my-blog-post", "article-123"
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

// Default creates a Validator with preset rules (snake, sort).
func Default() *Validator {
	return New(
		SetRules(Snake(), Sort()),
	)
}
