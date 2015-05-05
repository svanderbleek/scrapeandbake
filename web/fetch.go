package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Query struct {
	Url    string
	Result interface{}
	Error  error
}

func (q *Query) Fetch() {
	body, err := fetchBody(q.Url)
	if err == nil {
		err = readResult(body, q.Result)
	}
	if err != nil {
		q.Error = err
	}
}

func fetchBody(url string) ([]byte, error) {
	var body []byte
	response, err := http.Get(url)
	if err == nil {
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
	}
	return body, err
}

func readResult(body []byte, result interface{}) error {
	err := json.Unmarshal(body, result)
	return err
}
