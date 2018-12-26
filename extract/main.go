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

	"github.com/atecce/canon/lib"
)

const domain = "https://www.gutenberg.org/"

func fetchFiles(root string) error {

	sem := make(chan struct{}, 10)

	authorCollector := colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		lib.Log(0, r.URL.Path, "", "INFO", r.Method)
	})

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		// remove pilcrows from author name
		author := filepath.Join(root, strings.Replace(e.ChildText("a"), "¶", "", -1))

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

					// remove forward slashes and new lines
					name := lib.RemoveNewlines(strings.Replace(title, "/", "|", -1))

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

	return nil
}

func fetchTarball(root string) error {
	return nil
}

func fetch(url, path string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(path)
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

	if err := fetchFiles("gutenberg"); err != nil {
		log.Fatal(err)
	}
}
