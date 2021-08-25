package main

import (
	"Airport-Scarpper/Scrappe"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.POST("/flights", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, Scrappe.Scrappe(c.FormValue("airportCode"),c.FormValue("airportQueryType")))
	})
	e.POST("/airport", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, Scrappe.ScrappeAiportCode(c.FormValue("value")))
	})
	e.GET("/", func(c echo.Context) error {
		return c.File("home.html")
	})
	e.Logger.Fatal(e.Start(":8080"))
}