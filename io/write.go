package io

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

func WriteJSON(path string, obj interface{}) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(obj); err != nil {
		return err
	}
	return nil
}

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
