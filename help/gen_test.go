package help_test

import (
	"github.com/google/uuid"
	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUuid(t *testing.T) {
	v := help.Uuid()
	_, err := uuid.Parse(v)
	assert.NoError(t, err)
}
