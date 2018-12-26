package fetch

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yhat/scrape"

	"github.com/gocolly/colly"

	"github.com/atecce/canon/lib"
)

// Files hits https://gutenberg.org and writes the text into files in a directory
//
// fetching files is fast because it's parallelizable and has a low memory
// footprint because of the ability to simply pass res.Body to io.Copy. it
// also isolates failure well between each file
//
// however, it can create a mess at the destination
func Files(root string) error {

	sem := make(chan struct{}, 10)

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		// remove pilcrows from author name
		author := filepath.Join(root, strings.Replace(e.ChildText("a"), "¶", "", -1))

		if _, err := os.Stat(author); os.IsNotExist(err) {
			if mkErr := os.Mkdir(author, 0700); mkErr != nil {
				lib.Log(nil, author, "", "ERR", "failed to mkdir: "+err.Error())
			}
		}

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				sem <- struct{}{}

				// TODO try again on err?
				go func(href, title string) {
					defer func() {
						<-sem
					}()

					// remove forward slashes and new lines
					name := lib.RemoveNewlines(strings.Replace(title, "/", "|", -1))

					url := domain + href + ".txt.utf-8"
					if strings.Contains(url, "wikipedia") {
						return
					}

					path := filepath.Join(author, name+".txt.gz")
					lib.Log(nil, url, path, "INFO", "checking for kbPath")
					if _, err := os.Stat(path); os.IsNotExist(err) {
						lib.Log(nil, url, path, "INFO", "not on kbfs. fetching")
						if err := fetchFile(url, path); err != nil {
							lib.Log(nil, url, path, "ERR", "fetching: "+err.Error())
						}
					}
				}(scrape.Attr(node.FirstChild, "href"), node.FirstChild.FirstChild.Data)
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}

	return nil
}

func fetchFile(url, path string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}

	return nil
}
