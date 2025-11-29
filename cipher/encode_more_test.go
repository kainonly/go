package cipher_test

import (
	"testing"

	"github.com/kainonly/go/cipher"
	"github.com/stretchr/testify/assert"
)

func TestEncode_EmptyData(t *testing.T) {
	x, err := cipher.New("6ixSiEXaqxsJTozbnxQ76CWdZXB2JazK")
	assert.NoError(t, err)
	ct, err := x.Encode([]byte{})
	assert.NoError(t, err)
	assert.NotEmpty(t, ct)
}
