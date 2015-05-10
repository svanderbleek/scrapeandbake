package scraper

type Query struct {
	body     string
	selector string
}

func (query *Query) Body() string {
	return query.body
}

func (query *Query) Selector() string {
	return query.selector
}
