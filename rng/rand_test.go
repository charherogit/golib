package rng

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWightRandIndex(t *testing.T) {
	WeightRandI([]int{9, 9, 9})
	assert.Equal(t, WeightRandI([]int{1}), 0)
	assert.Equal(t, WeightRandI([]int{0, 0, 1}), 2)
	assert.Panics(t, func() { WeightRandI([]int{0, 0, 0}) })
	assert.Greater(t, WeightRandI([]int{0, 0, 5, 2, 9, 1}), 1)
	assert.Less(t, WeightRandI([]int{2, 1, 0, 0}), 2)
}

func TestRandIn(t *testing.T) {
	assert.Panics(t, func() { Range(0, 0) })
	assert.Panics(t, func() { Range(5, 5) })

	assert.Equal(t, Range(4, 5), 4)

	n := Range(0, 5)
	assert.True(t, n >= 0 && n < 5)

	n = Range(0, 1)
	assert.True(t, n >= 0 && n < 1)

	n = Range(2, 4)
	assert.True(t, n == 2 || n == 3)

	n = Range(5, 19)
	assert.True(t, n >= 5 && n < 19)
}
