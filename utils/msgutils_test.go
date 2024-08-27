package utils_test

import (
	"testing"

	"github.com/nekoite/go-napcat/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsRawMessageApiResp(t *testing.T) {
	assert := assert.New(t)
	assert.True(utils.IsRawMessageApiResp([]byte(`{"echo": "hello"}`)))
	assert.True(utils.IsRawMessageApiResp([]byte(`{"echo": "hello", "extra": "world"}`)))
	assert.False(utils.IsRawMessageApiResp([]byte(`{"something": "1"}`)))
}
