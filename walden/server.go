package main

import (
	"bufio"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

func readline(r io.Reader, k int) ([]byte, error) {
	sc := bufio.NewScanner(r)
	i := 0
	for sc.Scan() {
		i++
		if i == k {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Bytes(), sc.Err()
		}
	}
	return nil, io.EOF
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		f, _ := os.Open("/home/atec/sentences.txt")

		sentence, _ := readline(f, rand.Intn(57))
		w.WriteHeader(http.StatusOK)
		w.Write(sentence)
	})

	http.ListenAndServe(":80", nil)

}
