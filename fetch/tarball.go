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

	"github.com/atecce/canon/lib"
	"github.com/gocolly/colly"
	"github.com/yhat/scrape"
)

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

				if _, err := write(url, path, tw); err != nil {
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

func write(url, path string, tw *tar.Writer) (int64, error) {
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var size int64
	if res.ContentLength == -1 {

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return 0, err
		}

		size = int64(len(b))

		if err := tw.WriteHeader(&tar.Header{
			Name: path,
			Size: size,
			Mode: 0444,
		}); err != nil {
			return 0, err
		}

		lib.Log(&size, url, path, "INFO", "writing")
		if _, err := tw.Write(b); err != nil {
			return 0, err
		}

	} else {

		size = res.ContentLength

		if err := tw.WriteHeader(&tar.Header{
			Name: path,
			Size: size,
			Mode: 0444,
		}); err != nil {
			return 0, err
		}

		lib.Log(&size, url, path, "INFO", "writing")
		if n, err := io.Copy(tw, res.Body); err != nil {
			return n, err
		}
	}

	return size, nil
}
