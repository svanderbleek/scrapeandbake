package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func ScrapeAttributes(query Query) []string {
	document, err := NewDocument(query.Body)
	if err == nil {
		return document.Scrape(query)
	} else {
		log.Printf("Scraper error %v", err)
		return []string{}
	}
}

type Query struct {
	Body      string
	Selector  string
	Attribute string
}

type Document interface {
	Scrape(Query) []string
}

type GoQueryDocument struct {
	*goquery.Document
}

func NewDocument(html string) (Document, error) {
	reader := strings.NewReader(html)
	document, err := goquery.NewDocumentFromReader(reader)
	return &GoQueryDocument{document}, err
}

func (document GoQueryDocument) Scrape(query Query) []string {
	results := document.Find(query.Selector).Map(func(index int, selection *goquery.Selection) string {
		result, _ := selection.Attr(query.Attribute)
		return result
	})
	return results
}
