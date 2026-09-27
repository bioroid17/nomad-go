package main

import (
	"nomad-go/scrapper"
	"strings"

	"github.com/labstack/echo/v5"
)

func handleHome(c *echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c *echo.Context) error {
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	return nil
}


func main() {
	e := echo.New()
	
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}