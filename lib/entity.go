package lib

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jdkato/prose/chunk"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var (
	tokenizer = tokenize.NewTreebankWordTokenizer()
	tagger    = tag.NewPerceptronTagger()

	re        *regexp.Regexp
	segmenter *sentences.DefaultSentenceTokenizer
)

func init() {
	segmenter, _ = english.NewSentenceTokenizer(nil)

	re = regexp.MustCompile(`\n`)
}

func NewEnts(r io.ReadCloser) (map[string]uint, error) {

	defer r.Close()

	entities := make(map[string]uint)
	sc := bufio.NewScanner(r)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		// text := string(re.ReplaceAll(data, []byte(" ")))

		sents := segmenter.Tokenize(string(data))
		if len(sents) > 0 {
			first := sents[0].Text
			return len(first), []byte(first), nil
		}

		return 0, nil, nil
	})
	for sc.Scan() {

		text := sc.Text()

		println()
		println("BEGIN SENT")
		println(text)
		println("END SENT")
		println()

		// chomp the boilerplate at the end
		i := strings.Index(text, "End of the Project Gutenberg EBook")
		if i == -1 {
			extractEntities(text, entities)
		} else {
			extractEntities(text[:i], entities)
			break
		}
	}

	return entities, nil
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

func NewEntsFromPath(path string) (map[string]uint, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return NewEnts(f)
}

func NewEntsFromURL(url, path string) (map[string]uint, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return NewEnts(res.Body)
}

func scanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {

	sents := tokenizer.Tokenize(string(data))
	if len(sents) > 0 {
		first := sents[0]
		return len(first), []byte(first), nil
	}

	return 0, nil, nil
}
