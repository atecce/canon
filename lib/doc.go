package lib

import (
	"bufio"
	"compress/gzip"
	"net/http"
	"os"

	prose "gopkg.in/jdkato/prose.v2"
)

// Doc represents a document with named entities extracted
type Doc struct {
	Text     string
	Entities []Entity
}

// Entity contains the text and label along with the amount of occurences
type Entity struct {
	Text, Label string
	Count       uint
}

// NewDoc constucts a Doc with a url from gutenberg.org
func NewDoc(url, path string) (*Doc, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc := new(Doc)
	entities := make(map[prose.Entity]uint)
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	sc := bufio.NewScanner(res.Body)
	for sc.Scan() {

		chunk := sc.Text()

		_, err = w.Write([]byte(chunk))
		if err != nil {
			return nil, err
		}

		docChunk, err := prose.NewDocument(chunk)
		if err != nil {
			return nil, err
		}

		for _, ent := range docChunk.Entities() {
			if count, ok := entities[ent]; ok {
				entities[ent] = count + 1
			} else {
				entities[ent] = 1
			}
		}

		for ent, count := range entities {
			doc.Entities = append(doc.Entities, Entity{
				Text:  ent.Text,
				Label: ent.Label,
				Count: count,
			})
		}

	}

	// chomp the boilerplate at the end
	// corpus := string(text)
	// i := strings.Index(corpus, "End of the Project Gutenberg EBook")
	// if i == -1 {
	// 	Log(int64(len(corpus)), url, "", "WARN", "no license at end")
	// } else {
	// 	corpus = corpus[:i]
	// }

	return doc, nil
}
