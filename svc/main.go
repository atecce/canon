package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	authors      []string
	authorsLower []string

	lowerMap = make(map[string]string)

	ctx context.Context

	collection *mongo.Collection
)

func init() {

	ctx = context.TODO()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("canon").Collection("entities")

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
			"/authors?search={pattern}":      "search for an author",
			"/authors/{author}":              "get an author by name",
			"/authors/{author}/works/{work}": "get an authors work by name",
		})
	})

	e.GET("/authors", func(c echo.Context) error {

		res := c.Response()
		pattern := c.QueryParam("search")

		// TODO attempt to filter on query

		items, err := collection.Distinct(ctx, "author", bson.D{})
		if err != nil {
			return err
		}

		for _, item := range items {
			itemStr := item.(string)
			if strings.Contains(strings.ToLower(itemStr), strings.ToLower(pattern)) {
				res.Write([]byte(itemStr + "\n"))
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

	e.StartTLS(":443", "/etc/canon/server.crt", "/etc/canon/server.key")
}
