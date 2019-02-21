package fetch

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/atecce/canon/io"

	"github.com/atecce/canon/lib"
)

type EntitiesFetcher struct {
	Root string
	Sem  chan struct{}
	Ext  string
}

func (ef *EntitiesFetcher) MkRoot() error {
	return io.Mkdir(ef.Root)
}

func (ef *EntitiesFetcher) MkAuthorDir(name string) error {
	return io.Mkdir(filepath.Join(ef.Root, name))
}

func (ef *EntitiesFetcher) Fetch(url, path string) error {

	ef.Sem <- struct{}{}

	// TODO try again on err?
	go func() {
		defer func() {
			<-ef.Sem
		}()

		fullPath := filepath.Join(ef.Root, path) + ef.Ext

		logrus.WithFields(logrus.Fields{
			"path": fullPath,
		}).Info("checking")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {

			logrus.WithFields(logrus.Fields{
				"url":  url,
				"path": fullPath,
			}).Info("not on fs. getting ents")
			ents, err := lib.NewEntsFromURL(url, fullPath)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"url":  url,
					"path": fullPath,
				}).Error("getting ents:", err)
			}

			logrus.WithFields(logrus.Fields{
				"url":  url,
				"path": fullPath,
			}).Info("writing json")

			var writeErr error
			switch ef.Ext {
			case ".json":
				writeErr = io.WriteJSON(fullPath, ents)
			case ".json.gz":
				writeErr = io.WriteGzippedJSON(fullPath, ents)
			default:
				println("invalid extension")
				os.Exit(1)
			}
			if writeErr != nil {
				logrus.WithFields(logrus.Fields{
					"url":  url,
					"path": fullPath,
				}).Error("writing json:", err)
				return
			}
		}
	}()

	return nil
}
