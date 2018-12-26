package fetch

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yhat/scrape"

	"github.com/gocolly/colly"

	"github.com/atecce/canon/lib"
)

const domain = "https://www.gutenberg.org/"

var authorCollector *colly.Collector

func init() {
	authorCollector = colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		lib.Log(nil, r.URL.Path, "", "INFO", r.Method)
	})
}

// Files hits https://gutenberg.org and writes the text into files in a directory
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

func write(url, path string, tw *tar.Writer) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	size := int64(len(b))

	if err := tw.WriteHeader(&tar.Header{
		Name: path,
		Size: size,
		Mode: 0444,
	}); err != nil {
		return err
	}

	lib.Log(&size, url, path, "INFO", "writing")
	if _, err := tw.Write(b); err != nil {
		return err
	}

	return nil
}

// Tarball hits https://gutenberg.org and writes the text directly into a tarball
func Tarball(name string) error {

	f, _ := os.Create(name)
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		// remove pilcrows from author name
		author := strings.Replace(e.ChildText("a"), "¶", "", -1)

		for _, node := range e.DOM.Next().Children().Nodes {
			if node.FirstChild.FirstChild != nil {

				// remove forward slashes and new lines
				name := lib.RemoveNewlines(strings.Replace(node.FirstChild.FirstChild.Data, "/", "|", -1))

				url := domain + scrape.Attr(node.FirstChild, "href") + ".txt.utf-8"
				if strings.Contains(url, "wikipedia") {
					return
				}

				path := filepath.Join(author, name+".txt")

				if err := write(url, path, tw); err != nil {
					lib.Log(nil, url, path, "ERR", "writing: "+err.Error())
				}
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
