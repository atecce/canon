package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	// TODO pool of goroutines on a channel
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	filepath.Walk(lib.Dir, func(textPath string, info os.FileInfo, err error) error {

		// TODO try again on err?
		lib.Log(0, textPath, "", "INFO", "walking")
		if err != nil {
			lib.Log(0, textPath, "", "ERR", "walking")
			return nil
		}

		if !strings.Contains(textPath, ".txt") {
			return nil
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(jsonPath string) {
			defer func() {
				wg.Done()
				<-sem
			}()

			lib.Log(0, textPath, jsonPath, "INFO", "checking kbfs")
			if _, err := os.Stat(jsonPath); os.IsNotExist(err) {

				size := info.Size()

				lib.Log(size, textPath, jsonPath, "INFO", "reading")
				doc, err := lib.NewDoc(textPath)
				if err != nil {
					lib.Log(size, textPath, jsonPath, "ERR", "creating doc: "+err.Error())
					return
				}

				lib.Log(size, textPath, jsonPath, "INFO", "writing")
				if err := writeJSON(doc, jsonPath); err != nil {
					lib.Log(size, textPath, jsonPath, "ERR", "writing: "+err.Error())
					return
				}
			}
		}(strings.Replace(textPath, ".txt.", ".json.", -1))

		return nil
	})
	wg.Wait()
}
