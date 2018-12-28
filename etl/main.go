package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"

	"github.com/atecce/canon/lib"
)

const domain = "https://www.gutenberg.org/"

func writeJSON(doc *lib.Doc, path string) error {
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

func main() {

	filepath.Walk("gutenberg", func(path string, info os.FileInfo, err error) error {

		if strings.Contains(path, ".txt.") {
			println(path)
			doc, _ := lib.NewDocFromPath(path)
			pretty.Println(doc)
		}

		return nil
	})
}

// func temp() {
// 	workers, err := strconv.Atoi(os.Getenv("WORKERS"))
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to read in worker configuration: %v\n", err)
// 		fmt.Fprintln(os.Stderr, "defaulting to 16")
// 		workers = 16
// 	}

// 	sem := make(chan struct{}, workers)

// 	authorCollector := colly.NewCollector()

// 	authorCollector.OnRequest(func(r *colly.Request) {
// 		lib.Log(r.URL.Path, "", "INFO", r.Method)
// 	})

// 	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

// 		// remove pilcrows from author name
// 		author := filepath.Join("gutenberg", strings.Replace(e.ChildText("a"), "¶", "", -1))

// 		if _, err := os.Stat(author); os.IsNotExist(err) {
// 			if mkErr := os.Mkdir(author, 0700); mkErr != nil {
// 				lib.Log(author, "", "ERR", "failed to mkdir: "+err.Error())
// 			}
// 		}

// 		for _, node := range e.DOM.Next().Children().Nodes {
// 			if node.FirstChild.FirstChild != nil {

// 				sem <- struct{}{}

// 				// TODO try again on err?
// 				go func(href, title string) {
// 					defer func() {
// 						<-sem
// 					}()

// 					// remove forward slashes and new lines
// 					name := lib.RemoveNewlines(strings.Replace(title, "/", "|", -1))

// 					url := domain + href + ".txt.utf-8"
// 					if strings.Contains(url, "wikipedia") {
// 						return
// 					}

// 					textPath := filepath.Join(author, name+".txt.gz")
// 					lib.Log(url, textPath, "INFO", "checking for path")
// 					if _, err := os.Stat(textPath); os.IsNotExist(err) {

// 						lib.Log(url, textPath, "INFO", "not on fs. creating new doc")
// 						doc, err := lib.NewDoc(url, textPath)
// 						if err != nil {
// 							lib.Log(url, textPath, "ERR", "creating new doc: "+err.Error())
// 						}

// 						jsonPath := strings.Replace(textPath, ".txt.", ".json.", -1)
// 						lib.Log(url, jsonPath, "INFO", "writing")
// 						if err := writeJSON(doc, jsonPath); err != nil {
// 							lib.Log(url, jsonPath, "ERR", "writing: "+err.Error())
// 							return
// 						}
// 					}
// 				}(scrape.Attr(node.FirstChild, "href"), node.FirstChild.FirstChild.Data)
// 			}
// 		}
// 	})

// 	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
// 		authorCollector.Visit(domain + "browse/authors/" + string(letter))
// 	}
// }
