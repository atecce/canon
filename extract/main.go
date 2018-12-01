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

	"github.com/atecce/canon/common"
)

const domain = "https://www.gutenberg.org/"

type document struct {
	Text     string
	Entities []entity
}

type entity struct {
	Text, Label string
	Count       uint
}

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
	var wg sync.WaitGroup

	authorCollector := colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		log.Println("[INFO] get", r.URL)
	})

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := filepath.Join(common.Dir, e.ChildText("a"))
		if _, err := os.Stat(author); os.IsNotExist(err) {
			if mkErr := os.Mkdir(author, 0700); mkErr != nil {
				log.Println("[ERR]", err)
			}
		}

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				wg.Add(1)
				sem <- struct{}{}

				// TODO try again on err?
				go func(href, title string) {
					defer wg.Done()
					defer func() {
						<-sem
					}()

					// strip forward slash and new lines
					name := common.StripNewlines(strings.Replace(title, "/", "|", -1))

					wwwURL := domain + href + ".txt.utf-8"
					if strings.Contains(wwwURL, "wikipedia") {
						return
					}

					kbPath := filepath.Join(author, name+".txt.gz")
					if _, err := os.Stat(kbPath); os.IsNotExist(err) {
						log.Println("[INFO]", kbPath, "not on kbfs. fetching from", wwwURL)
						if err := fetch(wwwURL, kbPath); err != nil {
							log.Println("[ERR] fetching:", err)
						}
					}
				}(scrape.Attr(node.FirstChild, "href"), node.FirstChild.FirstChild.Data)
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}

	wg.Wait()
}
