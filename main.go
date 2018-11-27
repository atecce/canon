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

func fetch(wwwURL, kbURL string) error {

	res, err := http.Get(wwwURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(kbURL)
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

func extractDoc(url string) (*document, error) {

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

	// chomp the boilerplate at the end
	corpus := string(text)
	i := strings.Index(corpus, "End of the Project Gutenberg EBook")
	if i == -1 {
		log.Println("[WARN] no license at end of", url)
	} else {
		corpus = corpus[:i]
	}

	proseDoc, err := prose.NewDocument(string(text))
	if err != nil {
		return nil, err
	}

	entities := make(map[prose.Entity]uint)
	for _, ent := range proseDoc.Entities() {
		if count, ok := entities[ent]; ok {
			entities[ent] = count + 1
		} else {
			entities[ent] = 1
		}
	}

	doc := document{
		Text: proseDoc.Text,
	}
	for ent, count := range entities {
		doc.Entities = append(doc.Entities, entity{
			Text:  ent.Text,
			Label: ent.Label,
			Count: count,
		})
	}

	return &doc, nil
}

func main() {

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

		// TODO pool of goroutines on a channel
		var wg sync.WaitGroup
		for i, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				wg.Add(1)

				// TODO try again on err?
				go func(href, title string) {
					defer wg.Done()

					// strip forward slash and new lines
					name := common.StripNewlines(strings.Replace(title, "/", "|", -1))

					wwwURL := domain + href + ".txt.utf-8"
					if strings.Contains(wwwURL, "wikipedia") {
						return
					}

					kbTextURL := filepath.Join(author, name+".txt.gz")
					if _, err := os.Stat(kbTextURL); os.IsNotExist(err) {
						log.Println("[INFO]", kbTextURL, "not on kbfs. fetching from", wwwURL)
						if err := fetch(wwwURL, kbTextURL); err != nil {
							log.Println("[ERR] fetching:", err)
						}
					}

					kbJSONURL := filepath.Join(author, name+".json.gz")
					if _, err := os.Stat(kbJSONURL); os.IsNotExist(err) {

						log.Println("[INFO]", kbJSONURL, "not on kbfs. extracting doc")
						doc, err := extractDoc(kbTextURL)
						if err != nil {
							log.Println("[ERR] extracting doc for", kbTextURL+":", err)
							return
						}

						log.Println("[INFO] writing", kbJSONURL)
						f, err := os.Create(kbJSONURL)
						if err != nil {
							log.Println("[ERR] creating file at", kbJSONURL+":", err)
							return
						}
						defer f.Close()

						w := gzip.NewWriter(f)
						defer w.Close()

						if err := json.NewEncoder(w).Encode(doc); err != nil {
							log.Println("[ERR] encoding json at", kbJSONURL+":", err)
							return
						}
					}
				}(scrape.Attr(node.FirstChild, "href"), node.FirstChild.FirstChild.Data)

				if i != 0 && i%10 == 0 {
					wg.Wait()
				}
			}
		}
		wg.Wait()
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}
