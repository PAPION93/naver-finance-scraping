package main

import (
	"net/http"

	"github.com/PAPION93/naver-finance-scraping/scrapper"
	"github.com/labstack/echo"
)

const fileName string = "./csv/기관매매기준.csv"

func handleScrapeCSV(c echo.Context) error {
	return c.Attachment(fileName, fileName)
}

func handleScrapeJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, "Success")
}

func handleSaveCSV(c echo.Context) error {
	scrapper.Scrape()
	return c.JSON(http.StatusOK, "Success")
}

func main() {
	e := echo.New()
	e.GET("/scrape/csv", handleScrapeCSV)
	e.GET("/scrape/json", handleScrapeJSON)
	e.GET("/scrape", handleSaveCSV)

	e.Logger.Fatal(e.Start(":8080"))
}
