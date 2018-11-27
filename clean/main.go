package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/common"
)

func main() {
	filepath.Walk(common.Dir, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "\n") {
			newPath := common.StripNewlines(path)
			log.Println("[INFO] renaming", path, "to", newPath)
			err := os.Rename(path, newPath)
			if err != nil {
				log.Println("[ERR]", err)
			}
		}
		return nil
	})
}
