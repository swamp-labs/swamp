package httpreq

import (
	"net/http"
	"net/http/httptrace"
)

// Request defines an http request
type Request struct {
	Method          string
	URL             string
	QueryParameters map[string]string
	Assertions      []assertion
	SaveValues      []saveValue
	Name            string
}

type assertion struct{} //TO DO

type saveValue struct{} //TO DO

// AddQueryParam func is called to add query parameters to the http request
func (r *Request) AddQueryParam(key string, value string) *Request {
	r.QueryParameters[key] = value
	return r
}

// Execute generate send the http request using client trace
func (r *Request) Execute() error {
	client := &http.Client{}
	clientTrace := Trace{}
	req, _ := http.NewRequest(r.Method, r.URL, nil)

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
		return err
	}
	defer resp.Body.Close()

	clientTrace.Done()

	return nil
}
