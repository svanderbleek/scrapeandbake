package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

type Querier interface {
	Body() string
	Selector() string
	Result(*goquery.Selection) string
}

func Scrape(query Querier) []string {
	body := query.Body()
	document, err := NewDocument(body)
	if err == nil {
		return document.Scrape(query)
	} else {
		log.Printf("Scraper error %v", err)
		return []string{}
	}
}

type Document interface {
	Scrape(Querier) []string
}

type GoQueryDocument struct {
	*goquery.Document
}

func NewDocument(html string) (Document, error) {
	reader := strings.NewReader(html)
	document, err := goquery.NewDocumentFromReader(reader)
	return &GoQueryDocument{document}, err
}

func (document GoQueryDocument) Scrape(query Querier) []string {
	selector := query.Selector()
	selections := document.Find(selector)
	results := selections.Map(func(index int, selection *goquery.Selection) string {
		return query.Result(selection)
	})
	return results
}
