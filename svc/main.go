package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ctx context.Context

	collection *mongo.Collection

	// kept in memory at all times for perf on /authors?search
	lowerAuthors   []string
	lowerAuthorMap = make(map[string]string)
)

func init() {

	ctx = context.TODO()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("canon").Collection("entities")

	authorDocs, err := collection.Distinct(ctx, "author", bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for _, doc := range authorDocs {

		author := doc.(string)
		lowerAuthor := strings.ToLower(author)

		lowerAuthorMap[lowerAuthor] = author
		lowerAuthors = append(lowerAuthors, lowerAuthor)
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

		res.WriteHeader(http.StatusOK)

		var matches []string
		for _, lowerAuthor := range lowerAuthors {
			if strings.Contains(lowerAuthor, strings.ToLower(pattern)) {
				matches = append(matches, lowerAuthor)
			}
		}

		for _, match := range matches {
			if _, err := res.Write([]byte(lowerAuthorMap[match] + "\n")); err != nil {
				return err
			}
		}

		return nil
	})

	e.GET("/authors/:author", func(c echo.Context) error {

		author := c.Param("author")

		cur, err := collection.Find(ctx, bson.D{
			{
				"author", author,
			},
		}, &options.FindOptions{
			Projection: map[string]bool{
				"_id":  false,
				"work": true,
			},
		})
		if err != nil {
			return err
		}

		var names []string
		for cur.Next(ctx) {

			var work bson.D
			cur.Decode(&work)

			names = append(names, work[0].Value.(string))
		}

		return c.JSON(http.StatusOK, names)
	})

	e.GET("/authors/:author/works/:work", func(c echo.Context) error {

		author := c.Param("author")
		work := c.Param("work")

		res := collection.FindOne(ctx, bson.D{
			{
				"_id", author + work,
			},
		}, &options.FindOneOptions{
			Projection: map[string]bool{
				"_id":      false,
				"entities": true,
			},
		})

		var doc bson.D
		if err := res.Decode(&doc); err != nil {
			return err
		}

		ents := make(map[string]int64)
		for _, ent := range doc[0].Value.(bson.D) {
			ents[ent.Key] = ent.Value.(int64)
		}

		return c.JSON(http.StatusOK, ents)
	})

	e.StartTLS(":443", "/etc/canon/server.crt", "/etc/canon/server.key")
}
