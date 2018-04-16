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

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

type KV struct {
	key   string
	value string
}

type Config struct {
	// Whether this is an AND or OR matching
	matchType string

	// Values to match on
	matchOn []KV

	// Json, text, etc...
	logType string

	// Which keys to print for matching records
	printKeys []string
}

func StructuredLoggingGrep(config Config) error {

	reader := bufio.NewReader(os.Stdin)

	//wg := sync.WaitGroup{}

	config.matchOn = []KV{
		KV{"key2", "value2"},
		KV{"key5", "5"},
	}

	config.printKeys = []string{"key2", "key3", "key4", "key5", "complexKey"}
	config.matchType = "OR"

	for {
		// TODO(vishen): This is super inefficient...
		text, err := reader.ReadBytes('\n')
		if err != nil {
			//wg.Wait()
			return errors.Wrapf(err, "error reading from stdin: %s", err)
		}

		/*
			wg.Add(1)
			go func(line []byte) {
				defer wg.Done()
				searchLine(line, config)
			}(text[:len(text)-1])
		*/

		searchLine(text[:len(text)-1], config)
	}

	return nil
}

func searchLine(line []byte, config Config) {
	valuesToPrint := make([]KV, 0, len(config.printKeys))

	matched := false
	for _, v := range config.matchOn {
		// TODO(vishen): Is there a better way to check equality than forcing
		// everything to strings and comparing them?
		vs, _, _, _ := jsonparser.Get(line, v.key)
		if fmt.Sprintf("%s", vs) == v.value {
			matched = true
			for _, pk := range config.printKeys {
				pkv, _, _, err := jsonparser.Get(line, pk)
				if err != nil {
					continue
				}
				valuesToPrint = append(valuesToPrint, KV{pk, fmt.Sprintf("%s", pkv)})
			}
		} else if config.matchType == "AND" {
			return
		}
	}

	if !matched {
		return
	}

	if len(valuesToPrint) == 0 {
		fmt.Printf("%s\n", line)
	} else {
		var buffer bytes.Buffer
		buffer.WriteString("{")
		for i, v := range valuesToPrint {
			buffer.WriteString(fmt.Sprintf("\"%s\": \"%s\"", v.key, v.value))
			if i != len(valuesToPrint)-1 {
				buffer.WriteString(", ")
			}
		}
		buffer.WriteString("}")
		fmt.Println(buffer.String())
	}

}
