package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yhat/scrape"

	"github.com/gocolly/colly"
)

func main() {

	const domain = "https://www.gutenberg.org/"

	authorCollector := colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("[INFO] get", r.URL)
	})

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := e.ChildText("a")

		path := filepath.Join("/", "keybase", "public", "atec", "data", "gutenberg", author)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = os.Mkdir(path, 0700)
		}

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				title := node.FirstChild.FirstChild.Data
				href := scrape.Attr(node.FirstChild, "href")

				url := domain + href + ".txt.utf-8"
				name := filepath.Join(path, title+".txt.gz")

				fmt.Println("[INFO] getting", url)
				res, err := http.Get(url)
				if err != nil {
					fmt.Println("[ERR]", err)
					continue
				}
				defer res.Body.Close()

				fmt.Println("[INFO] creating", name)
				f, err := os.Create(name)
				if err != nil {
					fmt.Println("[ERR]", err)
					continue
				}
				defer f.Close()

				w := gzip.NewWriter(f)
				defer w.Close()

				fmt.Println("[INFO] copying")
				_, err = io.Copy(w, res.Body)
				if err != nil {
					fmt.Println("[ERR]", err)
					continue
				}
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
