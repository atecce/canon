package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/lib"
)

func removeInvalidChars(str string) string {
	ret := str
	for _, char := range []string{" ", "\"", "\\", "<", "|", ",", ">", "/", "?"} {
		ret = strings.Replace(ret, char, "", -1)
	}
	return ret
}

type Entities struct {
	Values []lib.Entity `json:"values"`
}

func main() {

	// res, _ := http.Get("http://canon.atec.pub/_aliases")
	// b, _ := ioutil.ReadAll(res.Body)
	// var authors map[string]interface{}
	// json.Unmarshal(b, &authors)

	sem := make(chan struct{}, 10)

	filepath.Walk(".corpora/gutenberg", func(path string, info os.FileInfo, err error) error {

		// TODO try again on err?

		if strings.Contains(path, ".json.") {

			author := strings.ToLower(removeInvalidChars(filepath.Base(filepath.Dir(path))))
			// // TODO check per title and not just per author
			// if _, done := authors[author]; done {
			// 	log.Println("[INFO] author", author, "already", http.MethodPut)
			// 	return nil
			// }

			sem <- struct{}{}

			go func(author, title string) {

				defer func() {
					<-sem
				}()

				u := url.URL{
					Scheme: "http",
					Host:   "canon.atec.pub",
					Path:   filepath.Join(author, "entities", title),
				}

				f, err := os.Open(path)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				defer f.Close()

				r, _ := gzip.NewReader(f)

				var entities []lib.Entity
				json.NewDecoder(r).Decode(&entities)

				temp := Entities{
					Values: entities,
				}

				b, _ := json.Marshal(&temp)

				log.Println("[INFO]", http.MethodPut, u.String())
				req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewReader(b))
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				req.Header.Add("Content-Type", "application/json")
				// req.Header.Add("Content-Encoding", "gzip")

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				defer res.Body.Close()

				log.Println("[INFO]", res.Status)
				// TODO try again on error
				// if res.StatusCode != http.StatusCreated {
				// 	return
				// }

				b, err = ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				log.Println("[INFO]", string(b))

			}(author, removeInvalidChars(strings.Split(info.Name(), ".")[0]))

		}
		return nil
	})
}
