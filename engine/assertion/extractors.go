package assertion

import (
	"encoding/json"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
)

// getFromRegex search regular expression in []byte passed in arguments
// First returned value is a booleaan (true means exp match against raw, false instead)
// func also returns a string which concat all successive matches of the expression
func getFromRegex(raw []byte, exp string) (string, error) {
	var s string
	re, err := regexp.Compile(exp)
	if err != nil {
		return "", err
	}
	result := re.FindAllSubmatch(raw, -1)
	for i := range result {
		if len(result[i]) > 1 {
			s = s + string(result[i][1])
		}
	}
	return s, nil
}

// validateWithRegex check if the regex expression match with the
// raw []byte passed in parameters, returns boolean with the result
func validateWithRegex(raw []byte, exp string) (bool, error) {
	re, err := regexp.Compile(exp)
	if err != nil {
		return false, err
	}
	v := re.Match(raw)
	return v, nil
}

// getFromJsonPath reads a path from a decoded JSON array
// and returns the corresponding value or an error.
// Also returns a boolean that indicates if a match has been found
func getFromJsonPath(raw []byte, exp string) (bool, string, error) {
	var body interface{}
	err := json.Unmarshal(raw, &body)
	if err != nil {
		return false, "", err
	}
	result, err := jsonpath.Get(exp, body)
	if err != nil {
		return false, "", err
	}
	if result != nil {
		variable := result.(string)

		return true, variable, nil
	}
	return false, "", nil
}
