## craigjr

Packages: `craigslist`, `crawler`, `visitor`, `scraper`, `proxy`, `store`

## config

```
export IN_CLOAK_API_KEY=
export MASHAPE_API_KEY=
export KING_PROXY_API_KEY=
export ELASTIC_SEARCH_URL=
```

## example

Crawler:

```go
import (
	"github.com/rentapplication/craigjr/craigslist"
	"github.com/rentapplication/craigjr/crawler"
	"github.com/rentapplication/craigjr/proxy"
)

proxy := proxy.NewList()
proxy.LoadDefault()

posts := make(chan crawler.PaginationIterator)
go craigslist.StreamCities(posts)

crawlers := crawler.NewPool(proxy)
go crawlers.Crawl(posts)
```

Visitor:

```go
import (
	"github.com/rentapplication/craigjr/proxy"
	"github.com/rentapplication/craigjr/store"
)

proxy := proxy.NewList()
proxy.LoadDefault()

visitors := visitor.NewPool(proxy)
go visitors.Visit(crawlers.Urls)
```

## agenda

- [x] Proxy List as Http Transport
- [x] Posts scrapes full Urls for each Post
- [x] Pull from City List for Crawlers to produce Post Urls
- [x] Visitors to pull from Url stream
- [x] Index Posts into elastic search
- [x] More Proxy Sources
- [ ] Balance pools based on limited resources (automatically preferred)
