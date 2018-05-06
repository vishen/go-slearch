package structured_log_search

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type StructuredLogType int

const (
	StructuredLogTypeJson = StructuredLogType(iota)
	StructuredLogTypeText
)

type StructuredLogMatchType int

const (
	StructuredLogMatchTypeAnd = StructuredLogMatchType(iota)
	StructuredLogMatchTypeOr
)

type KV struct {
	Key         string
	Value       string
	RegexString string
}

type Config struct {
	// Whether this is an AND or OR matching
	MatchType StructuredLogMatchType

	// Json, text, etc...
	LogType StructuredLogType

	// Values to match on
	MatchOn []KV

	// Which keys to print for matching records
	PrintKeys []string

	// String to split the key on
	KeySplitString string
}

type StructuredLogFormatter interface {
	GetValueFromLine(config Config, line []byte, key string) string
	PrintFoundValues(config Config, valuesFound []KV)
}

func StructuredLoggingSearch(config Config) error {

	// TODO(vishen): Take anything an input and output interface into this func
	// that can be used in bufio.NewReader(
	reader := bufio.NewReader(os.Stdin)

	// TODO(vishen): Allow configuration to be able to use a max number
	// of goroutines
	wg := sync.WaitGroup{}

	for {
		// TODO(vishen): This is super inefficient...
		text, err := reader.ReadBytes('\n')
		if err != nil {
			wg.Wait()
			return errors.Wrapf(err, "error reading from stdin: %s", err)
		}

		wg.Add(1)
		go func(line []byte) {
			defer wg.Done()
			searchLine(config, line)
		}(text[:len(text)-1])

	}

	return nil
}

func searchableKey(key, splitKeyOnString string) []string {
	if splitKeyOnString == "" {
		return []string{key}
	}
	return strings.Split(key, splitKeyOnString)
}

func searchLine(config Config, line []byte) {
	valuesToPrint := make([]KV, 0, len(config.PrintKeys))

	// TODO(vishen): Change this to a register approach?
	var formatter StructuredLogFormatter
	switch config.LogType {
	case StructuredLogTypeText:
		formatter = textLogFormatter{}
	default:
		formatter = jsonLogFormatter{}
	}

	found := false
	for _, v := range config.MatchOn {
		foundValue := formatter.GetValueFromLine(config, line, v.Key)

		matched := false
		if v.Value != "" {
			matched = foundValue == v.Value
		} else if v.RegexString != "" {
			matched, _ = regexp.MatchString(v.RegexString, foundValue)
		}

		if !matched && config.MatchType == StructuredLogMatchTypeAnd {
			return
		}

		if matched {
			found = matched
		}

		if len(valuesToPrint) > 0 {
			continue
		}

	}

	if !found && len(config.MatchOn) > 0 {
		return
	}

	for _, pk := range config.PrintKeys {
		pkv := formatter.GetValueFromLine(config, line, pk)
		if pkv == "" {
			continue
		}
		valuesToPrint = append(valuesToPrint, KV{Key: pk, Value: fmt.Sprintf("%s", pkv)})
	}

	// TODO(vishen): It is possible to have config.printKeys that don't match
	// any line, this should NOT print the entire line!
	if len(valuesToPrint) == 0 {
		fmt.Printf("%s\n", line)
	} else {
		formatter.PrintFoundValues(config, valuesToPrint)
	}
}
