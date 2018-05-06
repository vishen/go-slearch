package structured_log_search

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrNoMatchingKeyValues = errors.New("no matching key values found")
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
	FormatFoundValues(config Config, valuesFound []KV) string
}

func StructuredLoggingSearch(config Config) error {

	type lineResult struct {
		lineNumber uint64
		original   string
		result     string
		err        error
	}
	resultsChan := make(chan lineResult)

	doneChan := make(chan bool, 1)

	go func() {
		receivedLineResults := map[uint64]lineResult{}
		currentLineNumber := uint64(0)

		for lr := range resultsChan {
			receivedLineResults[lr.lineNumber] = lr

			for {
				foundLineResult, ok := receivedLineResults[currentLineNumber]
				if !ok {
					break
				}
				if foundLineResult.err != nil {
					if foundLineResult.err != ErrNoMatchingKeyValues {
						fmt.Printf("Error on line %d :%s : %s\n", foundLineResult.lineNumber, foundLineResult.original, foundLineResult.err)
					}
				} else {
					fmt.Println(foundLineResult.result)
				}
				currentLineNumber++
			}
		}

		doneChan <- true

	}()

	// TODO(vishen): Take anything an input and output interface into this func
	// that can be used in bufio.NewReader()
	reader := bufio.NewReader(os.Stdin)

	// TODO(vishen): Allow configuration to be able to use a max number
	// of goroutines
	wg := sync.WaitGroup{}

	for i := uint64(0); ; i++ {
		// TODO(vishen): This is super inefficient...
		text, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		wg.Add(1)
		go func(lineNumber uint64, line []byte) {
			defer wg.Done()
			result, err := searchLine(config, line)
			resultsChan <- lineResult{
				original:   string(line),
				lineNumber: lineNumber,
				result:     result,
				err:        err,
			}
		}(i, text[:len(text)-1])

	}

	wg.Wait()
	close(resultsChan)

	<-doneChan

	return nil
}

func searchLine(config Config, line []byte) (string, error) {
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
			return "", ErrNoMatchingKeyValues
		}

		if matched {
			found = matched
		}

		if len(valuesToPrint) > 0 {
			continue
		}

	}

	if !found && len(config.MatchOn) > 0 {
		return "", ErrNoMatchingKeyValues
	}

	for _, pk := range config.PrintKeys {
		pkv := formatter.GetValueFromLine(config, line, pk)
		if pkv == "" {
			continue
		}
		valuesToPrint = append(valuesToPrint, KV{Key: pk, Value: fmt.Sprintf("%s", pkv)})
	}

	// TODO(vishen): It is possible to have config.printKeys that don't match
	// any line, this should NOT print the entire line? Currently it kind of
	// seems alright to default to printing the line if no matching valuesToPrint
	// are found.
	if len(valuesToPrint) == 0 {
		if len(config.PrintKeys) == 0 {
			return string(line), nil
		}
		return "", ErrNoMatchingKeyValues
	} else {
		return formatter.FormatFoundValues(config, valuesToPrint), nil
	}
}
