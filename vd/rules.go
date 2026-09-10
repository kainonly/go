package vd

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// ============================================================================
// 中国本地化验证规则
// ============================================================================

// BankCard 校验中国银行卡号（16-19 位数字，含 Luhn 校验）。
// 示例："6222021234567890123"
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

// luhnCheck 使用 Luhn 算法校验数字字符串。
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

// LicensePlate 校验中国车牌号。
// 支持普通车牌和新能源汽车车牌。
// 示例："京A12345"、"沪A12345D"、"粤B123456"
func LicensePlate() Rule {
	return Rule{
		Tag: "license_plate",
		Fn:  licensePlateValidation,
	}
}

var licensePlateRegex = sync.OnceValue(func() *regexp.Regexp {
	// 普通车牌：京A12345（省份+字母后 5 个字符）
	// 新能源车牌：京A123456 或 京AD12345（6 个字符）
	return regexp.MustCompile(`^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼][A-Z][A-Z0-9]{5,6}$`)
})

func licensePlateValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return licensePlateRegex().MatchString(s)
}

// USCC 校验统一社会信用代码。
// 中国组织机构使用的 18 位代码。
// 示例："91310000MA1FL8TQ32"
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

// ChineseWord 校验字符串仅包含中文字符。
// 示例："中国"、"你好世界"
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

// ChineseName 校验中文姓名（2-6 个汉字，可含 ·）。
// 示例："张三"、"欧阳修"、"古力娜扎·迪丽热巴"
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
// 密码强度验证规则
// ============================================================================

// PasswordWeak 校验弱密码（至少 6 个字符）。
// 示例："123456"
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

// PasswordMedium 校验中等强度密码。
// 至少 8 个字符，须同时包含字母和数字。
// 示例："password123"
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

// PasswordStrong 校验强密码。
// 至少 8 个字符，须同时包含大写字母、小写字母、数字和特殊字符。
// 示例："Password123!"
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
// 特殊格式验证规则
// ============================================================================

// ObjectID 校验 MongoDB ObjectId（24 位十六进制字符）。
// 注意：validator 自带的 "mongodb" 用于连接字符串，此规则用于 ObjectId。
// 示例："507f1f77bcf86cd799439011"
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

// Snowflake 校验 Snowflake ID（正整数，通常 18-19 位）。
// 示例："1234567890123456789"
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

// Version 校验语义化版本格式。
// 注意：validator 自带 "semver"，此规则是更简单的版本检查。
// 示例："1.0.0"、"2.1.3"
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
// 数据结构验证规则
// ============================================================================

// SafeString 校验字符串不含注入攻击相关的危险字符。
// 拒绝：< > " ' ` ; & | $ \ 以及空字节。
// 示例："safe text 123"
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

// AlphaNumDash 校验可含连字符和下划线的字母数字字符串。
// 比 alphanum 更宽松，适用于 slug 和标识符。
// 示例："my-item_123"
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

// AlphaNumSpace 校验可含空格的字母数字字符串。
// 示例："Hello World 123"
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

// Decimal 校验十进制数字字符串（可含小数部分）。
// 示例："123.45"、"0.001"
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

// PositiveDecimal 校验正的十进制数字字符串。
// 示例："123.45"、"0.001"
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
	// 确认是真正的正数（而非 "0.000"）
	f, err := strconv.ParseFloat(s, 64)
	return err == nil && f > 0
}

// ============================================================================
// 常用格式验证规则
// ============================================================================

// Domain 校验不含协议的域名。
// 注意：validator 自带 "fqdn"，此规则是更简单的域名检查。
// 示例："example.com"、"sub.example.co.uk"
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

// FilePath 校验文件路径格式（Unix 或 Windows 风格）。
// 示例："/home/user/file.txt"、"C:\Users\file.txt"
func FilePath() Rule {
	return Rule{
		Tag: "file_path",
		Fn:  filePathValidation,
	}
}

var filePathRegex = sync.OnceValue(func() *regexp.Regexp {
	// Unix：/path/to/file 或相对路径
	// Windows：C:\path\to\file 或 \\server\share
	return regexp.MustCompile(`^([a-zA-Z]:\\|\\\\|/)?[\w\-. /\\]+$`)
})

func filePathValidation(fl FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok || s == "" {
		return false
	}
	return filePathRegex().MatchString(s)
}

// FileName 校验文件名（不含路径分隔符）。
// 示例："document.pdf"、"image_01.png"
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

// FileExt 校验文件扩展名（可带或不带点）。
// 示例：".pdf"、"pdf"、".tar.gz"
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

// Color 校验颜色代码（不带 # 前缀的十六进制）。
// 注意：validator 自带的 "hexcolor" 带 # 前缀，此规则不带。
// 示例："ff5733"、"FFF"、"abc123"
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
// 通信验证规则
// ============================================================================

// TelPhone 校验座机号码（中国格式）。
// 示例："010-12345678"、"0755-1234567"、"02112345678"
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

// QQ 校验 QQ 号（5-11 位数字，不以 0 开头）。
// 示例："12345"、"1234567890"
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

// WeChat 校验微信号格式。
// 6-20 个字符，以字母开头，仅限字母、数字和下划线。
// 示例："wxid_abc123"、"myWeChat_01"
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
// 地址验证规则
// ============================================================================

// ZipCode 校验中国邮政编码（6 位数字）。
// 示例："100000"、"518000"
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
// 代码命名规则
// ============================================================================

// Variable 校验变量名（编程命名惯例）。
// 以字母或下划线开头，仅限字母、数字和下划线。
// 示例："myVar"、"_private"、"MAX_VALUE"
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

// PascalCase 校验 PascalCase 格式。
// 示例："MyClass"、"UserService"
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

// CamelCase 校验 camelCase 格式。
// 示例："myVariable"、"getUserName"
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

// KebabCase 校验 kebab-case 格式。
// 示例："my-component"、"user-profile-card"
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

// UpperSnake 校验 UPPER_SNAKE_CASE 格式。
// 示例："MAX_VALUE"、"HTTP_STATUS_OK"
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
// 工具函数
// ============================================================================

// All 返回所有可用的自定义规则。
func All() []Rule {
	return []Rule{
		// 原有规则
		Snake(),
		Sort(),
		Phone(),
		IDCard(),
		Username(),
		Slug(),
		// 中国本地化
		BankCard(),
		LicensePlate(),
		USCC(),
		ChineseWord(),
		ChineseName(),
		// 密码强度
		PasswordWeak(),
		PasswordMedium(),
		PasswordStrong(),
		// 特殊格式
		ObjectID(),
		Snowflake(),
		Version(),
		// 数据结构
		SafeString(),
		AlphaNumDash(),
		AlphaNumSpace(),
		Decimal(),
		PositiveDecimal(),
		// 常用格式
		Domain(),
		FilePath(),
		FileName(),
		FileExt(),
		Color(),
		// 通信
		TelPhone(),
		QQ(),
		WeChat(),
		// 地址
		ZipCode(),
		// 代码命名
		Variable(),
		PascalCase(),
		CamelCase(),
		KebabCase(),
		UpperSnake(),
	}
}

// Common 返回常用规则（snake、sort、phone、idcard、username、slug、password_medium）。
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

// Chinese 返回所有中国本地化规则。
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

// NamingConvention 返回所有命名约定规则。
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

// NotBlank 校验字符串非空且非纯空白字符。
// 示例："hello"（有效）、"  "（无效）
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
