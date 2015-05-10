package scraper

import (
	"github.com/PuerkitoBio/goquery"
)

type TextQuery struct {
	Query
}

func NewTextQuery(body string, selector string, attribute string) *AttributeQuery {
	return &AttributeQuery{
		Query: Query{
			body:     body,
			selector: selector,
		},
		Attribute: attribute,
	}
}

func (query *TextQuery) Result(selection *goquery.Selection) string {
	return ""
}
