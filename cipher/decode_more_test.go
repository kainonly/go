package cipher_test

import (
	"testing"

	"github.com/kainonly/go/cipher"
	"github.com/stretchr/testify/assert"
)

func TestDecode_BadBase64(t *testing.T) {
	x, err := cipher.New("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK")
	assert.NoError(t, err)
	_, err = x.Decode("@@@notbase64@@@")
	assert.Error(t, err)
}
