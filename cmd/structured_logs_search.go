// Copyright © 2018 Jonathan Pentecost <pentecostjonathan@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

type StructuredLogType int

const (
	StructuredLogTypeJson = StructuredLogType(iota)
)

type StructuredLogMatchType int

const (
	StructuredLogMatchTypeAnd = StructuredLogMatchType(iota)
	StructuredLogMatchTypeOr
)

type KV struct {
	key         string
	value       string
	regexString string
}

type Config struct {
	// Whether this is an AND or OR matching
	matchType StructuredLogMatchType

	// Json, text, etc...
	logType StructuredLogType

	// Values to match on
	matchOn []KV

	// Which keys to print for matching records
	printKeys []string

	// String to split the key on
	keySplitString string
}

func StructuredLoggingSearch(config Config) error {

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
			searchLine(line, config)
		}(text[:len(text)-1])

	}

	return nil
}

func searchableKey(key, splitKeyOnString string) []string {
	return strings.Split(key, splitKeyOnString)
}

func searchLine(line []byte, config Config) {
	valuesToPrint := make([]KV, 0, len(config.printKeys))

	found := false
	for _, v := range config.matchOn {
		// TODO(vishen): Is there a better way to check equality than forcing
		// everything to strings and comparing them?
		keySplit := searchableKey(v.key, config.keySplitString)
		vs, _, _, _ := jsonparser.Get(line, keySplit...)

		matched := false
		if v.value != "" {
			matched = fmt.Sprintf("%s", vs) == v.value
		} else if v.regexString != "" {
			matched, _ = regexp.MatchString(v.regexString, fmt.Sprintf("%s", vs))
		}

		if !matched && config.matchType == StructuredLogMatchTypeAnd {
			return
		}

		if matched {
			found = matched
		}

		if len(valuesToPrint) > 0 {
			continue
		}

		for _, pk := range config.printKeys {
			pkv, _, _, err := jsonparser.Get(line, pk)
			if err != nil {
				continue
			}
			valuesToPrint = append(valuesToPrint, KV{key: pk, value: fmt.Sprintf("%s", pkv)})
		}

	}

	if !found && len(config.matchOn) > 0 {
		return
	}

	if len(valuesToPrint) == 0 {
		fmt.Printf("%s\n", line)
	} else {
		var buffer bytes.Buffer
		buffer.WriteString("{")
		for i, v := range valuesToPrint {
			buffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"", v.key, v.value))
			if i != len(valuesToPrint)-1 {
				buffer.WriteString(", ")
			}
		}
		buffer.WriteString("}")
		fmt.Println(buffer.String())
	}

}
