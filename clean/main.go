package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atecce/canon/lib"
)

func main() {

	tarball, err := os.Create(filepath.Join(lib.Dir, "text.tar.gz"))
	if err != nil {
		log.Fatal(err)
	}

	gzw := gzip.NewWriter(tarball)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	if err := filepath.Walk(lib.Dir, func(path string, info os.FileInfo, err error) error {

		if strings.Contains(path, ".txt.") && info.Mode().IsRegular() {

			log.Println("tarring", path)

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(strings.Replace(path, lib.Dir, "", -1), string(filepath.Separator))

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		log.Fatal(err)
	}
}
