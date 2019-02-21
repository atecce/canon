package fetch

import (
	"os"
	"path/filepath"

	"github.com/atecce/canon/io"
	"github.com/sirupsen/logrus"
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
	Ext  string
}

func (ff *FileFetcher) MkRoot() error {
	return io.Mkdir(ff.Root)
}

func (ff *FileFetcher) MkAuthorDir(name string) error {
	return io.Mkdir(filepath.Join(ff.Root, name))
}

func (ff *FileFetcher) Fetch(url, path string) error {

	ff.Sem <- struct{}{}

	go func() {
		defer func() {
			<-ff.Sem
		}()

		fullPath := filepath.Join(ff.Root, path)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {

			var getErr error

			logrus.WithFields(logrus.Fields{
				"url":  url,
				"path": fullPath,
			}).Info("getting file")
			switch ff.Ext {
			case ".txt":
				getErr = io.GetFile(url, fullPath+ff.Ext)
			case ".txt.gz":
				getErr = io.GetGzippedFile(url, fullPath+ff.Ext)
			default:
				println("invalid extension")
				os.Exit(1)
			}
			if getErr != nil {
				logrus.WithFields(logrus.Fields{
					"url":  url,
					"path": fullPath,
				}).Error("fetching:", err)
			}
		}
	}()
	return nil
}
