package help_test

import (
	"testing"

	"github.com/kainonly/go/help"
	"github.com/stretchr/testify/assert"
)

func TestIsEmpty_FullKinds(t *testing.T) {
	// string
	assert.True(t, help.IsEmpty(""))
	assert.False(t, help.IsEmpty("a"))
	// array
	var arr0 [0]int
	var arr1 = [1]int{1}
	assert.True(t, help.IsEmpty(arr0))
	assert.False(t, help.IsEmpty(arr1))
	// map
	var mNil map[string]int
	m0 := map[string]int{}
	m1 := map[string]int{"a": 1}
	assert.True(t, help.IsEmpty(mNil))
	assert.True(t, help.IsEmpty(m0))
	assert.False(t, help.IsEmpty(m1))
	// slice
	var sNil []int
	s0 := make([]int, 0)
	s1 := []int{1}
	assert.True(t, help.IsEmpty(sNil))
	assert.True(t, help.IsEmpty(s0))
	assert.False(t, help.IsEmpty(s1))
	// bool
	assert.True(t, help.IsEmpty(false))
	assert.False(t, help.IsEmpty(true))
	// int/uint
	assert.True(t, help.IsEmpty(int(0)))
	assert.False(t, help.IsEmpty(int(1)))
	assert.True(t, help.IsEmpty(uint(0)))
	assert.False(t, help.IsEmpty(uint(1)))
	// float
	assert.True(t, help.IsEmpty(float32(0)))
	assert.False(t, help.IsEmpty(float32(1)))
	assert.True(t, help.IsEmpty(float64(0)))
	assert.False(t, help.IsEmpty(float64(1)))
	// interface/ptr
	var i interface{}
	assert.True(t, help.IsEmpty(i))
	p := help.Ptr(0)
	assert.False(t, help.IsEmpty(p))
	// func
	var fn func()
	assert.True(t, help.IsEmpty(fn))
	fn = func() {}
	assert.False(t, help.IsEmpty(fn))
}
