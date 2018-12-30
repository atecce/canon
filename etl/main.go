package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/fs"
	"github.com/atecce/canon/lib"
)

const domain = "https://www.gutenberg.org/"

func main() {

	sem := make(chan struct{}, 16)

	filepath.Walk(".corpora/gutenberg", func(path string, info os.FileInfo, err error) error {

		if strings.Contains(path, ".txt.") {

			println(path)

			if _, err := os.Stat(path); !os.IsNotExist(err) {
				return nil
			}

			sem <- struct{}{}

			go func(textPath, jsonPath string) {
				defer func() {
					<-sem
				}()

				println(textPath)

				ents, err := lib.NewEntsFromPath(textPath)
				if err != nil {
					log.Fatal(err)
				}

				if err := fs.WriteJSON(jsonPath, &ents); err != nil {
					log.Fatal(err)
				}

			}(path, strings.Replace(path, ".txt.", ".json.", -1))
		}

		return nil
	})
}
