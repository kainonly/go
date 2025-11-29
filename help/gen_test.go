package help_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestUuid(t *testing.T) {
	v := help.Uuid()
	_, err := uuid.Parse(v)
	assert.NoError(t, err)
}

func TestSID(t *testing.T) {
	v1 := help.SID()
	v2 := help.SID()
	assert.NotEmpty(t, v1)
	assert.NotEmpty(t, v2)
	assert.NotEqual(t, v1, v2)
}
