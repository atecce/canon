package fetch

import (
	"os"
	"path/filepath"

	"github.com/atecce/canon/fs"
	"github.com/atecce/canon/lib"
)

// FileFetcher hits https://gutenberg.org and writes the text into files in a directory
//
// fetching files is fast because it's parallelizable and has a low memory
// footprint because of the ability to simply pass res.Body to io.Copy. it
// also isolates failure well between each file
//
// however, it can create a mess at the destination
type FileFetcher struct {
	Root string
	Sem  chan struct{}
}

func (ff *FileFetcher) MkRoot() error {
	return fs.Mkdir(ff.Root)
}

func (ff *FileFetcher) MkAuthorDir(name string) error {
	return fs.Mkdir(filepath.Join(ef.Root, name))
}

func (ff *FileFetcher) Fetch(url, path string) error {

	ff.Sem <- struct{}{}

	go func() {
		defer func() {
			<-ff.Sem
		}()

		fullPath := filepath.Join(ff.Root, path)

		lib.Log(nil, url, fullPath, "INFO", "checking for path")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			lib.Log(nil, url, fullPath, "INFO", "not on fs. fetching")
			if err := fs.GetFile(url, fullPath+".txt.gz"); err != nil {
				lib.Log(nil, url, fullPath, "ERR", "fetching: "+err.Error())
			}
		}
	}()
	return nil
}
