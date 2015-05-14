## craigjr

Packages: `craigslist`, `crawler`, `visitor`, `scraper`, `proxy`

```
go build && ./craigjr
```

see `demo.go` for usage

## config

```
export IN_CLOAK_API_KEY=
export PROX_IS_RIGHT_API_KEY=
export KING_PROXY_API_KEY=
export ELASTIC_SEARCH_URL=
```

## agenda

- [x] Proxy List as Http Transport
- [x] Posts scrapes full Urls for each Post
- [x] Pull from City List for Crawlers to produce Post Urls
- [x] Visitors to pull from Url stream
- [ ] Index Posts into elastic search
