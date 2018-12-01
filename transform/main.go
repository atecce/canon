package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

func newDoc(path string) (*document, error) {

	f, err := os.Open(path)
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
		log(path, "", "WARN", "no license at end")
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

func writeJSON(doc *document, path string) error {
	f, err := os.Create(path)
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

func log(in, out, level, msg string) {
	b, _ := json.Marshal(struct {
		Time  time.Time `json:"time"`
		In    string    `json:"in"`
		Out   string    `json:"out"`
		Level string    `json:"level"`
		Msg   string    `json:"msg"`
	}{
		time.Now(),
		in,
		out,
		level,
		msg,
	})
	fmt.Println(string(b))
}

func main() {

	// TODO pool of goroutines on a channel
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	filepath.Walk(common.Dir, func(textPath string, info os.FileInfo, err error) error {

		// TODO try again on err?

		if !strings.Contains(textPath, ".txt") {
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(jsonPath string) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			if _, err := os.Stat(jsonPath); os.IsNotExist(err) {

				log(textPath, jsonPath, "INFO", "not on kbfs. extracting doc")
				doc, err := newDoc(textPath)
				if err != nil {
					log(textPath, jsonPath, "ERR", "extracting doc: "+err.Error())
					return
				}

				log(textPath, jsonPath, "INFO", "writing")
				if err := writeJSON(doc, jsonPath); err != nil {
					log(textPath, jsonPath, "ERR", "writing: "+err.Error())
					return
				}
			}
		}(strings.Replace(textPath, ".txt.", ".json.", -1))

		return nil
	})
	wg.Wait()
}
