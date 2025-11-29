package help_test

import (
	"testing"

	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestMapToSignText_MoreTypes(t *testing.T) {
	m := map[string]any{
		"a": float32(1.5),
		"b": nil,
		"c": []int{1, 2},
		"d": struct{ X int }{X: 2},
		"e": uint(3),
	}
	r := help.MapToSignText(m)
	assert.Contains(t, r, "a=1.5")
	assert.Contains(t, r, "c=[1 2]")
	assert.Contains(t, r, "d={2}")
	assert.Contains(t, r, "e=3")
}
