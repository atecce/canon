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

	"github.com/gocolly/colly"
)

const domain = "https://www.gutenberg.org/"

var dir = filepath.Join("/", "keybase", "public", "atec", "data", "gutenberg")

func main() {

	authorCollector := colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		log.Println("[INFO] get", r.URL)
	})

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := filepath.Join(dir, e.ChildText("a"))
		if _, err := os.Stat(author); os.IsNotExist(err) {
			if mkErr := os.Mkdir(author, 0700); mkErr != nil {
				log.Println("[ERR]", err)
			}
		}

		// TODO pool of goroutines on a channel
		var wg sync.WaitGroup
		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				wg.Add(1)

				// TODO try again on err?
				go func(href, title string) {
					defer wg.Done()

					wwwURL := domain + href + ".txt.utf-8"
					kbURL := filepath.Join(author, strings.Replace(title, "/", "|", -1)+".txt.gz")

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
				}(scrape.Attr(node.FirstChild, "href"), node.FirstChild.FirstChild.Data)
			}
		}
		wg.Wait()
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
