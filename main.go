package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var errRequestFailed = errors.New("Request failed")
var urls []string

func main() {

	stockList, err := getStockList()
	if err != nil {
		log.Fatal(err)
	}

	for _, stockCode := range stockList {
		urls = append(urls, "https://finance.naver.com/item/frgn.nhn?code="+stockCode)
	}

	for _, url := range urls {
		hitURL(url)
	}

}

func getStockList() (map[int]string, error) {

	listFile, err := os.Open("./stock_list.csv")
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

func hitURL(url string) error {

	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return errRequestFailed
	}

	return nil
}
