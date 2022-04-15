package assertion

// BodyAssert define checks to execute against the response body
type BodyAssert struct {
	JsonPath string
	Regex    string
	Variable string
}

// validateBody verify if body match with the expression given
// in assertions.Body, it can be regex, jsonPath (or other in the future)
func (a *Assertion) validateBody(raw []byte, m map[string]string) (bool, error) {
	valid := true
	for _, exp := range a.Body {
		var e expression
		if exp.JsonPath != "" {
			e = &JsonPath{jsonpath: exp.JsonPath, variable: exp.Variable}

		}
		if exp.Regex != "" {
			e = &Regex{regex: exp.Regex, variable: exp.Variable}
		}
		matched, _ := e.validate(raw, m)
		valid = valid && matched
	}
	return valid, nil
}

type expression interface {
	validate(raw []byte, m map[string]string) (bool, error)
}

type JsonPath struct {
	jsonpath string `yaml:"jsonpath"`
	variable string `yaml:"variable"`
}

type Regex struct {
	regex    string `yaml:"regex"`
	variable string `yaml:"variable"`
}

func (j *JsonPath) validate(raw []byte, m map[string]string) (bool, error) {
	matched, result, err := getFromJsonPath(raw, j.jsonpath)
	if err != nil {
		return false, err
	}
	if j.variable != "" && matched {
		m[j.variable] = result
	}
	return matched, nil
}

func (r *Regex) validate(raw []byte, m map[string]string) (bool, error) {
	matched, err := validateWithRegex(raw, r.regex)
	if err != nil {
		return false, err
	}
	if r.variable != "" && matched {
		result, err := getFromRegex(raw, r.regex)
		if err != nil {
			return false, err
		}
		m[r.variable] = result
	}
	return matched, nil

}
