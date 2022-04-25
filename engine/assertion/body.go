package assertion

// BodyAssert define checks to execute against the response body
type BodyAssert struct {
	JsonPath string
	Regex    string
	Variable string
}

// validateBody verify if body match with the BodyAssertion given
// in assertions.Body, it can be regex, jsonPath (or other in the future)
func (a *Assertion) validateBody(raw []byte, m map[string]string) (bool, error) {
	valid := true
	for _, exp := range a.Body {
		matched, _ := exp.validate(raw, m)
		valid = valid && matched
	}
	return valid, nil
}

type BodyAssertion interface {
	// TODO : m parameter shouldn't be a part of the validate function & function should be renamed to validateBody
	validate(raw []byte, m map[string]string) (bool, error)
}

type jsonPath struct {
	jsonpath string `yaml:"jsonpath"`
	variable string `yaml:"variable"`
}

type regex struct {
	regex    string `yaml:"regex"`
	variable string `yaml:"variable"`
}

func NewRegexAssertion(expression string, variable string) BodyAssertion {
	return &regex{
		regex:    expression,
		variable: variable,
	}
}

func NewJsonPathAssertion(jsonpath string, variable string) BodyAssertion {
	return &jsonPath{
		jsonpath: jsonpath,
		variable: variable,
	}
}

func (j *jsonPath) validate(raw []byte, m map[string]string) (bool, error) {
	matched, result, err := getFromJsonPath(raw, j.jsonpath)
	if err != nil {
		return false, err
	}
	if j.variable != "" && matched {
		m[j.variable] = result
	}
	return matched, nil
}

func (r *regex) validate(raw []byte, m map[string]string) (bool, error) {
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
