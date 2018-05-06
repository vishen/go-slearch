package structured_log_search

import (
	"bytes"
	"fmt"
)

type textLogFormatter struct{}

func (t textLogFormatter) GetValueFromLine(config Config, line []byte, key string) string {

	// Loop through the characters in the `line` and find a matching `key`
	// and it's value, take into account some values might be surronded
	// in ' or " and have multiple spaces in the value.
	pos := 0
	for {
		if pos+len(key) > len(line) {
			break
		}

		if string(line[pos:pos+len(key)]) != key {
			pos += 1
			continue
		}

		// If the next character isn't the '=' then we don't
		// have a match
		if line[pos+len(key)] != '=' {
			pos += 1
			continue
		}

		// Eat '='
		pos += len(key) + 1

		var eatUntil byte = ' '
		switch line[pos] {
		case '"':
			eatUntil = '"'
			pos += 1
		case '\'':
			eatUntil = '\''
			pos += 1
		}
		startPos := pos

		// Eat up until the eatUntil character and return value
		for {
			if pos >= len(line) || line[pos] == eatUntil || line[pos] == '\n' {
				// If we are escaping, ignore
				if line[pos-1] != '\\' {
					break
				}
			}
			pos += 1
		}

		if startPos >= pos || pos > len(line) {
			return ""
		}

		return string(line[startPos:pos])
	}

	return ""
}

func (t textLogFormatter) FormatFoundValues(config Config, valuesFound []KV) string {
	var buffer bytes.Buffer
	for _, v := range valuesFound {
		buffer.WriteString(fmt.Sprintf("%s=\"%s\" ", v.Key, v.Value))
	}
	return buffer.String()
}
