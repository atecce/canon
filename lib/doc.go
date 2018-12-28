package lib

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"os"

	"github.com/jdkato/prose/tokenize"
)

// Doc represents a document with named entities extracted
type Doc struct {
	Tokens []string `json:"tokens"`
}

// WriteJSON serizlizes the doc to a gzipped file
func (doc *Doc) WriteJSON(path string) error {
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

// Entity contains the text and label along with the amount of occurences
type Entity struct {
	Text, Label string
	Count       uint
}

// NewDocFromPath constructs a doc from a path
func NewDocFromPath(path string) (*Doc, error) {

	doc := new(Doc)

	f, _ := os.Open(path)
	defer f.Close()

	r, _ := gzip.NewReader(f)
	defer r.Close()

	tokenizer := tokenize.NewTreebankWordTokenizer()

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		doc.Tokens = append(doc.Tokens, tokenizer.Tokenize(sc.Text())...)
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

// NewDocFromURL constucts a Doc from a url
// func NewDocFromURL(url, path string) (*Doc, error) {

// 	res, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	doc := new(Doc)
// 	entities := make(map[prose.Entity]uint)

// 	// TODO
// 	f, err := os.Create(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	w := gzip.NewWriter(f)
// 	defer w.Close()

// 	sc := bufio.NewScanner(res.Body)
// 	for sc.Scan() {

// 		chunk := sc.Text()

// 		if _, err := w.Write([]byte(chunk)); err != nil {
// 			return nil, err
// 		}

// 		docChunk, err := prose.NewDocument(chunk)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, ent := range docChunk.Entities() {
// 			if count, ok := entities[ent]; ok {
// 				entities[ent] = count + 1
// 			} else {
// 				entities[ent] = 1
// 			}
// 		}

// 		for ent, count := range entities {
// 			doc.Entities = append(doc.Entities, Entity{
// 				Text:  ent.Text,
// 				Label: ent.Label,
// 				Count: count,
// 			})
// 		}
// 	}

// 	// chomp the boilerplate at the end
// 	// corpus := string(text)
// 	// i := strings.Index(corpus, "End of the Project Gutenberg EBook")
// 	// if i == -1 {
// 	// 	Log(int64(len(corpus)), url, "", "WARN", "no license at end")
// 	// } else {
// 	// 	corpus = corpus[:i]
// 	// }

// 	return doc, nil
// }
