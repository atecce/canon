package lib

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"strings"

	prose "gopkg.in/jdkato/prose.v2"
)

type Doc struct {
	Text     string
	Entities []Entity
}

type Entity struct {
	Text, Label string
	Count       uint
}

func NewDoc(path string) (*Doc, error) {

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
		Log(int64(len(corpus)), path, "", "WARN", "no license at end")
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

	doc := Doc{
		Text: proseDoc.Text,
	}
	for ent, count := range entities {
		doc.Entities = append(doc.Entities, Entity{
			Text:  ent.Text,
			Label: ent.Label,
			Count: count,
		})
	}

	return &doc, nil
}
