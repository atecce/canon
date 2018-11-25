package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"

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

		var wg sync.WaitGroup
		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				wg.Add(1)

				// TODO try again on err?
				go func(node html.Node) {
					defer wg.Done()

					wwwURL := domain + scrape.Attr(node.FirstChild, "href") + ".txt.utf-8"
					// TODO rm all forward slashes from title
					kbURL := filepath.Join(path, node.FirstChild.FirstChild.Data+".txt.gz")

					if strings.Contains(wwwURL, "wikipedia") {
						return
					}

					if _, err := os.Stat(kbURL); os.IsNotExist(err) {

						log.Println("[INFO] get", wwwURL)
						res, err := http.Get(wwwURL)
						if err != nil {
							log.Println("[ERR]", err)
							return
						}
						defer res.Body.Close()

						log.Println("[INFO] creating", kbURL)
						f, err := os.Create(kbURL)
						if err != nil {
							log.Println("[ERR]", err)
							return
						}
						defer f.Close()

						w := gzip.NewWriter(f)
						defer w.Close()

						log.Println("[INFO] copying", wwwURL, "to", kbURL)
						_, err = io.Copy(w, res.Body)
						if err != nil {
							log.Println("[ERR]", err)
							return
						}
					}
				}(*node)
			}
		}
		wg.Wait()
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
