package help_test

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestOkFail(t *testing.T) {
	r1 := help.Ok()
	assert.Equal(t, int64(0), r1.Code)
	assert.Equal(t, "ok", r1.Message)

	r2 := help.Fail(1234, "oops")
	assert.Equal(t, int64(1234), r2.Code)
	assert.Equal(t, "oops", r2.Message)
}

func TestE(t *testing.T) {
	e := help.E(2002, "bad")
	assert.True(t, e.IsType(errors.ErrorTypePublic))
	meta, ok := e.Meta.(*help.ErrorMeta)
	assert.True(t, ok)
	assert.Equal(t, int64(2002), meta.Code)
}
