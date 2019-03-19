package lib

import (
	"bufio"
	"io"
	"regexp"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var segmenter *sentences.DefaultSentenceTokenizer

func init() {
	// TODO support languages not English
	segmenter, _ = english.NewSentenceTokenizer(nil)
}

var newlines = regexp.MustCompile(`\r?\n`)

func NewSentenceScanner(r io.Reader) *bufio.Scanner {

	sc := bufio.NewScanner(r)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		sents := segmenter.Tokenize(string(data))
		if len(sents) > 0 {
			first := sents[0].Text
			return len(first), newlines.ReplaceAll([]byte(first), []byte(" ")), nil
		}

		return 0, nil, nil
	})

	return sc
}
