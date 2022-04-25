package assertion

import (
	"io/ioutil"
	"net/http"
)

// Assertion defines checks to execute against an http response
type Assertion struct {
	Body    []BodyAssertion
	Code    []int
	Headers []map[string][]string
}

// AssertResponse executes all assertions defined by user
// it calls validateCodeStatus, validateHeaders and validateBody
// to verify each kind of assertions
func AssertResponse(a Assertion, resp *http.Response, m map[string]string) (valid bool, err error) {
	validBody := true
	validCode := a.validateCodeStatus(resp.StatusCode)
	validHeaders := a.validateHeaders(&resp.Header)
	if resp.ContentLength > 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		if a.Body != nil {
			validBody, err = a.validateBody(body, m)
			if err != nil {
				return false, err
			}
		}
	}
	validResp := validCode && validBody && validHeaders
	return validResp, nil
}
