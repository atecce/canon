package fetch

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/fs"

	"github.com/atecce/canon/lib"
	"github.com/gocolly/colly"
	"github.com/yhat/scrape"
)

type Fetcher interface {
	GetAuthor() string
	GetTitle() string
	Join(string, string) string
	Fetch(name string) error
	crawl()
}

func crawl(fetcher Fetcher) {
	authorCollector.OnHTML("h2", func(e *colly.HTMLElement) {

		author := fetcher.GetAuthor()

		for _, node := range e.DOM.Next().Children().Nodes {
			child := node.FirstChild
			grandchild := child.FirstChild
			if grandchild != nil {

				url := domain + scrape.Attr(child, "href") + ".txt.utf-8"
				if strings.Contains(url, "wikipedia") {
					return
				}

				title := grandchild.Data

				path := fetcher.Join(author, title)
			}
		}
	})

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		authorCollector.Visit(domain + "browse/authors/" + string(letter))
	}
}

type Tarballer struct {
	tw *tar.Writer

	url  string
	path string
}

// Tarball hits https://gutenberg.org and writes the text directly into a tarball
//
// tarball produces a single clean artifact which is easily moved around with file
// operations
//
// however, it is not parallelizable or resilient to failure. when it exits you will
// mostly likely need to start it from scratch. in addition, because you need to
// write file sizes in tar headers, it can create a considerable footprint counting
// bytes in memory
func Tarball(name string) error {

	f, err := os.Create(name)
	if err != nil {
		return err
	}
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

				if err := fs.GetTarFile(url, path, tw); err != nil {
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
