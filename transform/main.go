package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/jdkato/prose.v2"

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

func newDoc(url string) (*document, error) {

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

	proseDoc, err := prose.NewDocument(corpus)
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

func writeJSON(doc *document, url string) error {
	f, err := os.Create(url)
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	if err := json.NewEncoder(w).Encode(doc); err != nil {
		return err
	}
	return nil
}

func main() {

	// TODO pool of goroutines on a channel
	var wg sync.WaitGroup
	filepath.Walk(common.Dir, func(path string, info os.FileInfo, err error) error {

		// TODO try again on err?

		if !strings.Contains(path, ".txt") {
			return nil
		}

		JSONURL := strings.Replace(path, ".txt.", ".json.", -1)
		if _, err := os.Stat(JSONURL); os.IsNotExist(err) {

			log.Println("[INFO]", JSONURL, "not on kbfs. extracting doc")
			doc, err := newDoc(path)
			if err != nil {
				log.Println("[ERR] extracting doc for", path+":", err)
				return nil
			}

			log.Println("[INFO] writing", JSONURL)
			if err := writeJSON(doc, JSONURL); err != nil {
				log.Println("[ERR] writing doc to json for", path+":", err)
				return nil
			}
		}
		return nil
	})
	wg.Wait()
}
