package common

import (
	"path/filepath"
	"strings"
)

// Dir is the root directory where all the text data is stored
var Dir = filepath.Join("/", "keybase", "public", "atec", "data", "gutenberg")

// StripNewlines is used for filenames
func StripNewlines(str string) string {
	return strings.Replace(str, "\n", "", -1)
}
