package help_test

import (
	"testing"

	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestValidator_CustomRules(t *testing.T) {
	vd := help.Validator()
	type S struct {
		Name string `vd:"snake"`
		Sort string `vd:"sort"`
	}
	ok := S{Name: "user_name", Sort: "created_at:1"}
	err := vd.ValidateStruct(&ok)
	assert.NoError(t, err)

	bad1 := S{Name: "UserName", Sort: "created_at:1"}
	err = vd.ValidateStruct(&bad1)
	assert.Error(t, err)

	bad2 := S{Name: "user_name", Sort: "created_at:2"}
	err = vd.ValidateStruct(&bad2)
	assert.Error(t, err)
}
