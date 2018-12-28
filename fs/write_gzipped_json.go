package fs

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

func WriteGzippedJSON(path string, obj interface{}) error {
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
