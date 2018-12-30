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
	"os"

	"github.com/atecce/canon/fetch"
	"github.com/spf13/cobra"
)

func usage() {
	println("usage: canon crawl [files | tarball]")
	os.Exit(1)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "crawl gutenberg for corpora",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			usage()
		}

		var fetcher fetch.Fetcher
		switch args[0] {
		case "files":
			fetcher = &fetch.FileFetcher{
				Root: "gutenberg",
				Sem:  make(chan struct{}, 10),
			}
			break
		case "tarball":
			fetcher = &fetch.TarballFetcher{
				Root: "gutenberg.tar.gz",
			}
			break
		default:
			usage()
		}

		fetch.Crawl(fetcher)
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
