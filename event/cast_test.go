package event_test

import (
	"testing"

	"github.com/nekoite/go-napcat/event"
	"github.com/stretchr/testify/assert"
)

type testStruct struct{}

type testStruct2 struct{}

func TestGetAs(t *testing.T) {
	assert := assert.New(t)
	expected := &testStruct{}
	actual := event.GetAs[testStruct](expected)
	assert.Equal(expected, actual)
	actual = event.GetAs[testStruct](nil)
	assert.Nil(actual)
	actual2 := event.GetAs[testStruct2](expected)
	assert.Nil(actual2)
}

func TestGetAsUnsafe(t *testing.T) {
	assert := assert.New(t)
	expected := &testStruct{}
	actual := event.GetAsUnsafe[testStruct](expected)
	assert.Equal(expected, actual)
	assert.Panics(func() {
		event.GetAsUnsafe[testStruct2](expected)
	})
}

func TestGetAsOrError(t *testing.T) {
	assert := assert.New(t)
	expected := &testStruct{}
	actual, err := event.GetAsOrError[testStruct](expected)
	assert.Equal(expected, actual)
	assert.Nil(err)
	actual, err = event.GetAsOrError[testStruct](nil)
	assert.Nil(actual)
	assert.NotNil(err)
	actual2, err := event.GetAsOrError[testStruct2](expected)
	assert.Nil(actual2)
	assert.NotNil(err)
}
