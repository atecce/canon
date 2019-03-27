package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atec.pub/canon/lib"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirupsen/logrus"
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

			// TODO checkpoint

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if strings.Contains(path, ".txt") {

					sem <- struct{}{}

					go func(path string, info os.FileInfo) {
						defer func() {
							<-sem
						}()

						// TODO maybe dedup author and work info from fileinfo
						author := filepath.Base(filepath.Dir(path))
						work := strings.TrimSuffix(info.Name(), ".txt")

						ents, err := lib.NewEntsFromPath(path)
						if err != nil {
							logrus.WithFields(logrus.Fields{
								"author": author,
								"work":   work,
							}).Error(err)
							return
						}

						entities := struct {
							ID       string `bson:"_id"`
							Author   string
							Work     string
							Entities map[string]uint
						}{
							author + work,
							author,
							work,
							ents,
						}

						logrus.WithFields(logrus.Fields{
							"author": author,
							"work":   work,
						}).Info("extracting entities")

						_, err = collection.InsertOne(context.TODO(), entities)
						if err != nil {
							logrus.WithFields(logrus.Fields{
								"author": author,
								"work":   work,
							}).Error(err)
							return
						}

					}(path, info)
				}

				return nil
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(entitiesCmd)
}
