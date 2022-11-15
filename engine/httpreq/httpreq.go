package httpreq

import (
	as "github.com/swamp-labs/swamp/engine/assertion"
	ts "github.com/swamp-labs/swamp/engine/templateString"
	"net/http"
	"net/http/httptrace"
)

// Request struct defines all parameters for http requests to execute
type Request struct {
	Name            string              `yaml:"name"`
	Method          string              `yaml:"method"`
	Protocol        string              `yaml:"protocol"`
	Headers         []map[string]string `yaml:"headers"`
	URL             ts.TemplateString   `yaml:"url"`
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

	url, _ := r.URL.ToString(m)
	req, _ := http.NewRequest(r.Method, url, nil)
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
	v, err := r.Assertions.AssertResponse(resp, m)
	if err != nil {
		return false, err
	}
	r.displayResult(resp, m, v, s)

	return v, nil
}
