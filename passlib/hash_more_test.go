package passlib_test

import (
	"testing"

	"github.com/kainonly/go/passlib"
	"github.com/stretchr/testify/assert"
)

func TestHash_LongPassword(t *testing.T) {
	pw := make([]byte, 100)
	for i := 0; i < len(pw); i++ {
		pw[i] = 'a'
	}
	h, err := passlib.Hash(string(pw))
	assert.NoError(t, err)
	assert.NotEmpty(t, h)
}
