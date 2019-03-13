// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var segmenter *sentences.DefaultSentenceTokenizer

func newSentenceScanner(r io.Reader) *bufio.Scanner {

	segmenter, _ = english.NewSentenceTokenizer(nil)

	sc := bufio.NewScanner(r)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		sents := segmenter.Tokenize(string(data))
		if len(sents) > 0 {
			first := sents[0].Text
			return len(first), []byte(first), nil
		}

		return 0, nil, nil
	})

	return sc
}

var sentencesCmd = &cobra.Command{
	Use:   "sentences [dir]",
	Short: "segment sentences",
	Run: func(cmd *cobra.Command, args []string) {

		argc := len(args)

		if argc == 0 {

			sc := newSentenceScanner(os.Stdin)
			for sc.Scan() {
				println()
				println("BEGIN")
				os.Stdout.Write(sc.Bytes())
				os.Stdout.Write([]byte("\n"))
				println("END")
				println()
			}

		} else if argc == 1 {

			// sem := make(chan struct{}, 16)

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				println()
				os.Stdout.Write([]byte(path))
				os.Stdout.Write([]byte("\n"))
				println()

				f, _ := os.Open(path)
				defer f.Close()

				sc := newSentenceScanner(f)
				for sc.Scan() {
					println()
					println("BEGIN")
					os.Stdout.Write(sc.Bytes())
					os.Stdout.Write([]byte("\n"))
					println("END")
					println()
				}

				os.Stdout.Write([]byte(path))
				os.Stdout.Write([]byte("\n"))

				return nil
			})

		} else {
			panic("too many args")
		}
	},
}

func init() {
	rootCmd.AddCommand(sentencesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sentencesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sentencesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
