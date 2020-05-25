package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
)

type requestResult struct {
	url   string
	datas []string
}

var errRequestFailed = errors.New("Request failed")
var urls []string
var stockFile = "stock_list.csv"

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

	// select
	for _, url := range urls {
		go hitURI(url, c)
	}

	// receive
	for i := 0; i < len(urls); i++ {
		result := <-c
		results = append(results, result)
	}

	writeResultDataToCSV(results)

}

func writeResultDataToCSV(results []requestResult) {

	file, err := os.Create("results.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "날짜", "종가", "전일비", "등락률", "거래량", "기관순 매매량", "외국인 순매매량", "외국인 보유주수", "외국인 보유율"}

	err = w.Write(headers)
	checkErr(err)

	for i := range results {
	L1:
		for j, datas := range results[i].datas {

			var bodies []string
			bodies = append(bodies, results[i].url)

			// 매매 동향 없는 경우
			if len(datas) == 0 {
				break
			}

			dataSlices := strings.Fields(strings.TrimSpace(datas))
			for k := range dataSlices {

				// 종가 10,000 이하 제거
				if j == 0 && k == 1 {
					price, _ := strconv.Atoi(strings.Replace(dataSlices[k], ",", "", 1))
					if price < 10000 {
						break L1
					}
				}

				// 기관 동향이 0 또는 하락 종목 제거
				if j == 0 && k == 5 && (strings.Contains(dataSlices[k], "-") || dataSlices[k] == "0") {
					break L1
				}
				bodies = append(bodies, dataSlices[k])
			}
			err = w.Write(bodies)
			checkErr(err)
		}
	}
}

func getStockListByCSV() (map[int]string, error) {

	listFile, err := os.Open(stockFile)
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

	doc.Find(".tc").EachWithBreak(func(i int, s *goquery.Selection) bool {

		if i >= 5 {
			return false
		}

		data, _ := s.Parent().Html()
		datas = append(datas, cleanString(data))

		return true
	})

	c <- requestResult{url: url, datas: datas}

} //end of hitURI()

func cleanString(str string) string {
	cleanString := strip.StripTags(str)
	return strings.Join(strings.Fields(strings.TrimSpace(cleanString)), " ")
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
