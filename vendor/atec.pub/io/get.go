package io

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

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

		if _, err := io.Copy(tw, res.Body); err != nil {
			return err
		}
	}
	return nil
}
