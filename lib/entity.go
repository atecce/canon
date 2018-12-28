package lib

import (
	"bufio"
	"compress/gzip"
	"os"

	"github.com/jdkato/prose/chunk"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
)

var (
	tokenizer = tokenize.NewTreebankWordTokenizer()
	tagger    = tag.NewPerceptronTagger()
)

// Entity contains the text and the amount of occurences
type Entity struct {
	Text  string
	Count uint
}

func NewEntsFromPath(path string) (*[]Entity, error) {

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

	entities := make(map[string]uint)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		for _, entity := range chunk.Chunk(tagger.Tag(tokenizer.Tokenize(sc.Text())), chunk.TreebankNamedEntities) {
			if count, ok := entities[entity]; ok {
				entities[entity] = count + 1
			} else {
				entities[entity] = 1
			}
		}
	}

	var ents []Entity
	for ent, count := range entities {
		ents = append(ents, Entity{
			Text:  ent,
			Count: count,
		})
	}

	// chomp the boilerplate at the end
	// corpus := string(text)
	// i := strings.Index(corpus, "End of the Project Gutenberg EBook")
	// if i == -1 {
	// 	Log(int64(len(corpus)), url, "", "WARN", "no license at end")
	// } else {
	// 	corpus = corpus[:i]
	// }

	return &ents, nil
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
