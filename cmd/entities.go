package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"atec.pub/canon/lib"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var entitiesCmd = &cobra.Command{
	Use:   "entities [dir]",
	Short: "extract entities",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			ents, err := lib.NewEnts(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to construct entities: %v\n", err)
				os.Exit(1)
			}
			if err := json.NewEncoder(os.Stdout).Encode(&ents); err != nil {
				fmt.Fprintf(os.Stderr, "failed to encode json: %v\n", err)
				os.Exit(1)
			}
		} else {

			// set up mongo
			client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("localhost:27017"))

			collection := client.Database("canon").Collection("entities")

			sem := make(chan struct{}, 16)

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if strings.Contains(path, ".txt") {

					// TODO maybe dedup author and work info from fileinfo
					author := filepath.Base(filepath.Dir(path))
					work := strings.TrimSuffix(info.Name(), ".txt")

					sem <- struct{}{}

					go func(textPath, jsonPath string) {
						defer func() {
							<-sem
						}()

						ents, err := lib.NewEntsFromPath(textPath)
						if err != nil {
							log.Fatal(err)
						}

						pretty.Println(ents)

						entities := struct {
							Author   string
							Work     string
							Entities map[string]uint
						}{
							author,
							work,
							ents,
						}

						pretty.Println(entities)

						_, err = collection.InsertOne(context.TODO(), entities)
						if err != nil {
							log.Fatal(err)
						}

					}(path, strings.Replace(path, ".txt", ".json", -1))
				}

				return nil
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(entitiesCmd)
}
