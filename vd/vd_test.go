package vd_test

import (
	"testing"

	"github.com/kainonly/go/vd"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	v := vd.New()
	assert.NotNil(t, v)
	assert.NotNil(t, v.Engine())
	assert.NotNil(t, v.Core())
}

func TestNewWithTag(t *testing.T) {
	v := vd.New(vd.SetTag("validate"))
	assert.NotNil(t, v)
}

func TestDefault(t *testing.T) {
	v := vd.Default()
	assert.NotNil(t, v)
}

func TestSnakeValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Snake()))

	type TestStruct struct {
		Field string `vd:"snake"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid snake_case", "user_name", false},
		{"valid single word", "name", false},
		{"valid with underscore", "created_at", false},
		{"valid multiple underscores", "user_first_name", false},
		{"invalid uppercase", "User_Name", true},
		{"invalid with number", "user1", true},
		{"invalid with hyphen", "user-name", true},
		{"invalid empty", "", true},
		{"invalid space", "user name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(TestStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSortValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Sort()))

	type TestStruct struct {
		Field string `vd:"sort"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid ascending", "created_at:1", false},
		{"valid descending", "name:-1", false},
		{"valid with underscore", "user_name:1", false},
		{"invalid without direction", "name", true},
		{"invalid direction 0", "name:0", true},
		{"invalid direction 2", "name:2", true},
		{"invalid uppercase", "Name:1", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(TestStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPhoneValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Phone()))

	type TestStruct struct {
		Field string `vd:"phone"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid phone 138", "13800138000", false},
		{"valid phone 199", "19912345678", false},
		{"valid phone 177", "17712345678", false},
		{"invalid too short", "1380013800", true},
		{"invalid too long", "138001380001", true},
		{"invalid prefix 10", "10800138000", true},
		{"invalid prefix 12", "12800138000", true},
		{"invalid with letters", "1380013800a", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(TestStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIDCardValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.IDCard()))

	type TestStruct struct {
		Field string `vd:"idcard"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid id card", "110101199003071234", false},
		{"valid with X", "11010119900307123X", false},
		{"valid with x", "11010119900307123x", false},
		{"invalid too short", "11010119900307123", true},
		{"invalid too long", "1101011990030712345", true},
		{"invalid year", "110101189003071234", true},
		{"invalid month", "110101199013071234", true},
		{"invalid day", "110101199003321234", true},
		{"invalid prefix 0", "010101199003071234", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(TestStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsernameValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Username()))

	type TestStruct struct {
		Field string `vd:"username"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid simple", "john", false},
		{"valid with underscore", "john_doe", false},
		{"valid with numbers", "john123", false},
		{"valid mixed", "John_Doe123", false},
		{"valid min length", "abc", false},
		{"valid max length", "abcdefghijklmnopqrst", false},
		{"invalid too short", "ab", true},
		{"invalid too long", "abcdefghijklmnopqrstu", true},
		{"invalid start with number", "1john", true},
		{"invalid start with underscore", "_john", true},
		{"invalid with hyphen", "john-doe", true},
		{"invalid with space", "john doe", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(TestStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSlugValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Slug()))

	type TestStruct struct {
		Field string `vd:"slug"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid simple", "hello", false},
		{"valid with hyphen", "hello-world", false},
		{"valid with numbers", "post123", false},
		{"valid mixed", "my-blog-post-123", false},
		{"invalid uppercase", "Hello", true},
		{"invalid underscore", "hello_world", true},
		{"invalid start hyphen", "-hello", true},
		{"invalid end hyphen", "hello-", true},
		{"invalid double hyphen", "hello--world", true},
		{"invalid space", "hello world", true},
		{"invalid empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(TestStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateVar(t *testing.T) {
	v := vd.New()

	err := v.ValidateVar("test@example.com", "email")
	assert.NoError(t, err)

	err = v.ValidateVar("invalid-email", "email")
	assert.Error(t, err)
}

func TestRegisterRule(t *testing.T) {
	v := vd.New()

	err := v.RegisterRule(vd.Rule{
		Tag: "even",
		Fn: func(fl vd.FieldLevel) bool {
			num, ok := fl.Field().Interface().(int)
			if !ok {
				return false
			}
			return num%2 == 0
		},
	})
	assert.NoError(t, err)

	type TestStruct struct {
		Num int `vd:"even"`
	}

	err = v.Validate(TestStruct{Num: 4})
	assert.NoError(t, err)

	err = v.Validate(TestStruct{Num: 3})
	assert.Error(t, err)
}

func TestMultipleRules(t *testing.T) {
	v := vd.New(vd.SetRules(
		vd.Snake(),
		vd.Sort(),
		vd.Phone(),
	))

	type User struct {
		Name  string `vd:"snake"`
		Sort  string `vd:"sort"`
		Phone string `vd:"phone"`
	}

	err := v.Validate(User{
		Name:  "user_name",
		Sort:  "created_at:-1",
		Phone: "13800138000",
	})
	assert.NoError(t, err)
}

func TestCombinedWithBuiltInRules(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Snake()))

	type User struct {
		Name  string `vd:"snake,min=3,max=20"`
		Email string `vd:"required,email"`
		Age   int    `vd:"gte=0,lte=150"`
	}

	// Valid case
	err := v.Validate(User{
		Name:  "user_name",
		Email: "test@example.com",
		Age:   25,
	})
	assert.NoError(t, err)

	// Invalid snake case
	err = v.Validate(User{
		Name:  "UserName",
		Email: "test@example.com",
		Age:   25,
	})
	assert.Error(t, err)

	// Invalid email
	err = v.Validate(User{
		Name:  "user_name",
		Email: "invalid",
		Age:   25,
	})
	assert.Error(t, err)

	// Invalid age
	err = v.Validate(User{
		Name:  "user_name",
		Email: "test@example.com",
		Age:   -1,
	})
	assert.Error(t, err)
}
