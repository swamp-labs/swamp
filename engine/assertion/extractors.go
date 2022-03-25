package assertion

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/yalp/jsonpath"
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
	if result != nil {
		for i := range result {
			s = s + string(result[i][1])
		}
		log.Println("concat value:", s)
	}
	return s, nil
}

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
// also returns a boolean to indicate if a match has been found
func getFromJsonPath(raw []byte, exp string) (bool, string, error) {
	var body interface{}
	err := json.Unmarshal(raw, &body)
	if err != nil {
		return false, "", nil
	}
	result, err := jsonpath.Read(body, exp)
	log.Println("Jsonpath result:", result, "expression:", exp)
	if err != nil {
		return false, "", nil
	}
	if result != nil {
		log.Println(result)
		variable := result.(string)

		return true, variable, nil
	}

	return false, "", nil
}
