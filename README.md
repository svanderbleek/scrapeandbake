## craigjr

Packages: `craigslist`, `crawler`, `scraper`, `proxy`

```
go build && ./craigjr
```

see `demo.go` for usage

## Config

```
export IN_CLOAK_API_KEY=
export PROX_IS_RIGHT_API_KEY=
export KING_PROXY_API_KEY=
```

## Agenda

- [x] Proxy List as Http Transport
- [x] Posts scrapes full Urls for each Post
- [ ] Pull from City List for Crawlers to produce Post Urls
- [ ] Go routines to pull from UrlQueue and scrape Post data
