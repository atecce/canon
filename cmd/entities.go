// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/lib"
	"github.com/atecce/io"
	"github.com/spf13/cobra"
)

var entitiesCmd = &cobra.Command{
	Use:   "entities [dir]",
	Short: "extract entities",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			ents, err := lib.NewEnts(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to construct entities: %v\n", err)
				os.Exit(1)
			}
			if err := json.NewEncoder(os.Stdout).Encode(&ents); err != nil {
				fmt.Fprintf(os.Stderr, "failed to encode json: %v\n", err)
				os.Exit(1)
			}
		} else {
			sem := make(chan struct{}, 16)

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if strings.Contains(path, ".txt.") {

					println(path)

					sem <- struct{}{}

					go func(textPath, jsonPath string) {
						defer func() {
							<-sem
						}()

						println(textPath)

						ents, err := lib.NewEntsFromPath(textPath)
						if err != nil {
							log.Fatal(err)
						}

						if err := io.WriteJSON(jsonPath, &ents); err != nil {
							log.Fatal(err)
						}

					}(path, strings.Replace(path, ".txt.", ".json.", -1))
				}

				return nil
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(entitiesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
