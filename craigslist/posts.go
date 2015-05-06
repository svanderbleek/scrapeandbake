package craigslist

import (
	"fmt"
	"github.com/rentapplication/craigjr/proxy"
	"github.com/rentapplication/craigjr/scraper"
	"net/url"
)

const (
	postsPage = "http://%v.craigslist.org/search/apa?s=%v"
)

type Posts struct {
	City string
	page int
	urls []string
}

func (p *Posts) Next() bool {
	offset := p.page * 100
	pageUrl := fmt.Sprintf(postsPage, p.City, offset)
	p.urls = postUrls(pageUrl)
	p.page++
	return len(p.urls) > 0
}

func (p *Posts) Items() []string {
	return p.urls
}

func postUrls(url string) []string {
	response := proxy.MustGet(url)
	hrefs := scrapeHrefs(response.Body)
	return absoluteUrls(url, hrefs)
}

func scrapeHrefs(body string) []string {
	query := scraper.Query{
		Body:      body,
		Selector:  "#searchform .row .hdrlnk",
		Attribute: "href",
	}
	return scraper.ScrapeAttributes(query)
}

func absoluteUrls(baseUrl string, hrefs []string) []string {
	var urls []string
	base, _ := url.Parse(baseUrl)
	for _, href := range hrefs {
		relative, _ := url.Parse(href)
		absolute := base.ResolveReference(relative).String()
		urls = append(urls, absolute)
	}
	return urls
}
