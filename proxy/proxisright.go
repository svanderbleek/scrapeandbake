package proxy

import (
	"encoding/json"
	"fmt"
)

const (
	apiList = "https://theproxisright.com/api/proxy/get?onlyActive=true&onlySupportsCraigslist=true&apiKey=%v"
)

type apiResult struct {
	List `json:"list"`
}

func ProxIsRight() List {
	return List{}
}

func listUrl() {
	return fmt.Sprintf(apiList, os.Getenv("PROX_API_KEY"))
}

func fetchList() apiResult {
}
