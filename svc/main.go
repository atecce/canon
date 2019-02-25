package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var authors []string

func init() {
	fis, err := ioutil.ReadDir("/var/canon/gutenberg/")
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {
		author := fi.Name()
		authors = append(authors, strings.TrimSuffix(author, filepath.Ext(author)))
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.GET("/authors", func(c echo.Context) error {

		pattern := c.QueryParam("search")

		res := c.Response()

		res.WriteHeader(http.StatusOK)

		for _, author := range fuzzy.Find(pattern, authors) {
			if _, err := res.Write([]byte(author + "\n")); err != nil {
				return err
			}
		}

		return nil
	})

	e.GET("authors/:author", func(c echo.Context) error {

		author := c.Param("author")

		fis, err := ioutil.ReadDir("/var/canon/gutenberg/" + author)
		if err != nil {
			return err
		}

		var names []string
		for _, fi := range fis {
			name := fi.Name()
			names = append(names, strings.TrimSuffix(name, filepath.Ext(name)))
		}

		return c.JSON(http.StatusOK, names)
	})

	e.GET("authors/:author/works/:work", func(c echo.Context) error {

		author := c.Param("author")
		work := c.Param("work")

		b, err := ioutil.ReadFile("/var/canon/gutenberg/" + author + "/" + work + ".json")
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, string(b))
	})

	e.StartTLS(":443", "/etc/canon/server.crt", "/etc/canon/server.key")
}
