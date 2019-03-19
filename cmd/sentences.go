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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"atec.pub/canon/lib"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var sentencesCmd = &cobra.Command{
	Use:   "sentences [dir]",
	Short: "segment sentences",
	Run: func(cmd *cobra.Command, args []string) {

		argc := len(args)

		if argc == 0 {

			sc := lib.NewSentenceScanner(os.Stdin)
			for sc.Scan() {
				println()
				println("BEGIN")
				os.Stdout.Write(sc.Bytes())
				os.Stdout.Write([]byte("\n"))
				println("END")
				println()
			}

		} else if argc == 1 {

			// TODO sort out err handling and logging

			sem := make(chan struct{}, 16)

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if info.IsDir() {
					return nil
				}

				sem <- struct{}{}
				go func(path string, info os.FileInfo) {
					defer func() {
						<-sem
					}()

					f, err := os.Open(path)
					if err != nil {
						pretty.Println(err)
						return
					}
					defer f.Close()

					var i uint
					sc := lib.NewSentenceScanner(f)
					for sc.Scan() {

						sentence := struct {
							// TODO possibly put with concatenation of these id
							Author string `json:"author"`
							Work   string `json:"work"`
							I      uint   `json:"i"`

							Text string `json:"text"`
						}{

							filepath.Base(filepath.Dir(path)),
							strings.TrimSuffix(info.Name(), ".txt"),
							i,

							sc.Text(),
						}

						b, err := json.Marshal(sentence)
						if err != nil {
							pretty.Println(err)
							continue
						}

						pretty.Println(string(b))

						req, err := http.NewRequest(http.MethodPost, "http://localhost:9200/sentences/_doc/", bytes.NewReader(b))
						if err != nil {
							pretty.Println(err)
							continue
						}
						req.Header.Add("Content-Type", "application/json")

						res, err := http.DefaultClient.Do(req)
						if err != nil {
							pretty.Println(err)
							continue
						}

						pretty.Println(res.Status)

						b, err = ioutil.ReadAll(res.Body)
						if err != nil {
							pretty.Println(err)
							continue
						}

						pretty.Println(string(b))

						i++

						println()
					}

				}(path, info)

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
