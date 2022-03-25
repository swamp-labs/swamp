package assertion

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type Assertion struct {
	Body    []bodyAssertion
	Code    []int
	Headers []map[string]string
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

func AssertResponse(as Assertion, resp *http.Response) (valid bool, variables map[string]string, err error) {
	variables = make(map[string]string)
	var validCode bool
	validBody := true
	validHeaders := true
	if resp.ContentLength > 0 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, nil, err
		}
		if as.Code == nil {
			if resp.StatusCode > 199 && resp.StatusCode < 299 {
				validCode = true
			}
		} else {
			for _, code := range as.Code {
				if resp.StatusCode == code {
					validCode = true
				}
			}
		}
		if as.Body != nil {
			for _, ba := range as.Body {
				v, _, _ := ba.validateBody(body, variables)
				if !v {
					validBody = false
				}
			}
		}
	}
	log.Println("variables:", variables)
	validResp := validCode && validBody && validHeaders
	return validResp, nil, nil
}

// validateBody assert body
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
	return false, nil, errors.New("no assertion kind found")
}

func validateHeaders() {}
