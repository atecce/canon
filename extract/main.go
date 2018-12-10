package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yhat/scrape"

	"github.com/gocolly/colly"

	"github.com/atecce/canon/lib"
)

const domain = "https://www.gutenberg.org/"

func fetch(wwwURL, kbPath string) error {

	res, err := http.Get(wwwURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(kbPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}

	return nil
}

func main() {

	// TODO pool of goroutines on a channel
	sem := make(chan struct{}, 10)

	authorCollector := colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		lib.Log(0, r.URL.Path, "", "INFO", r.Method)
	})

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := filepath.Join(lib.Dir, e.ChildText("a"))
		if _, err := os.Stat(author); os.IsNotExist(err) {
			if mkErr := os.Mkdir(author, 0700); mkErr != nil {
				lib.Log(0, author, "", "ERR", "failed to mkdir: "+err.Error())
			}
		}

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				sem <- struct{}{}

				// TODO try again on err?
				go func(href, title string) {
					defer func() {
						<-sem
					}()

					// remove forward slash and new lines
					name := strings.Replace(lib.RemoveNewlines(strings.Replace(title, "/", "|", -1)), "¶", "", -1)

					wwwURL := domain + href + ".txt.utf-8"
					if strings.Contains(wwwURL, "wikipedia") {
						return
					}

					kbPath := filepath.Join(author, name+".txt.gz")
					lib.Log(0, wwwURL, kbPath, "INFO", "checking for kbPath")
					if _, err := os.Stat(kbPath); os.IsNotExist(err) {
						lib.Log(0, wwwURL, kbPath, "INFO", "not on kbfs. fetching")
						if err := fetch(wwwURL, kbPath); err != nil {
							lib.Log(0, wwwURL, kbPath, "ERR", "fetching: "+err.Error())
						}
					}
				}(scrape.Attr(node.FirstChild, "href"), node.FirstChild.FirstChild.Data)
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
