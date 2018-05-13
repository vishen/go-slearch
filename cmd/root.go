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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	// Load the structured log search formatters
	_ "github.com/vishen/go-slearch/formatters"
	"github.com/vishen/go-slearch/slearch"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-slearch",
	Short: "Search structured logs",
	Long:  `Read stuctured logs from STDIN and filter out lines based on exact and regex matches. Currently only supports JSON logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, _ := cmd.Flags().GetStringSlice("match")
		r, _ := cmd.Flags().GetStringSlice("regexp")
		k, _ := cmd.Flags().GetStringSlice("print_keys")
		t, _ := cmd.Flags().GetString("type")
		s, _ := cmd.Flags().GetString("search_type")
		d, _ := cmd.Flags().GetString("key_delimiter")
		v, _ := cmd.Flags().GetBool("verbose")

		config := slearch.Config{}

		makeKVSlice := func(values []string, regex bool) []slearch.KV {
			prevKey := ""
			kvs := make([]slearch.KV, 0, len(m))
			for _, v := range values {
				vSplit := strings.SplitN(v, "=", 2)

				var key string
				var value string

				if len(vSplit) == 1 {
					// TODO(): Should error if `prevKey` is empty
					key = strings.TrimSpace(prevKey)
					value = strings.TrimSpace(vSplit[0])
				} else {
					key = strings.TrimSpace(vSplit[0])
					value = strings.TrimSpace(vSplit[1])
				}
				prevKey = key

				kv := slearch.KV{Key: key}
				if regex {
					kv.RegexString = value
				} else {
					kv.Value = value
				}

				kvs = append(kvs, kv)

			}

			return kvs
		}

		config.MatchOn = makeKVSlice(m, false)
		config.MatchOn = append(config.MatchOn, makeKVSlice(r, true)...)
		config.PrintKeys = k
		config.KeySplitString = d
		config.Verbose = v

		if t == "" {
			t = "json"
		}
		config.LogFormatterType = strings.ToLower(t)

		if strings.ToLower(s) == "or" {
			config.MatchType = slearch.StructuredLogMatchTypeOr
		} else {
			config.MatchType = slearch.StructuredLogMatchTypeAnd
		}

		if err := slearch.StructuredLoggingSearch(config, os.Stdin, os.Stdout); err != nil {
			log.Printf("error searching structured logs: %s\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("type", "t", "json", "the log type to use: 'json' or 'text'")
	rootCmd.Flags().StringP("search_type", "s", "and", "the search type to use: 'and' or 'or'")
	rootCmd.Flags().StringP("key_delimiter", "d", "", "the string to split the key on for complex key queries")
	rootCmd.Flags().StringSliceP("match", "m", []string{}, "key and value to match on. eg: label1=value1")
	rootCmd.Flags().StringSliceP("regexp", "r", []string{}, "key and value to regex match on. eg: label1=value*")
	rootCmd.Flags().StringSliceP("print_keys", "p", []string{}, "keys to print if a match is found")
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")
}
