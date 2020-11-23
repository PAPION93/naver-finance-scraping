package scrapper

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
)

type requestResult struct {
	url   string
	datas []string
}

// Scrape func
func Scrape() {

	// var urls []string
	var results []string
	var requestResults []requestResult
	// c := make(chan requestResult)

	// get stock code list
	stockList, err := getStockListByCSV()
	checkErr(err)

	// url setting
	for i, stockCode := range stockList {
		url := "https://finance.naver.com/item/frgn.nhn?code=" + stockCode
		results = hitURI(url)
		requestResults = append(requestResults, requestResult{url: url, datas: results})
		fmt.Println(time.Now(), url, i)
	}

	// data processing
	processData := processing(requestResults)

	// write
	write(processData)
}

func processing(results []requestResult) [][]string {

	ps := [][]string{}
	for i := range results {

		tmp := [][]string{}

		for j, datas := range results[i].datas {

			s := []string{}

			// 매매 동향 없는 경우
			if len(datas) == 0 {
				break
			}

			dataSlices := strings.Fields(strings.TrimSpace(datas))

			// 가장 최근 날짜
			if j == 0 {

				// 종가 10,000 이하 제거
				closingPrice, _ := strconv.Atoi(strings.Replace(dataSlices[1], ",", "", 1))
				if closingPrice < 10000 {
					break
				}

				// 기관 동향이 0 또는 하락 종목 제거
				if strings.Contains(dataSlices[5], "-") || dataSlices[5] == "0" {
					break
				}

				// 기관동향 10,000 이상
				r := strings.NewReplacer(",", "", "+", "")
				price, _ := strconv.Atoi(r.Replace(dataSlices[5]))
				if price < 9000 {
					break
				}
			}

			// -1 day 하락일 경우 작성안함
			if j == 1 {
				if !strings.Contains(dataSlices[5], "-") {
					tmp = [][]string{}
					break
				}
			}

			s = append(s,
				results[i].url,
				dataSlices[0],
				dataSlices[1],
				dataSlices[2],
				dataSlices[3],
				dataSlices[4],
				dataSlices[5],
				dataSlices[6],
				dataSlices[7],
				dataSlices[8],
			)

			tmp = append(tmp, s)

		}

		for _, t := range tmp {
			ps = append(ps, t)
		}
	}

	return ps
}

func write(processData [][]string) {

	file, err := os.Create("./csv/기관매매기준.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "날짜", "종가", "전일비", "등락률", "거래량", "기관순 매매량", "외국인 순매매량", "외국인 보유주수", "외국인 보유율"}
	err = w.Write(headers)
	checkErr(err)

	for _, datas := range processData {

		err = w.Write(datas)
		checkErr(err)

	}
}

func getStockListByCSV() (map[int]string, error) {

	listFile, err := os.Open("./csv/stock_list.csv")
	checkErr(err)

	rdr := csv.NewReader(bufio.NewReader(listFile))
	rows, err := rdr.ReadAll()
	checkErr(err)

	stockMap := make(map[int]string)

	for i, row := range rows {
		stockMap[i] = row[0]
	}

	return stockMap, nil

} //end of getStockListByCSV()

func hitURI(url string) []string {

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	var datas []string
	doc.Find(".tc").EachWithBreak(func(i int, s *goquery.Selection) bool {

		if i >= 5 {
			return false
		}

		data, _ := s.Parent().Html()
		datas = append(datas, cleanString(data))

		return true
	})

	return datas

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
