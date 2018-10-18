package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	prose "gopkg.in/jdkato/prose.v2"
)

const (
	dir = "/keybase/private/atec/data/walden"

	url = "https://www.gutenberg.org/files/205/205-h/205-h.htm"

	walden = dir + "/index.html.gz"
)

func fetch() {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("getting")
		fmt.Println(err)
		os.Exit(1)
	}
	b := res.Body
	defer b.Close()

	f := createFile(walden)
	defer f.Close()

	if n, err := io.Copy(f, b); err != nil {
		fmt.Println("copying")
		fmt.Println(err)
		fmt.Println(n)
		os.Exit(1)
	}
}

func createFile(name string) *gzip.Writer {
	f, err := os.Create(name)
	if err != nil {
		fmt.Println("creating file")
		fmt.Println(err)
		os.Exit(1)
	}
	return gzip.NewWriter(f)
}

func openFile(url string) *gzip.Reader {
	f, err := os.Open(url)
	if err != nil {
		fmt.Println("opening")
		fmt.Println(err)
		os.Exit(1)
	}
	gz, _ := gzip.NewReader(f)
	return gz
}

type Category uint

const (
	Nature Category = iota
)

type text struct {
	category Category
	content string
}

func main() {

	// res, err := http.Get("http://localhost:9200")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// body := res.Body
	// defer body.Close()
	// b, _ := ioutil.ReadAll(body)
	// fmt.Println(string(b))
	// os.Exit(0)

	// TODO gross
	// get walden corpora
	fmt.Println("checking for file:", walden)
	if _, err := os.Stat(walden); err != nil {
		fmt.Println("can't find file, getting web page:", url)
		fetch()
	}

	r := openFile(walden)
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("reading all")
		fmt.Println(err)
		os.Exit(1)
	}
	corpus := string(b)

	doc, err := prose.NewDocument(corpus)
	if err != nil {
		fmt.Println("creating new prose doc")
		fmt.Println(err)
		os.Exit(1)
	}

	// sentences := dir + "/nature/sentences.txt.gz"
	// fmt.Println("creating file:", sentences)

	// w := createFile(sentences)
	// defer w.Close()

	regex := regexp.MustCompile(`\s+`)
	rdr, w := io.Pipe()
	httpBody := json.NewDecoder(rdr).Buffered()
	defer w.Close()

	sentences := make([]string, 57)

	// TODO filter out html tags?
	for _, sent := range doc.Sentences() {
		if strings.Contains(sent.Text, "Nature") {
			sentence := regex.ReplaceAllString(sent.Text, " ")
			sentences = append(sentences, sentence)
			fmt.Println(sentence)
			n, err := w.Write([]byte(sentence + "\n"))
			if err != nil {
				fmt.Println("writing sentence")
				fmt.Println(err)
				fmt.Println(n)
				os.Exit(1)
			}
		}
	}

	res, err := http.Post("http://localhost:9200/thoreau/walden/nature", "", httpBody)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	output, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(output))
}
