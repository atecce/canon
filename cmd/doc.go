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
	"os"
	"path/filepath"

	"github.com/atecce/canon/lib"
	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use: "doc",
	Run: func(cmd *cobra.Command, args []string) {
		path := filepath.Join(lib.Dir, args[0])
		doc, _ := lib.NewDoc(path)
		json.NewEncoder(os.Stdout).Encode(&doc)
	},
}

func init() {
	rootCmd.AddCommand(docCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
