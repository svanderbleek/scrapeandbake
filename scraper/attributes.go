package scraper

import (
	"github.com/PuerkitoBio/goquery"
)

type AttributeQuery struct {
	Query
	Attribute string
}

func NewAttributeQuery(body string, selector string, attribute string) *AttributeQuery {
	return &AttributeQuery{
		Query: Query{
			body:     body,
			selector: selector,
		},
		Attribute: attribute,
	}
}

func (query *AttributeQuery) Result(selection *goquery.Selection) string {
	result, _ := selection.Attr(query.Attribute)
	return result
}
