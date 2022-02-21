package httpreq

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Method          string
	URL             string
	QueryParameters map[string]string
	Assertions      []assertion
	SaveValues      []saveValue
	Name            string
}

type assertion struct{}

type saveValue struct{}

func (r *Request) InitQueryParam() *Request {
	r.QueryParameters = make(map[string]string)
	return r
}

func (r *Request) AddQueryParam(key string, value string) *Request {
	r.QueryParameters[key] = value
	return r
}

func (r *Request) Execute() error {
	client := &http.Client{}

	req, _ := http.NewRequest(r.Method, r.URL, nil)

	q := req.URL.Query()
	for key, value := range r.QueryParameters {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	log.Println(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(resp.Status)
	log.Println(string(responseBody))

	return nil
}
