package fs

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/atecce/canon/lib"
)

func Mkdir(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(name, 0700); mkErr != nil {
			return mkErr
		}
	}
	return nil
}

func WriteJSON(path string, obj interface{}) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}
	return nil
}

func GetFile(url, path string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return err
	}
	return nil
}

func GetGzippedFile(url, path string) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}
	return nil
}

func GetTarFile(url, path string, tw *tar.Writer) error {

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.ContentLength == -1 {

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		size := int64(len(b))

		if err := tw.WriteHeader(&tar.Header{
			Name: path,
			Size: size,
			Mode: 0444,
		}); err != nil {
			return err
		}

		lib.Log(&size, url, path, "INFO", "writing")
		if _, err := tw.Write(b); err != nil {
			return err
		}

	} else {

		size := res.ContentLength

		if err := tw.WriteHeader(&tar.Header{
			Name: path,
			Size: size,
			Mode: 0444,
		}); err != nil {
			return err
		}

		lib.Log(&size, url, path, "INFO", "writing")
		if _, err := io.Copy(tw, res.Body); err != nil {
			return err
		}
	}
	return nil
}
