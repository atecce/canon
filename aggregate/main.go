package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/atecce/canon/fs"

	"github.com/atecce/canon/lib"
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
		currentEntities map[string]uint
	)

	corporaEntities := make(map[string]uint)

	filepath.Walk(".corpora/gutenberg/", func(path string, info os.FileInfo, err error) error {

		name := info.Name()

		// skip base dir
		if name == "gutenberg" {
			return nil
		}

		// new author
		if info.IsDir() && name != currentDir {

			var authorEntities []lib.Entity
			if currentEntities != nil {
				for text, count := range currentEntities {

					authorEntities = append(authorEntities, lib.Entity{
						Text:  text,
						Count: count,
					})

					if corporaCount, ok := corporaEntities[text]; ok {
						corporaEntities[text] = count + corporaCount
					} else {
						corporaEntities[text] = count
					}
				}
			}

			// reset author
			currentDir = name
			currentEntities = make(map[string]uint)

			// // split dir name for author metadata
			// tmp := strings.Split(name, ",")

			// println()

			// if len(tmp) == 3 {

			// 	life := strings.Split(tmp[2], "-")

			// 	birthYear, _ := strconv.Atoi(life[0])
			// 	deathYear, _ := strconv.Atoi(life[1])

			// 	birth := time.Date(birthYear, time.January, 0, 0, 0, 0, 0, time.UTC)
			// 	death := time.Date(deathYear, time.January, 0, 0, 0, 0, 0, time.UTC)

			// 	author := Author{
			// 		FirstName: &tmp[1],
			// 		LastName:  &tmp[0],
			// 		Life:      [2]*time.Time{&birth, &death},
			// 		Entities:  authorEntities,
			// 	}

			// 	pretty.Println(author)

			// } else {
			// 	pretty.Println("TODO", tmp)
			// }
			// println()

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

				var entities []lib.Entity
				if err = json.NewDecoder(r).Decode(&entities); err != nil {
					fmt.Println(err)
					return nil
				}

				// aggregate entities from the authors work
				for _, entity := range entities {
					if count, ok := currentEntities[entity.Text]; ok {
						currentEntities[entity.Text] = count + entity.Count
					} else {
						currentEntities[entity.Text] = entity.Count
					}
				}

				f.Close()
				r.Close()
			}
		}

		return nil
	})

	var totalEntities []lib.Entity
	for text, count := range corporaEntities {
		totalEntities = append(totalEntities, lib.Entity{
			Text:  text,
			Count: count,
		})
	}

	sort.Slice(totalEntities, func(i, j int) bool {
		return totalEntities[i].Count > totalEntities[j].Count
	})

	fs.WriteJSON("entities.json.gz", &totalEntities)
}
