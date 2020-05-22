package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
)

type requestResult struct {
	url    string
	status string
	datas  []string
}

var errRequestFailed = errors.New("Request failed")
var urls []string

func main() {

	var results []requestResult
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
		go hitURI(url, c)
	}

	// receive
	for i := 0; i < len(urls); i++ {
		result := <-c
		results := append(results, result)
		fmt.Println(results)
	}

	// get data with goquery
	// for result := range results {
	// 	fmt.Println(result)
	// }

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

} //end of getStockListByCSV()

func hitURI(url string, c chan<- requestResult) {

	var datas []string

	fmt.Println("Checking:", url)

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".tc").Each(func(i int, s *goquery.Selection) {
		data, _ := s.Parent().Html()
		datas = append(datas, cleanString(data))
	})

	status := "success"
	c <- requestResult{url: url, status: status, datas: datas}

} //end of hitURI()

func cleanString(str string) string {
	stripped := strip.StripTags(str)
	return strings.Join(strings.Fields(strings.TrimSpace(stripped)), " ")
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
