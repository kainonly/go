package totp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckCode_WindowAdvance(t *testing.T) {
	x := &Totp{Secret: "JBSWY3DPEHPK3PXP", Window: 2, Counter: 1}
	code := Compute(x.Secret, int64(x.Counter))
	ok := x.CheckCode(code)
	assert.True(t, ok)
	assert.Equal(t, 2, x.Counter)
}
