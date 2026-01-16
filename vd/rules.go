package vd

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// ============================================================================
// Chinese Localization Rules (中国本地化验证规则)
// ============================================================================

// BankCard validates Chinese bank card number (16-19 digits with Luhn check).
// Example: "6222021234567890123"
func BankCard() Rule {
	return Rule{
		Tag: "bankcard",
		Fn:  bankcardValidation,
	}
}

var bankcardRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d{16,19}$`)
})

func bankcardValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || !bankcardRegex().MatchString(s) {
		return false
	}
	return luhnCheck(s)
}

// luhnCheck validates a number string using Luhn algorithm.
func luhnCheck(s string) bool {
	var sum int
	alt := false
	for i := len(s) - 1; i >= 0; i-- {
		n := int(s[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

// LicensePlate validates Chinese vehicle license plate number.
// Supports regular plates and new energy vehicle plates.
// Example: "京A12345", "沪A12345D", "粤B123456"
func LicensePlate() Rule {
	return Rule{
		Tag: "license_plate",
		Fn:  licensePlateValidation,
	}
}

var licensePlateRegex = sync.OnceValue(func() *regexp.Regexp {
	// Regular: 京A12345 (5 chars after province+letter)
	// New energy: 京A123456 or 京AD12345 (6 chars)
	return regexp.MustCompile(`^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼][A-Z][A-Z0-9]{5,6}$`)
})

func licensePlateValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return licensePlateRegex().MatchString(s)
}

// USCC validates Unified Social Credit Code (统一社会信用代码).
// 18-character code for Chinese organizations.
// Example: "91310000MA1FL8TQ32"
func USCC() Rule {
	return Rule{
		Tag: "uscc",
		Fn:  usccValidation,
	}
}

var usccRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$`)
})

func usccValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || len(s) != 18 {
		return false
	}
	return usccRegex().MatchString(s)
}

// ChineseWord validates that string contains only Chinese characters.
// Example: "中国", "你好世界"
func ChineseWord() Rule {
	return Rule{
		Tag: "chinese",
		Fn:  chineseWordValidation,
	}
}

func chineseWordValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

// ChineseName validates Chinese person name (2-6 Chinese characters, may include ·).
// Example: "张三", "欧阳修", "古力娜扎·迪丽热巴"
func ChineseName() Rule {
	return Rule{
		Tag: "chinese_name",
		Fn:  chineseNameValidation,
	}
}

var chineseNameRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[\p{Han}]{1,6}(·[\p{Han}]{1,6})*$`)
})

func chineseNameValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return chineseNameRegex().MatchString(s)
}

// ============================================================================
// Password Strength Rules (密码强度验证规则)
// ============================================================================

// PasswordWeak validates weak password (at least 6 chars).
// Example: "123456"
func PasswordWeak() Rule {
	return Rule{
		Tag: "password_weak",
		Fn:  passwordWeakValidation,
	}
}

func passwordWeakValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return len(s) >= 6
}

// PasswordMedium validates medium strength password.
// At least 8 chars, must contain letters and numbers.
// Example: "password123"
func PasswordMedium() Rule {
	return Rule{
		Tag: "password_medium",
		Fn:  passwordMediumValidation,
	}
}

func passwordMediumValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || len(s) < 8 {
		return false
	}
	var hasLetter, hasDigit bool
	for _, r := range s {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}

// PasswordStrong validates strong password.
// At least 8 chars, must contain uppercase, lowercase, number, and special char.
// Example: "Password123!"
func PasswordStrong() Rule {
	return Rule{
		Tag: "password_strong",
		Fn:  passwordStrongValidation,
	}
}

func passwordStrongValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || len(s) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

// ============================================================================
// Special Format Rules (特殊格式验证规则)
// ============================================================================

// ObjectID validates MongoDB ObjectId (24 hex characters).
// Note: validator has "mongodb" for connection string, this is for ObjectId.
// Example: "507f1f77bcf86cd799439011"
func ObjectID() Rule {
	return Rule{
		Tag: "objectid",
		Fn:  objectIDValidation,
	}
}

var objectIDRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
})

func objectIDValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return objectIDRegex().MatchString(s)
}

// Snowflake validates Snowflake ID (positive integer, typically 18-19 digits).
// Example: "1234567890123456789"
func Snowflake() Rule {
	return Rule{
		Tag: "snowflake",
		Fn:  snowflakeValidation,
	}
}

func snowflakeValidation(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind().String() {
	case "string":
		s := field.Interface().(string)
		if len(s) < 15 || len(s) > 19 {
			return false
		}
		_, err := strconv.ParseInt(s, 10, 64)
		return err == nil
	case "int64":
		return field.Int() > 0
	case "uint64":
		return field.Uint() > 0
	}
	return false
}

// Version validates semantic version format.
// Note: validator has "semver" but this is a simpler version check.
// Example: "1.0.0", "2.1.3"
func Version() Rule {
	return Rule{
		Tag: "version",
		Fn:  versionValidation,
	}
}

var versionRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d+\.\d+\.\d+$`)
})

func versionValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return versionRegex().MatchString(s)
}

// ============================================================================
// Data Structure Rules (数据结构验证规则)
// ============================================================================

// SafeString validates string contains no dangerous characters for injection.
// Rejects: < > " ' ` ; & | $ \ and null bytes.
// Example: "safe text 123"
func SafeString() Rule {
	return Rule{
		Tag: "safe_string",
		Fn:  safeStringValidation,
	}
}

var unsafeChars = []rune{'<', '>', '"', '\'', '`', ';', '&', '|', '$', '\\', '\x00'}

func safeStringValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	for _, r := range s {
		for _, unsafe := range unsafeChars {
			if r == unsafe {
				return false
			}
		}
	}
	return true
}

// AlphaNumDash validates alphanumeric string with dashes and underscores.
// More permissive than alphanum, useful for slugs and identifiers.
// Example: "my-item_123"
func AlphaNumDash() Rule {
	return Rule{
		Tag: "alphanumdash",
		Fn:  alphaNumDashValidation,
	}
}

var alphaNumDashRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
})

func alphaNumDashValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return alphaNumDashRegex().MatchString(s)
}

// AlphaNumSpace validates alphanumeric string with spaces.
// Example: "Hello World 123"
func AlphaNumSpace() Rule {
	return Rule{
		Tag: "alphanumspace",
		Fn:  alphaNumSpaceValidation,
	}
}

var alphaNumSpaceRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
})

func alphaNumSpaceValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return alphaNumSpaceRegex().MatchString(s)
}

// Decimal validates decimal number string with optional precision.
// Example: "123.45", "0.001"
func Decimal() Rule {
	return Rule{
		Tag: "decimal",
		Fn:  decimalValidation,
	}
}

var decimalRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^-?\d+(\.\d+)?$`)
})

func decimalValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return decimalRegex().MatchString(s)
}

// PositiveDecimal validates positive decimal number string.
// Example: "123.45", "0.001"
func PositiveDecimal() Rule {
	return Rule{
		Tag: "positive_decimal",
		Fn:  positiveDecimalValidation,
	}
}

var positiveDecimalRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d+(\.\d+)?$`)
})

func positiveDecimalValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" || s == "0" || s == "0.0" {
		return false
	}
	if !positiveDecimalRegex().MatchString(s) {
		return false
	}
	// Check it's actually positive (not just "0.000")
	f, err := strconv.ParseFloat(s, 64)
	return err == nil && f > 0
}

// ============================================================================
// Common Format Rules (常用格式验证规则)
// ============================================================================

// Domain validates domain name without protocol.
// Note: validator has "fqdn" but this is simpler domain check.
// Example: "example.com", "sub.example.co.uk"
func Domain() Rule {
	return Rule{
		Tag: "domain",
		Fn:  domainValidation,
	}
}

var domainRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
})

func domainValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || len(s) > 253 {
		return false
	}
	return domainRegex().MatchString(s)
}

// FilePath validates file path format (Unix or Windows style).
// Example: "/home/user/file.txt", "C:\Users\file.txt"
func FilePath() Rule {
	return Rule{
		Tag: "file_path",
		Fn:  filePathValidation,
	}
}

var filePathRegex = sync.OnceValue(func() *regexp.Regexp {
	// Unix: /path/to/file or relative/path
	// Windows: C:\path\to\file or \\server\share
	return regexp.MustCompile(`^([a-zA-Z]:\\|\\\\|/)?[\w\-. /\\]+$`)
})

func filePathValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return filePathRegex().MatchString(s)
}

// FileName validates file name (no path separators).
// Example: "document.pdf", "image_01.png"
func FileName() Rule {
	return Rule{
		Tag: "filename",
		Fn:  fileNameValidation,
	}
}

var fileNameRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[^<>:"/\\|?*\x00-\x1f]+$`)
})

func fileNameValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" || s == "." || s == ".." {
		return false
	}
	if len(s) > 255 {
		return false
	}
	return fileNameRegex().MatchString(s)
}

// FileExt validates file extension (with or without dot).
// Example: ".pdf", "pdf", ".tar.gz"
func FileExt() Rule {
	return Rule{
		Tag: "file_ext",
		Fn:  fileExtValidation,
	}
}

var fileExtRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^\.?[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*$`)
})

func fileExtValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" || s == "." {
		return false
	}
	return fileExtRegex().MatchString(s)
}

// Color validates color code (hex without # prefix).
// Note: validator has "hexcolor" with #, this is without.
// Example: "ff5733", "FFF", "abc123"
func Color() Rule {
	return Rule{
		Tag: "color",
		Fn:  colorValidation,
	}
}

var colorRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^([a-fA-F0-9]{3}|[a-fA-F0-9]{6}|[a-fA-F0-9]{8})$`)
})

func colorValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return colorRegex().MatchString(s)
}

// ============================================================================
// Communication Rules (通信相关验证规则)
// ============================================================================

// TelPhone validates telephone number (landline, Chinese format).
// Example: "010-12345678", "0755-1234567", "02112345678"
func TelPhone() Rule {
	return Rule{
		Tag: "tel",
		Fn:  telPhoneValidation,
	}
}

var telPhoneRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^0\d{2,3}-?\d{7,8}$`)
})

func telPhoneValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return telPhoneRegex().MatchString(s)
}

// QQ validates QQ number (5-11 digits, not starting with 0).
// Example: "12345", "1234567890"
func QQ() Rule {
	return Rule{
		Tag: "qq",
		Fn:  qqValidation,
	}
}

var qqRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[1-9]\d{4,10}$`)
})

func qqValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return qqRegex().MatchString(s)
}

// WeChat validates WeChat ID format.
// 6-20 chars, starts with letter, alphanumeric and underscore only.
// Example: "wxid_abc123", "myWeChat_01"
func WeChat() Rule {
	return Rule{
		Tag: "wechat",
		Fn:  weChatValidation,
	}
}

var weChatRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{5,19}$`)
})

func weChatValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return weChatRegex().MatchString(s)
}

// ============================================================================
// Address Rules (地址相关验证规则)
// ============================================================================

// ZipCode validates Chinese postal code (6 digits).
// Example: "100000", "518000"
func ZipCode() Rule {
	return Rule{
		Tag: "zipcode",
		Fn:  zipCodeValidation,
	}
}

var zipCodeRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^\d{6}$`)
})

func zipCodeValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return zipCodeRegex().MatchString(s)
}

// ============================================================================
// Code Rules (编码相关验证规则)
// ============================================================================

// Variable validates variable name (programming convention).
// Starts with letter or underscore, alphanumeric and underscore only.
// Example: "myVar", "_private", "MAX_VALUE"
func Variable() Rule {
	return Rule{
		Tag: "variable",
		Fn:  variableValidation,
	}
}

var variableRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
})

func variableValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return variableRegex().MatchString(s)
}

// PascalCase validates PascalCase format.
// Example: "MyClass", "UserService"
func PascalCase() Rule {
	return Rule{
		Tag: "pascal",
		Fn:  pascalCaseValidation,
	}
}

var pascalCaseRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
})

func pascalCaseValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return pascalCaseRegex().MatchString(s)
}

// CamelCase validates camelCase format.
// Example: "myVariable", "getUserName"
func CamelCase() Rule {
	return Rule{
		Tag: "camel",
		Fn:  camelCaseValidation,
	}
}

var camelCaseRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)
})

func camelCaseValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return camelCaseRegex().MatchString(s)
}

// KebabCase validates kebab-case format.
// Example: "my-component", "user-profile-card"
func KebabCase() Rule {
	return Rule{
		Tag: "kebab",
		Fn:  kebabCaseValidation,
	}
}

var kebabCaseRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
})

func kebabCaseValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return kebabCaseRegex().MatchString(s)
}

// UpperSnake validates UPPER_SNAKE_CASE format.
// Example: "MAX_VALUE", "HTTP_STATUS_OK"
func UpperSnake() Rule {
	return Rule{
		Tag: "upper_snake",
		Fn:  upperSnakeValidation,
	}
}

var upperSnakeRegex = sync.OnceValue(func() *regexp.Regexp {
	return regexp.MustCompile(`^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$`)
})

func upperSnakeValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return upperSnakeRegex().MatchString(s)
}

// ============================================================================
// Utility Functions
// ============================================================================

// All returns all available custom rules.
func All() []Rule {
	return []Rule{
		// Original rules
		Snake(),
		Sort(),
		Phone(),
		IDCard(),
		Username(),
		Slug(),
		// Chinese localization
		BankCard(),
		LicensePlate(),
		USCC(),
		ChineseWord(),
		ChineseName(),
		// Password strength
		PasswordWeak(),
		PasswordMedium(),
		PasswordStrong(),
		// Special formats
		ObjectID(),
		Snowflake(),
		Version(),
		// Data structure
		SafeString(),
		AlphaNumDash(),
		AlphaNumSpace(),
		Decimal(),
		PositiveDecimal(),
		// Common formats
		Domain(),
		FilePath(),
		FileName(),
		FileExt(),
		Color(),
		// Communication
		TelPhone(),
		QQ(),
		WeChat(),
		// Address
		ZipCode(),
		// Code naming
		Variable(),
		PascalCase(),
		CamelCase(),
		KebabCase(),
		UpperSnake(),
	}
}

// Common returns commonly used rules (snake, sort, phone, idcard, username, slug, password_medium).
func Common() []Rule {
	return []Rule{
		Snake(),
		Sort(),
		Phone(),
		IDCard(),
		Username(),
		Slug(),
		PasswordMedium(),
	}
}

// Chinese returns all Chinese localization rules.
func Chinese() []Rule {
	return []Rule{
		Phone(),
		IDCard(),
		BankCard(),
		LicensePlate(),
		USCC(),
		ChineseWord(),
		ChineseName(),
		TelPhone(),
		QQ(),
		WeChat(),
		ZipCode(),
	}
}

// NamingConvention returns all naming convention rules.
func NamingConvention() []Rule {
	return []Rule{
		Snake(),
		PascalCase(),
		CamelCase(),
		KebabCase(),
		UpperSnake(),
		Variable(),
	}
}

// NotBlank validates string is not empty and not only whitespace.
// Example: "hello" (valid), "  " (invalid)
func NotBlank() Rule {
	return Rule{
		Tag: "notblank",
		Fn:  notBlankValidation,
	}
}

func notBlankValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(s) != ""
}
