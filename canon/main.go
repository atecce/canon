package main

import (
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/kr/pretty"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

func main() {
	u := url.URL{
		Scheme: "https",
		Host:   "www.gutenberg.org",
	}

	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		u.Path = filepath.Join("browse", "authors", string(letter))
		res, _ := http.Get(u.String())
		defer res.Body.Close()
		root, _ := html.Parse(res.Body)
		authors := scrape.FindAll(root, func(n *html.Node) bool {
			return n.Data == "h2" && n.FirstChild.FirstChild != nil
		})
		for _, author := range authors {
			pretty.Println(author.FirstChild.FirstChild.Data)
		}
	}
}
