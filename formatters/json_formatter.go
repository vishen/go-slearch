package formatters

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/vishen/go-slearch/slearch"
)

func init() {
	slearch.Register("json", jsonLogFormatter{})
}

type jsonLogFormatter struct{}

func (j jsonLogFormatter) GetValueFromLine(config slearch.Config, line []byte, key string) (string, error) {
	trimedLine := bytes.TrimSpace(line)
	if trimedLine[0] != '{' && trimedLine[len(trimedLine)-1] != '}' {
		return "", slearch.ErrInvalidFormatForLine
	}
	keySplit := searchableKey(key, config.KeySplitString)
	vs, _, _, _ := jsonparser.Get(trimedLine, keySplit...)
	return fmt.Sprintf("%s", vs), nil
}

func (j jsonLogFormatter) FormatFoundValues(config slearch.Config, valuesFound []slearch.KV) string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, v := range valuesFound {
		buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", v.Key, v.Value))
		if i != len(valuesFound)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}")
	return buffer.String()

}

func searchableKey(key, splitKeyOnString string) []string {
	if splitKeyOnString == "" {
		return []string{key}
	}
	return strings.Split(key, splitKeyOnString)
}
