package io

import "os"

func Mkdir(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(name, 0700); mkErr != nil {
			return mkErr
		}
	}
	return nil
}
