package assertion

import (
	"io/ioutil"
	"net/http"
)

type Assertion interface {
	AssertResponse(resp *http.Response, m map[string]string) (bool, error)
}

func MakeRequestAssertion(body []BodyAssertion, code []int, headers []map[string][]string) Assertion {
	return &assertion{
		Body:    body,
		Code:    code,
		Headers: headers,
	}
}

// assertion defines checks to execute against an http response
type assertion struct {
	Body    []BodyAssertion
	Code    []int
	Headers []map[string][]string
}

// AssertResponse executes all assertions defined by user
// it calls validateCodeStatus, validateHeaders and validateBody
// to verify each kind of assertions
func (a *assertion) AssertResponse(resp *http.Response, m map[string]string) (valid bool, err error) {
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
