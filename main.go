package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/jdkato/prose.v2"

	"github.com/yhat/scrape"

	"github.com/gocolly/colly"
)

// TODO rm licenses at the end

const domain = "https://www.gutenberg.org/"

var dir = filepath.Join("/", "keybase", "public", "atec", "data", "gutenberg")

func readDoc(url string) (*prose.Document, error) {

	f, err := os.Open(url)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	text, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	doc, err := prose.NewDocument(string(text))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

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

					name := strings.Replace(title, "/", "|", -1)

					wwwURL := domain + href + ".txt.utf-8"
					if strings.Contains(wwwURL, "wikipedia") {
						return
					}

					kbTextURL := filepath.Join(author, name+".txt.gz")
					if _, err := os.Stat(kbTextURL); os.IsNotExist(err) {

						log.Println("[INFO]", kbTextURL, "not on kbfs. fetching")

						log.Println("[INFO] get", wwwURL)
						res, err := http.Get(wwwURL)
						if err != nil {
							log.Println("[ERR]", err)
							return
						}
						defer res.Body.Close()

						log.Println("[INFO] create", kbTextURL)
						f, err := os.Create(kbTextURL)
						if err != nil {
							log.Println("[ERR]", err)
							return
						}
						defer f.Close()

						w := gzip.NewWriter(f)
						defer w.Close()

						log.Println("[INFO] copy", wwwURL, "to", kbTextURL)
						if _, err := io.Copy(w, res.Body); err != nil {
							log.Println("[ERR] copy", wwwURL, "to", kbTextURL, ":", err)
							return
						}
					}

					kbJSONURL := filepath.Join(author, name+".json.gz")
					if _, err := os.Stat(kbJSONURL); os.IsNotExist(err) {

						log.Println("[INFO]", kbJSONURL, "not on kbfs. extracting entities")

						doc, err := readDoc(kbTextURL)
						if err != nil {
							log.Println("[ERR] reading", kbTextURL, ":", err)
							return
						}

						log.Println("[INFO] create", kbJSONURL)
						f, err := os.Create(kbJSONURL)
						if err != nil {
							log.Println("[ERR]", err)
							return
						}
						defer f.Close()

						w := gzip.NewWriter(f)
						defer w.Close()

						log.Println("[INFO] encode", kbJSONURL)
						if err := json.NewEncoder(w).Encode(doc.Entities()); err != nil {
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
