package help_test

import (
	"testing"

	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestIsEmpty_MoreBranches(t *testing.T) {
	assert.True(t, help.IsEmpty(map[string]string{}))
	assert.True(t, help.IsEmpty(([]int)(nil)))
	assert.True(t, help.IsEmpty(make([]int, 0)))
	type S struct{ A int }
	assert.True(t, help.IsEmpty(S{}))
	assert.False(t, help.IsEmpty(S{A: 1}))
	ch := make(chan int)
	assert.False(t, help.IsEmpty(ch))
}
