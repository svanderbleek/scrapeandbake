package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Query struct {
	Url    string
	Header *Header
	Result interface{}
	Error  error
}

type Header struct {
	Key   string
	Value string
}

func (q *Query) Fetch() {
	body, err := fetchBody(q)
	if err == nil {
		q.Result = string(body)
	} else {
		q.Error = err
	}
}

func (q *Query) FetchJson() {
	body, err := fetchBody(q)
	if err == nil {
		err = readResult(body, q.Result)
	}
	if err != nil {
		q.Error = err
	}
}

func fetchBody(q *Query) ([]byte, error) {
	var body []byte
	response, err := fetch(q)
	if err == nil {
		defer response.Body.Close()
		body, err = ioutil.ReadAll(response.Body)
	}
	return body, err
}

func fetch(q *Query) (*http.Response, error) {
	if q.Header != nil {
		client := &http.Client{}
		request, _ := http.NewRequest("GET", q.Url, nil)
		request.Header.Add(q.Header.Key, q.Header.Value)
		return client.Do(request)
	} else {
		return http.Get(q.Url)
	}
}

func readResult(body []byte, result interface{}) error {
	err := json.Unmarshal(body, result)
	return err
}
