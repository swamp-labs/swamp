package httpreq

import (
	"io"
	"log"
	"net/http"
	"net/http/httptrace"

	as "github.com/swamp-labs/swamp/engine/assertion"
)

// Request struct defines all parameters for http requests to execute
type Request struct {
	Name            string              `yaml:"name"`
	Verb            string              `yaml:"verb"`
	Protocol        string              `yaml:"protocol"`
	Headers         []map[string]string `yaml:"headers"`
	URL             string              `yaml:"url"`
	Body            string              `yaml:"body"`
	QueryParameters map[string]string   `yaml:"query_parameters"`
	Assertions      as.Assertion        `yaml:"assertions"`
}

// AddQueryParam func is called to add query parameters to the http request
func (r *Request) AddQueryParam(key string, value string) *Request {
	r.QueryParameters[key] = value
	return r
}

// Execute send the http request using client trace
// and call AssertResponse function to validate the response
func (r *Request) Execute(m map[string]string) (bool, error) {
	client := &http.Client{}
	clientTrace := Trace{}
	req, _ := http.NewRequest(r.Verb, r.URL, nil)
	q := req.URL.Query()
	if len(r.QueryParameters) > 0 {
		for key, value := range r.QueryParameters {
			q.Add(key, value)
		}
	}

	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace.Trace())
	req = req.WithContext(clientTraceCtx)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	s := clientTrace.Done()
	v, err := as.AssertResponse(r.Assertions, resp, m)
	if err != nil {
		return false, err
	}
	r.displayResult(resp, m, v, s)

	return v, nil
}

func (r *Request) displayResult(resp *http.Response, m map[string]string, b bool, s *Sample) {
	body, _ := io.ReadAll(resp.Body)
	log.Println("#######################################")
	log.Println("--------------- Request ---------------")
	log.Println("Request name :", r.Name)
	log.Println("Target :", r.Verb, r.URL)
	log.Println("---------------------------------------")
	log.Println("--------------- Response --------------")
	log.Println("Response code :", resp.StatusCode)
	log.Println("Response body :", string(body))
	log.Println("---------------------------------------")
	log.Println("---------- Assertions Result ----------")
	log.Println("Global validation :", b)
	log.Println("Returned variables :", m)
	log.Println("---------------------------------------")
	log.Println("------------ Request traces -----------")
	s.displayTrace()
	log.Println("#######################################")
}
