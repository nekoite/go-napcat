package utils_test

import (
	"testing"

	"github.com/nekoite/go-napcat/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSet[int]()
	assert.NotNil(s)
}

func TestNewSetFrom(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSetFrom(1, 2, 3)
	assert.NotNil(s)
	assert.EqualValues(map[int]utils.Void{1: {}, 2: {}, 3: {}}, s)
}

func TestSet_Add(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSet[int]()
	s.Add(1)
	assert.EqualValues(map[int]utils.Void{1: {}}, s)
}

func TestSet_Contains(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSetFrom(1, 2, 3)
	assert.True(s.Contains(1))
	assert.True(s.Contains(2))
	assert.True(s.Contains(3))
	assert.False(s.Contains(4))
}

func TestSet_Remove(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSetFrom(1, 2, 3)
	s.Remove(1)
	assert.EqualValues(map[int]utils.Void{2: {}, 3: {}}, s)
}

func TestSet_Len(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSetFrom(1, 2, 3)
	assert.Equal(3, s.Len())
}

func TestSet_Clear(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSetFrom(1, 2, 3)
	assert.Equal(3, s.Len())
	s.Clear()
	assert.Equal(0, s.Len())
	assert.Empty(s)
}

func TestSet_ToSlice(t *testing.T) {
	assert := assert.New(t)
	s := utils.NewSetFrom(1, 2, 3)
	assert.ElementsMatch([]int{1, 2, 3}, s.ToSlice())
}
