package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	// "github.com/PuerkitoBio/goquery"
)

type requestResult struct {
	url string
	status string
}

var errRequestFailed = errors.New("Request failed")
var urls []string

func main() {

	results := make(map[string]string)
	c := make(chan requestResult)

	// get stock code list
	stockList, err := getStockListByCSV()
	if err != nil {
		log.Fatal(err)
	}

	// url setting
	for _, stockCode := range stockList {
		urls = append(urls, "https://finance.naver.com/item/frgn.nhn?code="+stockCode)
	}

	// request
	for _, url := range urls {
		go hitURL(url, c)
	}

	// receive
	for i := 0; i< len(urls); i++ {
		result := <- c
		results[result.url] = result.status
	}

	// get data with goquery
	for url, status := range results {
		fmt.Println(url, status)
	}

}

func getStockListByCSV() (map[int]string, error) {

	listFile, err := os.Open("./test.csv")
	if err != nil {
		return nil, err
	}

	rdr := csv.NewReader(bufio.NewReader(listFile))
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, err
	}

	stockMap := make(map[int]string)

	for i, row := range rows {
		stockMap[i] = row[0]
	}

	return stockMap, nil

}

func hitURL(url string, c chan<- requestResult) {

	fmt.Println("Checking:", url)

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// checkErr(err)

	// doc.Find("외국인 기관 순매매 거래량에 관한표").Each(func(i int, s *goquery.Selection) {
	// 	// For each item found, get the band and title
	// 	band := s.Find("tah").Text()
	// 	fmt.Println(band)
	//   })

	// fmt.Println(doc)
	
	status := "success"
	c <- requestResult{url: url, status: status}

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}