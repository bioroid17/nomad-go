package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseUrl = "https://www.saramin.co.kr/zf_user/search/recruit?searchword=python"

func main() {
	totalPages := getPages()
	
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}
func getPage(page int) {
	pageUrl := baseUrl + "&recruitPage=" + strconv.Itoa(page + 1)
	fmt.Println("Requesting", pageUrl)
}
func getPages() int {
	pages := 0
	res, err := http.Get(baseUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
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