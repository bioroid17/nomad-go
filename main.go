package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
}

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
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobCards := doc.Find(".item_recruit")
	jobCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("value")
		fmt.Println(id)
		title := card.Find(".area_job > .job_tit > a").Text()
		fmt.Println(title)
		location := card.Find(".area_job > .job_condition > span:first-child").Text()
		fmt.Println(location)
		salary := card.Find(".area_job > .job_condition > span:nth-child(5)").Text()
		fmt.Println(salary)
	})
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