package slearch

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testLogFormatter struct{}

func (t testLogFormatter) GetValueFromLine(config Config, line []byte, key string) (string, error) {

	return "", nil
}

func (t testLogFormatter) FormatFoundValues(config Config, valuesFound []KV) string {
	return ""
}

func Test_isSoftError(t *testing.T) {
	assert := assert.New(t)

	assert.True(isSoftError(ErrNoMatchingKeyValues))
	assert.True(isSoftError(ErrNoMatchingPrintValues))
	assert.False(isSoftError(ErrInvalidFormatForLine))
	assert.False(isSoftError(errors.New("some random error")))
}

func Test_searchLine(t *testing.T) {

}
