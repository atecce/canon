package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"atec.pub/canon/lib"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sentencesCmd = &cobra.Command{
	Use:   "sentences [dir]",
	Short: "segment sentences",
	Run: func(cmd *cobra.Command, args []string) {

		argc := len(args)

		if argc == 0 {

			sc := lib.NewSentenceScanner(os.Stdin)
			for sc.Scan() {
				println()
				println("BEGIN")
				os.Stdout.Write(sc.Bytes())
				os.Stdout.Write([]byte("\n"))
				println("END")
				println()
			}

		} else if argc == 1 {

			// TODO sort out err handling and logging

			sem := make(chan struct{}, 16)

			filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {

				if info.IsDir() {
					return nil
				}

				sem <- struct{}{}
				go func(path string, info os.FileInfo) {
					defer func() {
						<-sem
					}()

					logrus.Info(path)

					f, err := os.Open(path)
					if err != nil {
						logrus.Error(err)
						return
					}
					defer f.Close()

					var i uint
					sc := lib.NewSentenceScanner(f)
					for sc.Scan() {

						author, work := lib.SplitAuthorWork(info, path)

						b, err := json.Marshal(struct {
							Author string `json:"author"`
							Work   string `json:"work"`
							I      uint   `json:"i"`

							Text string `json:"text"`
						}{
							author,
							work,
							i,

							sc.Text(),
						})
						if err != nil {
							logrus.Error(err)
							continue
						}

						logrus.Info(string(b))

						req, err := http.NewRequest(http.MethodPut, "http://localhost:9200/sentences/_doc/"+url.QueryEscape(author+work+string(i)), bytes.NewReader(b))
						if err != nil {
							logrus.Error(err)
							continue
						}
						req.Header.Add("Content-Type", "application/json")

						res, err := http.DefaultClient.Do(req)
						if err != nil {
							logrus.Error(err)
							continue
						}

						logrus.Info(res.Status)

						b, err = ioutil.ReadAll(res.Body)
						if err != nil {
							logrus.Error(err)
							continue
						}

						logrus.Info(string(b))

						i++
					}

				}(path, info)

				return nil
			})

		} else {
			panic("too many args")
		}
	},
}

func init() {
	rootCmd.AddCommand(sentencesCmd)
}
