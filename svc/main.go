package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/usr/local/var/canon")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
