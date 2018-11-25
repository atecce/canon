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
			_ = os.Mkdir(path, 0700)
		}

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				title := node.FirstChild.FirstChild.Data
				href := scrape.Attr(node.FirstChild, "href")

				if strings.Contains(href, "wikipedia") {
					continue
				}

				url := domain + href + ".txt.utf-8"
				name := filepath.Join(path, title+".txt.gz")

				log.Println("[INFO] getting", url)
				res, err := http.Get(url)
				if err != nil {
					log.Println("[ERR]", err)
					continue
				}
				defer res.Body.Close()

				log.Println("[INFO] creating", name)
				f, err := os.Create(name)
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
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
