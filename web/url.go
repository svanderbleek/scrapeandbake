package web

import (
	"net/url"
)

func FullUrl(baseUrl string, relativeUrl string) string {
	base, _ := url.Parse(baseUrl)
	relative, _ := url.Parse(relativeUrl)
	return base.ResolveReference(relative).String()
}
