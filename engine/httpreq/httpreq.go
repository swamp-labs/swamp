package httpreq

import (
	"log"
	"net/http"
	"net/http/httptrace"

	as "github.com/swamp-labs/swamp/engine/assertion"
)

// type QueryParameter struct {
// 	Key   string `yaml:"key"`
// 	Value string `yaml:"value"`
// }

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

// Execute generate send the http request using client trace
func (r *Request) Execute() (map[string]string, error) {
	client := &http.Client{}
	clientTrace := Trace{}
	req, _ := http.NewRequest(r.Verb, r.URL, nil)

	if r.Headers != nil {
		for i := range r.Headers {
			for key := range r.Headers[i] {
				log.Println("key", key, "value", r.Headers[i][key])
				req.Header.Set(key, r.Headers[i][key])
			}
		}
	}
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
		return nil, err
	}
	defer resp.Body.Close()
	clientTrace.Done()
	_, variables, err := as.AssertResponse(r.Assertions, resp)
	if err != nil {
		return nil, err
	}
	log.Println(resp.Header)

	return variables, nil
}
