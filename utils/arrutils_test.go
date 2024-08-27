package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	assert := assert.New(t)
	assert.Equal([]int{2, 4}, Filter([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }))
	assert.Equal([]string{"b", "c", "d"}, Filter([]string{"a", "b", "c", "d"}, func(x string) bool { return x > "a" }))
}
