package fetch

import (
	"github.com/atecce/canon/lib"
	"github.com/gocolly/colly"
)

const domain = "https://www.gutenberg.org/"

var authorCollector *colly.Collector

func init() {
	authorCollector = colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		lib.Log(nil, r.URL.Path, "", "INFO", r.Method)
	})
}
