package craigslist

import (
	"github.com/rentapplication/craigjr/scraper"
	"github.com/rentapplication/craigjr/web"
	"regexp"
)

type Post struct {
	Url         string `json:"_id"`
	Body        string `json:"body"`
	ContactInfo string `json:"info,omitempty"`
}

func NewPost(url, body string) *Post {
	return &Post{
		Url:  url,
		Body: body,
	}
}

var showContactInfo = regexp.MustCompile(`show contact info`)

func (post *Post) IsContactInfo() bool {
	return showContactInfo.MatchString(post.Body)
}

func (post *Post) ContactInfoUrl() string {
	query := scraper.NewAttributeQuery(post.Body, ".showcontact", "href")
	href := scraper.Scrape(query)[0]
	return web.FullUrl(post.Url, href)
}
