package structured_log_search

import (
	"bytes"
	"fmt"

	"github.com/buger/jsonparser"
)

type jsonLogFormatter struct{}

func (j jsonLogFormatter) GetValueFromLine(config Config, line []byte, key string) string {
	keySplit := searchableKey(key, config.KeySplitString)
	vs, _, _, _ := jsonparser.Get(line, keySplit...)
	return fmt.Sprintf("%s", vs)
}

func (j jsonLogFormatter) PrintFoundValues(config Config, valuesFound []KV) {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, v := range valuesFound {
		buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", v.Key, v.Value))
		if i != len(valuesFound)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}")
	fmt.Println(buffer.String())

}
