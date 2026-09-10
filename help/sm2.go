package help

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/smx509"
)

// SM2 相关错误。
var (
	ErrSM2InvalidPublicKey  = errors.New("sm2: invalid public key format")
	ErrSM2InvalidPrivateKey = errors.New("sm2: invalid private key format")
	ErrSM2InvalidSignature  = errors.New("sm2: invalid signature format")
)

// SM2UID 是 SM2 签名/验签的默认用户 ID。
// 即标准的 16 字节 UID "1234567812345678"。
var SM2UID = []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}

// PubKeySM2FromBase64 解析 base64 编码的 SM2 公钥。
// 密钥应为 PKIX 格式（DER 编码后再 base64）。
func PubKeySM2FromBase64(v string) (*ecdsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, err
	}
	key, err := smx509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	pubKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrSM2InvalidPublicKey
	}
	return pubKey, nil
}

// PrivKeySM2FromBase64 解析 base64 编码的 SM2 私钥。
// 密钥应为 PKCS8 格式（DER 编码后再 base64）。
func PrivKeySM2FromBase64(v string) (*sm2.PrivateKey, error) {
	der, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, err
	}
	key, err := smx509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	priKey, ok := key.(*sm2.PrivateKey)
	if !ok {
		return nil, ErrSM2InvalidPrivateKey
	}
	return priKey, nil
}

// Sm2Sign 使用 SM2 私钥对文本签名。
// 返回 base64 编码的 ASN.1 DER 签名。
func Sm2Sign(key *sm2.PrivateKey, text string) (string, error) {
	signature, err := key.Sign(rand.Reader, []byte(text), sm2.DefaultSM2SignerOpts)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Sm2Verify 使用 SM2 公钥验证签名。
// 签名应为 base64 编码的 ASN.1 DER 格式。
// 签名有效时返回 true。
func Sm2Verify(pubKey *ecdsa.PublicKey, text string, sign string) (bool, error) {
	b, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}
	return sm2.VerifyASN1WithSM2(pubKey, SM2UID, []byte(text), b), nil
}
