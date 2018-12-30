package fetch

import (
	"path/filepath"
	"strings"

	"github.com/atecce/canon/lib"
	"github.com/gocolly/colly"
	"github.com/yhat/scrape"
)

type Fetcher interface {
	MkRoot() error
	MkAuthorDir(name string) error
	Fetch(url, path string) error
}

func Crawl(fetcher Fetcher) {

	if err := fetcher.MkRoot(); err != nil {
		lib.Log(nil, "", "", "ERR", "making root: "+err.Error())
	}

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := strings.Replace(e.ChildText("a"), "¶", "", -1)

		if err := fetcher.MkAuthorDir(author); err != nil {
			lib.Log(nil, "", "", "ERR", "making author directory: "+err.Error())
		}

		for _, node := range e.DOM.Next().Children().Nodes {

			child := node.FirstChild
			grandchild := child.FirstChild

			if grandchild != nil {

				url := domain + scrape.Attr(child, "href") + ".txt.utf-8"
				if strings.Contains(url, "wikipedia") {
					return
				}

				title := strings.Replace(strings.Replace(grandchild.Data, "/", "|", -1), "\n", "", -1)

				path := filepath.Join(author, title)

				if err := fetcher.Fetch(url, path); err != nil {
					lib.Log(nil, url, path, "ERR", "fetching: "+err.Error())
				}
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
