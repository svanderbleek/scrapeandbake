package craigslist

import (
	"fmt"
	"github.com/rentapplication/craigjr/scraper"
	"net/url"
)

const (
	postsBase = "http://%v.craigslist.org/"
	postsPage = "http://%v.craigslist.org/search/apa?s=%v"
)

type Posts struct {
	City string
	*url.URL
	page int
	urls []string
	done bool
}

func NewPosts(city string) *Posts {
	base := fmt.Sprintf(postsBase, city)
	baseUrl, _ := url.Parse(base)
	return &Posts{
		City: city,
		URL:  baseUrl,
	}
}

func (p *Posts) Next() (string, bool) {
	offset := p.page * 100
	pageUrl := fmt.Sprintf(postsPage, p.City, offset)
	if !p.done {
		p.page++
	}
	return pageUrl, p.done
}

func (p *Posts) Items(body string) []string {
	hrefs := scrapeHrefs(body)
	if len(hrefs) > 0 {
		p.urls = p.absoluteUrls(hrefs)
	} else {
		p.done = true
		p.urls = []string{}
	}
	return p.urls
}

func scrapeHrefs(body string) []string {
	query := scraper.NewAttributeQuery(
		body,
		"#searchform .row .hdrlnk",
		"href",
	)
	return scraper.Scrape(query)
}

func (p *Posts) absoluteUrls(hrefs []string) []string {
	var urls []string
	for _, href := range hrefs {
		relative, _ := url.Parse(href)
		absolute := p.ResolveReference(relative)
		urls = append(urls, absolute.String())
	}
	return urls
}
