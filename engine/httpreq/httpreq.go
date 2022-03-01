package httpreq

import (
	"log"
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

// func (r *Request) InitQueryParam() *Request {
// 	r.QueryParameters = make(map[string]string)
// 	return r
// }

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
	if err != nil {
		return err
	}

	sample := clientTrace.Done()
	log.Println("######### SAMPLE #########")
	log.Println("EndTime: ", sample.EndTime)
	log.Println("ConnDuration: ", sample.ConnDuration)
	log.Println("ReqDuration: ", sample.ReqDuration)
	log.Println("WaitingConn: ", sample.WaitingConn)
	log.Println("Connecting: ", sample.Connecting)
	log.Println("TLSHandshaking: ", sample.TLSHandshaking)
	log.Println("Sending: ", sample.Sending)
	log.Println("WaitingResp: ", sample.WaitingResp)
	log.Println("Receiving: ", sample.Receiving)
	log.Println("ConnReused: ", sample.ConnReused)
	log.Println("ConnRemoteAddr: ", sample.ConnRemoteAddr)

	return nil
}
