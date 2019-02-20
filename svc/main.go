package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {

		fis, err := ioutil.ReadDir("/var/canon/gutenberg")
		if err != nil {
			return err
		}

		res := c.Response()

		res.WriteHeader(http.StatusOK)

		for _, fi := range fis {
			if _, err := res.Write([]byte(fi.Name() + "\n")); err != nil {
				return err
			}
			// time.Sleep(time.Minute)
		}

		return nil
	})

	e.GET("/:author", func(c echo.Context) error {

		author := c.Param("author")

		fis, err := ioutil.ReadDir("/var/canon/gutenberg/" + author)
		if err != nil {
			return err
		}

		var names []string
		for _, fi := range fis {
			names = append(names, fi.Name())
		}

		return c.JSON(http.StatusOK, names)
	})

	e.GET("/:author/:work", func(c echo.Context) error {

		author := c.Param("author")
		work := c.Param("work")

		b, err := ioutil.ReadFile("/var/canon/gutenberg/" + author + "/" + work)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, string(b))
	})

	e.StartTLS(":443", "/etc/canon/server.crt", "/etc/canon/server.key")
}
