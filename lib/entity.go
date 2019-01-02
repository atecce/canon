package lib

import (
	"bufio"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

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
	Text  string `json:"text"`
	Count uint   `json:"count"`
}

func NewEnts(r io.ReadCloser) (*[]Entity, error) {

	defer r.Close()

	entities := make(map[string]uint)
	sc := bufio.NewScanner(r)
	for sc.Scan() {

		text := sc.Text()

		// chomp the boilerplate at the end
		i := strings.Index(text, "End of the Project Gutenberg EBook")
		if i == -1 {
			extractEntities(text, entities)
		} else {
			extractEntities(text[:i], entities)
			break
		}
	}

	var ents []Entity
	for ent, count := range entities {
		ents = append(ents, Entity{
			Text:  ent,
			Count: count,
		})
	}

	sort.Slice(ents, func(i, j int) bool {
		return ents[i].Count > ents[j].Count
	})

	return &ents, nil
}

func extractEntities(text string, entities map[string]uint) {
	for _, entity := range chunk.Chunk(tagger.Tag(tokenizer.Tokenize(text)), chunk.TreebankNamedEntities) {
		if count, ok := entities[entity]; ok {
			entities[entity] = count + 1
		} else {
			entities[entity] = 1
		}
	}
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
