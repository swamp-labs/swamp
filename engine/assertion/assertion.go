package assertion

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type Assertion struct {
	Body    []bodyAssertion
	Code    []int
	Headers []map[string][]string
}

// Assertion defines a check to execute against response
type bodyAssertion struct {
	Kind       Kind
	Target     string
	Expression string `yaml:"exp"`
	Values     []int
	Variable   string
}

type Kind string

const (
	regex    Kind = "regex"
	jsonPath Kind = "jsonpath"
)

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
			for _, ba := range a.Body {
				v, _, _ := ba.validateBody(body, m)
				if !v {
					validBody = false
				}
			}
		}
	}
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

// validateBody verify if body matchs with the expression given
// in assertions.Body, it can be regex, jsonPath (or other in the future)
func (ba *bodyAssertion) validateBody(raw []byte, m map[string]string) (bool, map[string]string, error) {

	switch ba.Kind {
	case regex:
		v, err := validateWithRegex(raw, ba.Expression)
		if err != nil {
			return false, nil, err
		}
		if ba.Variable != "" && v {
			result, err := getFromRegex(raw, ba.Expression)
			if err != nil {
				return false, nil, err
			}
			if result != "" {
				m[ba.Variable] = result
			}
		}
		return v, m, nil
	case jsonPath:
		matched, result, err := getFromJsonPath(raw, ba.Expression)
		if err != nil {
			return false, nil, err
		}
		if result != "" && ba.Variable != "" {
			m[ba.Variable] = result
		}
		return matched, m, nil
	}
	return false, nil, errors.New("error: no assertion kind found")
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
		for key, values := range m {
			// loop over the assertions values to execute contains function
			for _, value := range values {
				// We use the AND operation (&&). This way if a value define in assertion
				// is not present, the request is considered not valid
				valid = valid && contains(headers.Values(key), value)
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
