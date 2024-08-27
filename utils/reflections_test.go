package utils_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nekoite/go-napcat/utils"
	"github.com/stretchr/testify/assert"
)

func TestWalkStructWithTag(t *testing.T) {
	assert := assert.New(t)
	type S struct {
		A int `json:"a"`
		B struct {
			C int `json:"c"`
		} `json:"b"`
		D any `json:"d"`
	}
	s := S{A: 1, B: struct {
		C int `json:"c"`
	}{C: 2}, D: struct {
		E string `json:"e"`
	}{E: "e"}}
	var keys []string
	err := utils.WalkStructLeafWithTag(s, func(v reflect.Value, tagPath []reflect.StructTag) error {
		keys = append(keys, tagPath[len(tagPath)-1].Get("json"))
		return nil
	})
	assert.Nil(err)
	assert.Equal([]string{"a", "c", "e"}, keys)

	keys = nil
	err = utils.WalkStructLeafWithTag(&s, func(v reflect.Value, tagPath []reflect.StructTag) error {
		keys = append(keys, tagPath[len(tagPath)-1].Get("json"))
		return nil
	})
	assert.Nil(err)
	assert.Equal([]string{"a", "c", "e"}, keys)
}

func TestWalkStructLeafWithTagError(t *testing.T) {
	assert := assert.New(t)
	type S struct {
		A int `json:"a"`
		B struct {
			C int `json:"c"`
		} `json:"b"`
	}
	s := S{A: 1, B: struct {
		C int `json:"c"`
	}{C: 2}}
	err := utils.WalkStructLeafWithTag(nil, func(v reflect.Value, tagPath []reflect.StructTag) error {
		return nil
	})
	assert.Nil(err)

	err = utils.WalkStructLeafWithTag(s, func(v reflect.Value, tagPath []reflect.StructTag) error {
		return errors.New("error")
	})
	assert.NotNil(err)
}
