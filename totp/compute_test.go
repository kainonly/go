package totp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompute_InvalidSecret(t *testing.T) {
	assert.Equal(t, -1, Compute("!@#$$%", 100))
}
