package craigslist

import (
	"../scraper"
)

const (
	posts = "http://%v.craigslist.org/search/apa?s=%v"
)

type Posts struct {
	City string
	page int
	urls []string
}

func (p *Posts) Next() bool {
	offset := p.page * 100
	pageUrl := fmt.Sprintf(posts, p.City, offset)
	p.urls = postUrls(pageUrl)
	p.page++
	return len(p.urls) > 0
}

func (p *Posts) Items() []string {
	return p.urls
}

func postUrls(page string) []string {
	postFilter := scraper.Filter{Attributes: []string{"data-id", "href", "class"}}
	return scraper.Urls(page, postFilter)
}
