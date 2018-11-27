package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filepath.Walk("/keybase/public/atec/data/gutenberg/", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "\n") {
			newPath := strings.Replace(path, "\n", "", -1)
			log.Println("[INFO] renaming", path, "to", newPath)
			err := os.Rename(path, newPath)
			if err != nil {
				log.Println("[ERR]", err)
			}
		}
		return nil
	})
}
