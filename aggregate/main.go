package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/atecce/canon/lib"
	"github.com/kr/pretty"
	prose "gopkg.in/jdkato/prose.v2"
)

type Author struct {
	FirstName *string
	LastName  *string
	Life      [2]*time.Time
	Entities  []lib.Entity
}

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

			var docEntities []lib.Entity

			// print out accumulated entities
			if currentEntities != nil {
				println()

				for ent, count := range currentEntities {
					docEntities = append(docEntities, lib.Entity{
						Label: ent.Label,
						Text:  ent.Text,
						Count: count,
					})
				}
			}

			// reset author
			currentDir = name
			currentEntities = make(map[prose.Entity]uint)

			// split dir name for author metadata
			tmp := strings.Split(name, ",")

			println()

			if len(tmp) == 3 {

				life := strings.Split(tmp[2], "-")

				birthYear, _ := strconv.Atoi(life[0])
				deathYear, _ := strconv.Atoi(life[1])

				birth := time.Date(birthYear, time.January, 0, 0, 0, 0, 0, time.UTC)
				death := time.Date(deathYear, time.January, 0, 0, 0, 0, 0, time.UTC)

				author := Author{
					FirstName: &tmp[1],
					LastName:  &tmp[0],
					Life:      [2]*time.Time{&birth, &death},
					Entities:  docEntities,
				}

				pretty.Println(author)

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
