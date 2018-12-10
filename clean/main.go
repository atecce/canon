package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kr/pretty"

	"github.com/atecce/canon/lib"
)

func main() {
	infos, _ := ioutil.ReadDir(lib.Dir)
	for _, info := range infos {
		path := filepath.Join(lib.Dir, info.Name())
		if strings.Contains(path, "¶") {
			newPath := strings.Replace(path, "¶", "", -1)
			log.Println("[INFO] renaming", path, "to", newPath)
			err := os.Rename(path, newPath)
			if err != nil {
				linkErr := err.(*os.LinkError)
				pretty.Logln("[ERR]", linkErr)
			}
		}
	}
}
