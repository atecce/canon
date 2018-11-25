package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yhat/scrape"

	"github.com/gocolly/colly"
)

const domain = "https://www.gutenberg.org/"

var dir = filepath.Join("/", "keybase", "public", "atec", "data", "gutenberg")

func main() {

	// TODO pick up where you left off

	authorCollector := colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		log.Println("[INFO] get", r.URL)
	})

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := e.ChildText("a")

		path := filepath.Join(dir, author)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if mkErr := os.Mkdir(path, 0700); mkErr != nil {
				log.Println("[ERR]", err)
			}
		}

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				wwwURL := domain + scrape.Attr(node.FirstChild, "href") + ".txt.utf-8"
				kbURL := filepath.Join(path, node.FirstChild.FirstChild.Data+".txt.gz")

				if strings.Contains(wwwURL, "wikipedia") {
					continue
				}

				if _, err := os.Stat(kbURL); os.IsNotExist(err) {

					log.Println("[INFO] getting", wwwURL)
					res, err := http.Get(wwwURL)
					if err != nil {
						log.Println("[ERR]", err)
						continue
					}
					defer res.Body.Close()

					log.Println("[INFO] creating", kbURL)
					f, err := os.Create(kbURL)
					if err != nil {
						log.Println("[ERR]", err)
						continue
					}
					defer f.Close()

					w := gzip.NewWriter(f)
					defer w.Close()

					log.Println("[INFO] copying")
					_, err = io.Copy(w, res.Body)
					if err != nil {
						log.Println("[ERR]", err)
						continue
					}
				}
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
