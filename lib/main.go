package lib

import (
	"encoding/json"
	"time"
)

// Log writes json to stderr
func Log(size *int64, in, out, level, msg string) {
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
