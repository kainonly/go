package passlib_test

import (
	"testing"

	"github.com/kainonly/go/passlib"
	"github.com/stretchr/testify/assert"
)

func TestVerify_BadSlices(t *testing.T) {
	// malformed base64 that will trigger decode errors
	err := passlib.Verify("a", "$argon2id$v=19$m=65536,t=4,p=1$#$#")
	assert.Error(t, err)
}
