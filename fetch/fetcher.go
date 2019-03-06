package fetch

import (
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/yhat/scrape"
)

type Fetcher interface {
	MkRoot() error
	MkAuthorDir(name string) error
	Fetch(url, path string) error
}

func Crawl(fetcher Fetcher) {

	if err := fetcher.MkRoot(); err != nil {
		logrus.Error("making root:", err)
	}

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := strings.Replace(e.ChildText("a"), "¶", "", -1)

		if err := fetcher.MkAuthorDir(author); err != nil {
			logrus.Error("making author directory:", err)
		}

		for _, node := range e.DOM.Next().Children().Nodes {

			child := node.FirstChild
			grandchild := child.FirstChild

			if grandchild != nil {

				url := domain + scrape.Attr(child, "href") + ".txt.utf-8"
				if strings.Contains(url, "wikipedia") {
					continue
				}

				title := strings.Replace(strings.Replace(grandchild.Data, "/", "|", -1), "\n", "", -1)

				path := filepath.Join(author, title)

				if err := fetcher.Fetch(url, path); err != nil {
					logrus.WithFields(logrus.Fields{
						"url":  url,
						"path": path,
					}).Error("fetching:", err)
				}
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
