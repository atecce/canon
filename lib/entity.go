package lib

import (
	"bufio"
	"compress/gzip"
	"io"
	"net/http"
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

func NewEnts(r io.ReadCloser) (*[]Entity, error) {

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

func NewEntsFromPath(path string) (*[]Entity, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}

	return NewEnts(r)
}

func NewEntsFromURL(url, path string) (*[]Entity, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return NewEnts(res.Body)
}
