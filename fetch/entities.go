package fetch

import (
	"os"
	"path/filepath"

	"github.com/atecce/canon/fs"

	"github.com/atecce/canon/lib"
)

type EntitiesFetcher struct {
	Root string
	Sem  chan struct{}
}

func (ef *EntitiesFetcher) MkRoot() error {
	if _, err := os.Stat(ef.Root); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(ef.Root, 0700); mkErr != nil {
			lib.Log(nil, ef.Root, "", "ERR", "failed to mkdir: "+err.Error())
		}
	}
	return nil
}

func (ef *EntitiesFetcher) MkAuthorDir(name string) error {
	if _, err := os.Stat(filepath.Join(ef.Root, name)); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(filepath.Join(ef.Root, name), 0700); mkErr != nil {
			lib.Log(nil, name, "", "ERR", "failed to mkdir: "+err.Error())
		}
	}
	return nil
}

func (ef *EntitiesFetcher) Fetch(url, path string) error {

	ef.Sem <- struct{}{}

	// TODO try again on err?
	go func() {
		defer func() {
			<-ef.Sem
		}()

		fullPath := filepath.Join(ef.Root, path) + ".json.gz"

		lib.Log(nil, url, fullPath, "INFO", "checking for path")
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {

			lib.Log(nil, url, fullPath, "INFO", "not on fs. getting ents from url")
			ents, err := lib.NewEntsFromURL(url, fullPath)
			if err != nil {
				lib.Log(nil, url, fullPath, "ERR", "getting ents: "+err.Error())
			}

			lib.Log(nil, url, fullPath, "INFO", "writing")
			if err := fs.WriteJSON(fullPath, ents); err != nil {
				lib.Log(nil, url, fullPath, "ERR", "writing: "+err.Error())
				return
			}
		}
	}()

	return nil
}
