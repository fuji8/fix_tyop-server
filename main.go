package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Golang, Hello, World!")
	})
	e.GET("/hogehoge", saveAndServeSound)
	e.Logger.Fatal(e.Start(":1323"))
}
