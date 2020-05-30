package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/PAPION93/naver-finance-scraping/scrapper"
	"github.com/labstack/echo"
)

type data struct {
	url   string
	datas []string
}

// headers := []string{"Link", "날짜", "종가", "전일비", "등락률", "거래량", "기관순 매매량", "외국인 순매매량", "외국인 보유주수", "외국인 보유율"}

const fileName string = "기관매매기준.csv"

func handleScrapeCSV(c echo.Context) error {
	// defer os.Remove(fileName)
	scrapper.Scrape()
	return c.JSON(http.StatusOK, "")
}

func handleScrapeJSON(c echo.Context) error {
	defer os.Remove(fileName)
	scrapper.Scrape()

	listFile, err := os.Open(fileName)
	checkErr(err)
	rdr := csv.NewReader(bufio.NewReader(listFile))
	rows, err := rdr.ReadAll()
	checkErr(err)

	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/scrape/csv", handleScrapeCSV)
	e.GET("/scrape/json", handleScrapeJSON)
	e.Logger.Fatal(e.Start(":1323"))
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
