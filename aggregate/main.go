package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/lib"
	"github.com/kr/pretty"
	prose "gopkg.in/jdkato/prose.v2"
)

func main() {

	var (
		currentDir      string
		currentEntities map[prose.Entity]uint
	)

	filepath.Walk(lib.Dir, func(path string, info os.FileInfo, err error) error {

		name := info.Name()

		// skip base dir
		if name == "gutenberg" {
			return nil
		}

		// new author
		if info.IsDir() && name != currentDir {

			// print out accumulated entities
			if currentEntities != nil {
				println()

				var docEntities []lib.Entity
				for ent, count := range currentEntities {
					docEntities = append(docEntities, lib.Entity{
						Label: ent.Label,
						Text:  ent.Text,
						Count: count,
					})
				}

				pretty.Println(docEntities)
			}

			// reset author
			currentDir = name
			currentEntities = make(map[prose.Entity]uint)

			// split dir name for author metadata
			tmp := strings.Split(name, ",")

			println()

			if len(tmp) == 3 {
				println("last name:", tmp[0])
				println("first name:", tmp[1])

				println("life:", tmp[2])
			} else {
				pretty.Println("TODO", tmp)
			}
			println()

		} else {

			// only get extracted docs
			if strings.Contains(path, ".json.") {
				println(path)

				// decode doc
				f, err := os.Open(path)
				if err != nil {
					fmt.Println(err)
					return nil
				}

				r, err := gzip.NewReader(f)
				if err != nil {
					fmt.Println(err)
					return nil
				}

				var doc lib.Doc
				if err = json.NewDecoder(r).Decode(&doc); err != nil {
					fmt.Println(err)
					return nil
				}

				// aggregate entities from the authors work
				for _, docEnt := range doc.Entities {
					proseEnt := prose.Entity{
						Label: docEnt.Label,
						Text:  docEnt.Text,
					}

					if count, ok := currentEntities[proseEnt]; ok {
						currentEntities[proseEnt] = count + docEnt.Count
					} else {
						currentEntities[proseEnt] = docEnt.Count
					}
				}
			}
		}

		return nil
	})
}
