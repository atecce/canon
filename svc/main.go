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

var (
	authors      []string
	authorsLower []string

	lowerMap = make(map[string]string)
)

func init() {
	fis, err := ioutil.ReadDir("/var/gutenberg/")
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {

		author := fi.Name()
		authorLower := strings.ToLower(author)

		authors = append(authors, author)
		authorsLower = append(authorsLower, authorLower)
		lowerMap[authorLower] = author
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"/authors?search={pattern}":      "fuzzy search for an author",
			"/authors/{author}":              "get an author by name",
			"/authors/{author}/works/{work}": "get an authors work by name",
		})
	})

	e.GET("/authors", func(c echo.Context) error {

		pattern := c.QueryParam("search")

		res := c.Response()

		res.WriteHeader(http.StatusOK)

		for _, authorLower := range fuzzy.Find(strings.ToLower(pattern), authorsLower) {
			if _, err := res.Write([]byte(lowerMap[authorLower] + "\n")); err != nil {
				return err
			}
		}

		return nil
	})

	e.GET("/authors/:author", func(c echo.Context) error {

		author := c.Param("author")

		fis, err := ioutil.ReadDir("/var/gutenberg/" + author)
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

	e.GET("/authors/:author/works/:work", func(c echo.Context) error {

		author := c.Param("author")
		work := c.Param("work")

		b, err := ioutil.ReadFile("/var/gutenberg/" + author + "/" + work + ".json")
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, string(b))
	})

	e.StartTLS(":443", "/etc/srv.crt", "/etc/srv.key")
}
