package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/lib"
	"github.com/kr/pretty"
	prose "gopkg.in/jdkato/prose.v2"
)

func main() {

	var currentDir string
	var currentEntities map[prose.Entity]uint

	filepath.Walk(lib.Dir, func(path string, info os.FileInfo, err error) error {

		name := info.Name()

		if name == "gutenberg" {
			return nil
		}

		if info.IsDir() && name != currentDir {

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

			currentDir = name
			currentEntities = make(map[prose.Entity]uint)

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
			if strings.Contains(path, ".json.") {
				println(path)

				f, _ := os.Open(path)

				r, _ := gzip.NewReader(f)

				var doc lib.Doc
				json.NewDecoder(r).Decode(&doc)

				// pretty.Println(doc.Entities)

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
