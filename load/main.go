package main

import (
	"compress/gzip"
	"io/ioutil"
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

		if strings.Contains(path, ".json.") {
			u := url.URL{
				Scheme: "http",
				Host:   "35.243.128.27",
				Path:   filepath.Join(strings.ToLower(removeInvalidChars(filepath.Base(filepath.Dir(path)))), "title", removeInvalidChars(strings.Split(info.Name(), ".")[0])),
			}
			f, _ := os.Open(path)
			defer f.Close()
			r, _ := gzip.NewReader(f)
			defer r.Close()
			println("[INFO]", u.String())
			req, _ := http.NewRequest("PUT", u.String(), r)
			req.Header.Add("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				// TODO try again?
				return nil
			}
			defer res.Body.Close()
			b, _ := ioutil.ReadAll(res.Body)
			println(string(b))
		}
		return nil
	})
}
