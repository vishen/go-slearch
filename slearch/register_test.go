package slearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	assert := assert.New(t)

	Register("test", testLogFormatter{})

	_, ok := formatters["test"]
	assert.True(ok)
}

func Test_getFormatter(t *testing.T) {
	assert := assert.New(t)

	Register("test", testLogFormatter{})

	_, ok := getFormatter("test")
	assert.True(ok)
}
