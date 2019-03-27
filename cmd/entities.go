package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atec.pub/canon/lib"

	"go.mongodb.org/mongo-driver/bson"
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

			ctx := context.TODO()

			sem := make(chan struct{}, 16)

			// set up mongo
			client, _ := mongo.Connect(ctx, options.Client().ApplyURI("localhost:27017"))

			collection := client.Database("canon").Collection("entities")

			// check for checkpoint
			cur, err := collection.Find(ctx, bson.D{}, &options.FindOptions{
				Sort: map[string]int{"_id": -1},
			})
			if err != nil {
				logrus.Error(err)
				os.Exit(1)
			}
			defer cur.Close(ctx)

			var start bool
			var checkpoint string
			var i uint
			for cur.Next(ctx) {
				var res bson.M
				cur.Decode(&res)

				i++
				if i == 16 {
					checkpoint = res["_id"].(string)
					break
				}
			}

			// walk
			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if strings.Contains(path, ".txt") {

					author, work := lib.SplitAuthorWork(info, path)

					if author+work == checkpoint {
						start = true
					}

					if start {

						sem <- struct{}{}

						go func(author, work string) {
							defer func() {
								<-sem
							}()

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

						}(author, work)
					}
				}

				return nil
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(entitiesCmd)
}
