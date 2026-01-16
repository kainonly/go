package vd_test

import (
	"testing"

	"github.com/kainonly/go/vd"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Chinese Localization Rules Tests
// ============================================================================

func TestBankCardValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.BankCard()))

	type TestStruct struct {
		Field string `vd:"bankcard"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 16 digits", "6222021234567890", true}, // random number, Luhn may fail
		{"valid luhn", "4532015112830366", false},     // valid Luhn
		{"invalid too short", "622202123456", true},
		{"invalid too long", "62220212345678901234", true},
		{"invalid with letters", "622202123456789a", true},
		{"invalid luhn", "1234567890123456", true},
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

func TestLicensePlateValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.LicensePlate()))

	type TestStruct struct {
		Field string `vd:"license_plate"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid beijing", "京A12345", false},
		{"valid shanghai", "沪B88888", false},
		{"valid new energy 6", "粤A123456", false},
		{"valid new energy D", "京AD12345", false},
		{"invalid too short", "京A1234", true},
		{"invalid too long", "京A1234567", true},
		{"invalid province", "XA12345", true},
		{"invalid lowercase", "京a12345", true},
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

func TestUSCCValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.USCC()))

	type TestStruct struct {
		Field string `vd:"uscc"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid uscc", "91310000MA1FL8TQ32", false},
		{"valid uscc 2", "91110108MA01KPGM09", false},
		{"invalid too short", "91310000MA1FL8TQ3", true},
		{"invalid too long", "91310000MA1FL8TQ321", true},
		{"invalid char I", "9I310000MA1FL8TQ32", true},
		{"invalid char O", "91O10000MA1FL8TQ32", true},
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

func TestChineseWordValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.ChineseWord()))

	type TestStruct struct {
		Field string `vd:"chinese"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid chinese", "中国", false},
		{"valid chinese long", "你好世界", false},
		{"invalid with english", "中国China", true},
		{"invalid with number", "中国123", true},
		{"invalid english only", "China", true},
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

func TestChineseNameValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.ChineseName()))

	type TestStruct struct {
		Field string `vd:"chinese_name"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 2 chars", "张三", false},
		{"valid 3 chars", "李小龙", false},
		{"valid minority", "古力娜扎·迪丽热巴", false},
		{"invalid too long", "张三李四王五赵六钱七", true},
		{"invalid with number", "张三3", true},
		{"invalid english", "Zhang", true},
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

// ============================================================================
// Password Strength Rules Tests
// ============================================================================

func TestPasswordWeakValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.PasswordWeak()))

	type TestStruct struct {
		Field string `vd:"password_weak"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 6 chars", "123456", false},
		{"valid 8 chars", "12345678", false},
		{"invalid 5 chars", "12345", true},
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

func TestPasswordMediumValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.PasswordMedium()))

	type TestStruct struct {
		Field string `vd:"password_medium"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid letters and numbers", "password123", false},
		{"valid mixed", "abc12345", false},
		{"invalid numbers only", "12345678", true},
		{"invalid letters only", "password", true},
		{"invalid too short", "pass1", true},
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

func TestPasswordStrongValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.PasswordStrong()))

	type TestStruct struct {
		Field string `vd:"password_strong"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid strong", "Password123!", false},
		{"valid strong 2", "MyP@ssw0rd", false},
		{"invalid no special", "Password123", true},
		{"invalid no uppercase", "password123!", true},
		{"invalid no lowercase", "PASSWORD123!", true},
		{"invalid no number", "Password!!!", true},
		{"invalid too short", "Pass1!", true},
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

// ============================================================================
// Special Format Rules Tests
// ============================================================================

func TestObjectIDValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.ObjectID()))

	type TestStruct struct {
		Field string `vd:"objectid"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid objectid", "507f1f77bcf86cd799439011", false},
		{"valid uppercase", "507F1F77BCF86CD799439011", false},
		{"invalid too short", "507f1f77bcf86cd79943901", true},
		{"invalid too long", "507f1f77bcf86cd7994390111", true},
		{"invalid chars", "507f1f77bcf86cd79943901g", true},
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

func TestSnowflakeValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Snowflake()))

	type StringStruct struct {
		Field string `vd:"snowflake"`
	}

	type Int64Struct struct {
		Field int64 `vd:"snowflake"`
	}

	// String tests
	stringTests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid snowflake", "1234567890123456789", false},
		{"valid 15 digits", "123456789012345", false},
		{"invalid too short", "12345678901234", true},
		{"invalid too long", "12345678901234567890", true},
		{"invalid with letters", "123456789012345a", true},
	}

	for _, tt := range stringTests {
		t.Run("string_"+tt.name, func(t *testing.T) {
			err := v.Validate(StringStruct{Field: tt.value})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Int64 tests
	t.Run("int64_valid", func(t *testing.T) {
		err := v.Validate(Int64Struct{Field: 1234567890123456789})
		assert.NoError(t, err)
	})

	t.Run("int64_invalid_zero", func(t *testing.T) {
		err := v.Validate(Int64Struct{Field: 0})
		assert.Error(t, err)
	})
}

func TestVersionValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Version()))

	type TestStruct struct {
		Field string `vd:"version"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid version", "1.0.0", false},
		{"valid version 2", "12.34.56", false},
		{"invalid no patch", "1.0", true},
		{"invalid with v", "v1.0.0", true},
		{"invalid with suffix", "1.0.0-beta", true},
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

// ============================================================================
// Data Structure Rules Tests
// ============================================================================

func TestSafeStringValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.SafeString()))

	type TestStruct struct {
		Field string `vd:"safe_string"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid safe", "Hello World 123", false},
		{"valid with dash", "my-text_here", false},
		{"invalid script tag", "<script>", true},
		{"invalid quotes", "test\"injection", true},
		{"invalid semicolon", "cmd;rm", true},
		{"invalid pipe", "cat|grep", true},
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

func TestDecimalValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Decimal()))

	type TestStruct struct {
		Field string `vd:"decimal"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid integer", "123", false},
		{"valid decimal", "123.45", false},
		{"valid negative", "-123.45", false},
		{"valid zero", "0", false},
		{"invalid letters", "12.3a", true},
		{"invalid double dot", "12..3", true},
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

// ============================================================================
// Common Format Rules Tests
// ============================================================================

func TestDomainValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Domain()))

	type TestStruct struct {
		Field string `vd:"domain"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid domain", "example.com", false},
		{"valid subdomain", "sub.example.com", false},
		{"valid co.uk", "example.co.uk", false},
		{"invalid with protocol", "http://example.com", true},
		{"invalid no tld", "example", true},
		{"invalid with path", "example.com/path", true},
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

func TestFileNameValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.FileName()))

	type TestStruct struct {
		Field string `vd:"filename"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid filename", "document.pdf", false},
		{"valid with spaces", "my document.pdf", false},
		{"valid no extension", "README", false},
		{"invalid with slash", "path/file.txt", true},
		{"invalid with backslash", "path\\file.txt", true},
		{"invalid dot only", ".", true},
		{"invalid double dot", "..", true},
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

func TestColorValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Color()))

	type TestStruct struct {
		Field string `vd:"color"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 3 hex", "fff", false},
		{"valid 6 hex", "ff5733", false},
		{"valid 8 hex", "ff573380", false},
		{"valid uppercase", "FF5733", false},
		{"invalid with hash", "#ff5733", true},
		{"invalid 4 hex", "ff57", true},
		{"invalid chars", "gggggg", true},
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

// ============================================================================
// Communication Rules Tests
// ============================================================================

func TestTelPhoneValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.TelPhone()))

	type TestStruct struct {
		Field string `vd:"tel"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid beijing", "010-12345678", false},
		{"valid shenzhen", "0755-1234567", false},
		{"valid no dash", "02112345678", false},
		{"invalid no area", "12345678", true},
		{"invalid mobile", "13800138000", true},
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

func TestQQValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.QQ()))

	type TestStruct struct {
		Field string `vd:"qq"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 5 digits", "12345", false},
		{"valid 10 digits", "1234567890", false},
		{"invalid start 0", "01234567", true},
		{"invalid too short", "1234", true},
		{"invalid too long", "123456789012", true},
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

func TestWeChatValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.WeChat()))

	type TestStruct struct {
		Field string `vd:"wechat"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid wechat", "wxid_abc123", false},
		{"valid simple", "myWeChat01", false},
		{"invalid start number", "1wechat", true},
		{"invalid too short", "wx123", true},
		{"invalid special char", "wx@chat", true},
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

// ============================================================================
// Code Naming Rules Tests
// ============================================================================

func TestPascalCaseValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.PascalCase()))

	type TestStruct struct {
		Field string `vd:"pascal"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid pascal", "MyClass", false},
		{"valid single", "User", false},
		{"valid with number", "User123", false},
		{"invalid lowercase start", "myClass", true},
		{"invalid underscore", "My_Class", true},
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

func TestCamelCaseValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.CamelCase()))

	type TestStruct struct {
		Field string `vd:"camel"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid camel", "myVariable", false},
		{"valid single", "name", false},
		{"valid with number", "user123", false},
		{"invalid uppercase start", "MyVariable", true},
		{"invalid underscore", "my_variable", true},
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

func TestKebabCaseValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.KebabCase()))

	type TestStruct struct {
		Field string `vd:"kebab"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid kebab", "my-component", false},
		{"valid single", "button", false},
		{"valid with number", "item-123", false},
		{"invalid uppercase", "My-Component", true},
		{"invalid underscore", "my_component", true},
		{"invalid double dash", "my--component", true},
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

func TestUpperSnakeValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.UpperSnake()))

	type TestStruct struct {
		Field string `vd:"upper_snake"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid upper snake", "MAX_VALUE", false},
		{"valid single", "VALUE", false},
		{"valid with number", "HTTP_200", false},
		{"invalid lowercase", "max_value", true},
		{"invalid mixed", "Max_Value", true},
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

// ============================================================================
// Utility Functions Tests
// ============================================================================

func TestAllRules(t *testing.T) {
	rules := vd.All()
	assert.Greater(t, len(rules), 30)

	// Verify all rules can be registered
	v := vd.New(vd.SetRules(rules...))
	assert.NotNil(t, v)
}

func TestCommonRules(t *testing.T) {
	rules := vd.Common()
	assert.Equal(t, 7, len(rules))
}

func TestChineseRules(t *testing.T) {
	rules := vd.Chinese()
	assert.Equal(t, 11, len(rules))
}

func TestNamingConventionRules(t *testing.T) {
	rules := vd.NamingConvention()
	assert.Equal(t, 6, len(rules))
}

func TestNotBlankValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.NotBlank()))

	type TestStruct struct {
		Field string `vd:"notblank"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid text", "hello", false},
		{"valid with spaces", "  hello  ", false},
		{"invalid empty", "", true},
		{"invalid spaces only", "   ", true},
		{"invalid tabs only", "\t\t", true},
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

func TestZipCodeValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.ZipCode()))

	type TestStruct struct {
		Field string `vd:"zipcode"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid beijing", "100000", false},
		{"valid shenzhen", "518000", false},
		{"invalid too short", "10000", true},
		{"invalid too long", "1000000", true},
		{"invalid with letters", "10000a", true},
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

func TestVariableValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.Variable()))

	type TestStruct struct {
		Field string `vd:"variable"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid var", "myVar", false},
		{"valid underscore start", "_private", false},
		{"valid with number", "var123", false},
		{"valid uppercase", "MAX_VALUE", false},
		{"invalid start number", "123var", true},
		{"invalid with dash", "my-var", true},
		{"invalid with space", "my var", true},
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

func TestAlphaNumDashValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.AlphaNumDash()))

	type TestStruct struct {
		Field string `vd:"alphanumdash"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid simple", "abc123", false},
		{"valid with dash", "my-item", false},
		{"valid with underscore", "my_item", false},
		{"valid mixed", "my-item_123", false},
		{"invalid with space", "my item", true},
		{"invalid with dot", "my.item", true},
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

func TestPositiveDecimalValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.PositiveDecimal()))

	type TestStruct struct {
		Field string `vd:"positive_decimal"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid integer", "123", false},
		{"valid decimal", "123.45", false},
		{"valid small", "0.001", false},
		{"invalid zero", "0", true},
		{"invalid zero decimal", "0.0", true},
		{"invalid negative", "-123", true},
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

func TestFileExtValidation(t *testing.T) {
	v := vd.New(vd.SetRules(vd.FileExt()))

	type TestStruct struct {
		Field string `vd:"file_ext"`
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid with dot", ".pdf", false},
		{"valid without dot", "pdf", false},
		{"valid tar.gz", ".tar.gz", false},
		{"valid tar.gz no dot", "tar.gz", false},
		{"invalid empty", "", true},
		{"invalid dot only", ".", true},
		{"invalid special char", ".pdf!", true},
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
