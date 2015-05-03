package scraper

import (
	"golang.org/x/net/html"
	"net/http"
)

type Filter struct {
	Attributes []string
}

func Urls(url string, filter Filter) []string {
	links := links(url)
	links = applyFilter(links, filter)
	return urls(links)
}

func links(url string) []*html.Token {
	tokenizer := tokenizer(url)
	var links []*html.Token
	link := nextLinkToken(tokenizer)
	for link != nil {
		links = append(links, link)
		link = nextLinkToken(tokenizer)
	}
	return links
}

func applyFilter(links []*html.Token, filter Filter) []*html.Token {
	var passed []*html.Token
	for _, link := range links {
		if pass(link, filter) {
			passed = append(links, link)
		}
	}
	return passed
}

func pass(link *html.Token, filter Filter) bool {
	for _, attribute := range filter.Attributes {
		_, found := getAttribute(link, attribute)
		if !found {
			return false
		}
	}
	return true
}

func urls(links []*html.Token) []string {
	var urls []string
	for _, link := range links {
		urls = append(urls, href(link))
	}
	return urls
}

func href(link *html.Token) string {
	href, _ := getAttribute(link, "href")
	return href
}

func getAttribute(token *html.Token, key string) (string, bool) {
	for _, attribute := range token.Attr {
		if attribute.Key == key {
			return attribute.Val, true
		}
	}
	return "", false
}

func tokenizer(url string) *html.Tokenizer {
	response, _ := http.Get(url)
	return html.NewTokenizer(response.Body)
}

func nextLinkToken(tokenizer *html.Tokenizer) *html.Token {
	for {
		tokenType := tokenizer.Next()
		if isErrorToken(tokenType) {
			return nil
		}
		token := tokenizer.Token()
		if isLinkToken(token) {
			return &token
		}
	}
}

func isErrorToken(tokenType html.TokenType) bool {
	return tokenType == html.ErrorToken
}

func isLinkToken(token html.Token) bool {
	return token.Type == html.StartTagToken && token.DataAtom.String() == "a"
}
