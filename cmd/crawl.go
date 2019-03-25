package cmd

import (
	"os"

	"atec.pub/canon/fetch"

	"github.com/spf13/cobra"
)

func usage() {
	println("usage: canon crawl [files | entities]")
	os.Exit(1)
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "crawl gutenberg for corpora",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			usage()
		}

		var gzipExt string
		if *gzipFlag {
			gzipExt = ".gz"
		}

		var fetcher fetch.Fetcher
		switch args[0] {
		case "files":
			fetcher = &fetch.FileFetcher{
				Root: "gutenberg",
				Sem:  make(chan struct{}, 10),
				Ext:  ".txt" + gzipExt,
			}
		case "entities":
			fetcher = &fetch.EntitiesFetcher{
				Root: "gutenberg",
				Sem:  make(chan struct{}, 10),
				Ext:  ".json" + gzipExt,
			}
		default:
			usage()
		}

		fetch.Crawl(fetcher)
	},
}

var gzipFlag *bool

func init() {
	rootCmd.AddCommand(crawlCmd)

	gzipFlag = crawlCmd.Flags().Bool("gzip", false, "gzips files")
}
