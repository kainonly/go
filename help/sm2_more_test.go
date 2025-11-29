package help_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/smx509"
	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

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
