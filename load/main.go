package main

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/common"
)

type author struct {
	Titles []string `json:"titles"`
}

func removeInvalidChars(str string) string {
	ret := str
	for _, char := range []string{" ", "\"", "\\", "<", "|", ",", ">", "/", "?"} {
		ret = strings.Replace(ret, char, "", -1)
	}
	return ret
}

func main() {
	filepath.Walk(common.Dir, func(path string, info os.FileInfo, err error) error {

		// TODO try again on err?

		if strings.Contains(path, ".json.") {

			author := strings.ToLower(removeInvalidChars(filepath.Base(filepath.Dir(path))))
			title := removeInvalidChars(strings.Split(info.Name(), ".")[0])

			u := url.URL{
				Scheme: "http",
				Host:   "35.243.128.27",
				Path:   filepath.Join(author, "title", title),
			}

			f, err := os.Open(path)
			if err != nil {
				log.Println("[ERR]", err)
				return nil
			}
			defer f.Close()

			r, err := gzip.NewReader(f)
			if err != nil {
				log.Println("[ERR]", err)
				return nil
			}
			defer r.Close()

			log.Println("[INFO]", u.String())
			req, err := http.NewRequest("PUT", u.String(), r)
			if err != nil {
				log.Println("[ERR]", err)
				return nil
			}
			req.Header.Add("Content-Type", "application/json")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("[ERR]", err)
				return nil
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("[ERR]", err)
				return nil
			}
			log.Println("[INFO]", string(b))
		}
		return nil
	})
}
