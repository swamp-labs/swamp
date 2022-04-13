package assertion

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Assertion defines checks to execute against an http response
type Assertion struct {
	Body    []bodyAssertion
	Code    []int
	Headers []map[string][]string
}

type bodyAssertion struct {
	Kind       Kind
	Expression string `yaml:"exp"`
	Variable   string `yaml:"variable"`
}

type Kind string

const (
	regex    Kind = "regex"
	jsonPath Kind = "jsonpath"
)

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
	log.Println("validCode:", validCode, "validBody:", validBody, "validHeaders:", validHeaders)
	validResp := validCode && validBody && validHeaders
	return validResp, nil
}

// validateCodeStatus verify if the returned http code
// matchs with at least one value provided by user
// In case user did not provide any code, we check code is 2XX
func (a *Assertion) validateCodeStatus(statusCode int) bool {

	if a.Code == nil {
		if statusCode > 199 && statusCode < 300 {
			return true
		}
	} else {
		for _, code := range a.Code {
			if code == statusCode {
				return true
			}
		}
	}
	return false
}

// validateBody verify if body match with the expression given
// in assertions.Body, it can be regex, jsonPath (or other in the future)
func (a *Assertion) validateBody(raw []byte, m map[string]string) (bool, error) {
	valid := true
	for _, b := range a.Body {
		switch b.Kind {
		case regex:
			matched, err := validateWithRegex(raw, b.Expression)
			if err != nil {
				return false, err
			}
			if b.Variable != "" && matched {
				result, err := getFromRegex(raw, b.Expression)
				if err != nil {
					return false, err
				}
				m[b.Variable] = result
			}
			valid = valid && matched
		case jsonPath:
			matched, result, err := getFromJsonPath(raw, b.Expression)
			if err != nil {
				return false, err
			}
			if b.Variable != "" && matched {
				m[b.Variable] = result
			}
			valid = valid && matched
		}

	}
	return valid, nil
}

// validateHeaders verify for each key []values inserted by user if
// the all the values exists for the associated key.
// Example : - Access-Control-Allow-Origin: ["*"]
// The function will find if the header exists with the associated value.
func (a *Assertion) validateHeaders(headers *http.Header) bool {
	valid := true
	// loop over table of maps in Assertion.Headers
	for _, m := range a.Headers {
		// loop over each map
		for headerName, headerValues := range m {
			// loop over the assertions values to execute contains function
			for _, headerValue := range headerValues {
				// We use the AND operation (&&). This way if a value define in assertion
				// is not present, the request is considered not valid
				valid = valid && contains(headers.Values(headerName), headerValue)
			}
		}
	}
	return valid
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
