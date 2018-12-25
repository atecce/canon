package lib

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

// Dir is the root directory where all the text data is stored
var Dir = filepath.Join("/", "keybase", "public", "atec", "data", "gutenberg")

// RemoveNewlines is used for filenames
func RemoveNewlines(str string) string {
	return strings.Replace(str, "\n", "", -1)
}

// Log writes json to stderr
func Log(in, out, level, msg string) {
	b, _ := json.Marshal(struct {
		Time  time.Time `json:"time"`
		In    string    `json:"in"`
		Out   string    `json:"out"`
		Level string    `json:"level"`
		Msg   string    `json:"msg"`
	}{
		time.Now(),
		in,
		out,
		level,
		msg,
	})
	println(string(b))
}
