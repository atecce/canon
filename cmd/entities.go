package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"atec.pub/canon/lib"
	"atec.pub/io"

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
			sem := make(chan struct{}, 16)

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if strings.Contains(path, ".txt") {

					println(path)

					sem <- struct{}{}

					go func(textPath, jsonPath string) {
						defer func() {
							<-sem
						}()

						println(jsonPath)

						ents, err := lib.NewEntsFromPath(textPath)
						if err != nil {
							log.Fatal(err)
						}

						if err := io.WriteJSON(jsonPath, &ents); err != nil {
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
