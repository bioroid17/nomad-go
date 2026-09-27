package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
}

// Scrape Saramin by a term
func Scrape(term string) {
	var baseUrl = "https://www.saramin.co.kr/zf_user/search/recruit?searchword=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseUrl)
	
	for i := 0; i < totalPages; i++ {
		go getPage(i, baseUrl, c)
	}
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}
func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	var jobSlices [][]string
	c := make(chan []string)
	for _, job := range jobs {
		go createJobSlice(job, c)
	}
	for i := 0; i < len(jobs); i++ {
		jobSlice := <-c
		jobSlices = append(jobSlices, jobSlice)
	}
	jwErr := w.WriteAll(jobSlices)
	checkErr(jwErr)
}
func createJobSlice(job extractedJob, c chan<- []string) {
	c <- []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?rec_idx=" + job.id, job.title, job.location, job.salary}
}
func getPage(page int, baseUrl string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
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
		go extractJob(card, c)
	})

	for i := 0; i< jobCards.Length(); i++ {
		job := <- c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".area_job > .job_tit > a").Text())
	location := cleanString(card.Find(".area_job > .job_condition > span:first-child").Text())
	salary := cleanString(card.Find(".area_job > .job_condition > span:nth-child(5)").Text())
	c <- extractedJob{
		id: id,
		title: title,
		location: location,
		salary: salary,
	}
}
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
func getPages(baseUrl string) int {
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