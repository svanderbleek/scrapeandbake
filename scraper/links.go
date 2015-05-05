package scraper

import (
	"../proxy"
)

type Query struct {
	Selector  string
	Attribute string
}

type Document interface {
	Filter(Query) []string
}

func Urls(reponse proxy.Response, query Query) []string {
	document := NewDocument(response.body)
	document.Filter(query)
}
