package structured_log_search

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type jsonLogFormatter struct{}

func (j jsonLogFormatter) GetValueFromLine(config Config, line []byte, key string) string {
	keySplit := searchableKey(key, config.KeySplitString)
	vs, _, _, _ := jsonparser.Get(line, keySplit...)
	return fmt.Sprintf("%s", vs)
}

func (j jsonLogFormatter) FormatFoundValues(config Config, valuesFound []KV) string {
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
