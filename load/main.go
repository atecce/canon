package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/atecce/canon/common"
)

func removeInvalidChars(str string) string {
	ret := str
	for _, char := range []string{" ", "\"", "\\", "<", "|", ",", ">", "/", "?"} {
		ret = strings.Replace(ret, char, "", -1)
	}
	return ret
}

func main() {

	res, _ := http.Get("http://canon.atec.pub/_aliases")
	b, _ := ioutil.ReadAll(res.Body)
	var aliases map[string]interface{}
	json.Unmarshal(b, &aliases)

	var last string
	for author := range aliases {
		if author > last {
			last = author
		}
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	filepath.Walk(common.Dir, func(path string, info os.FileInfo, err error) error {

		// TODO try again on err?

		if strings.Contains(path, ".json.") {

			author := strings.ToLower(removeInvalidChars(filepath.Base(filepath.Dir(path))))
			if author < last {
				return nil
			}

			wg.Add(1)
			sem <- struct{}{}

			go func(author, title string) {

				defer func() {
					wg.Done()
					<-sem
				}()

				u := url.URL{
					Scheme: "http",
					Host:   "canon.atec.pub",
					Path:   filepath.Join(author, "title", title),
				}

				f, err := os.Open(path)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				defer f.Close()

				r, err := gzip.NewReader(f)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				defer r.Close()

				log.Println("[INFO]", u.String())
				req, err := http.NewRequest("PUT", u.String(), r)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				req.Header.Add("Content-Type", "application/json")

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}
				defer res.Body.Close()

				b, err := ioutil.ReadAll(res.Body)
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
