package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/swamp-labs/swamp/engine/assertion"
	"gopkg.in/yaml.v3"
)

func castTemplateNodeToBodyAssertion[T bodyAssertionTemplate](node yaml.Node, dst T) (*T, error) {
	err := node.Decode(&dst)
	if err != nil {
		return nil, fmt.Errorf("fail to cast node : %v", err)
	}
	validate := validator.New()
	err = validate.Struct(dst)
	if err != nil {
		return nil, fmt.Errorf("can not validate structure : %v", err)
	}
	return &dst, nil
}

// decode transform an assertionsBlockTemplate into an assertion.Assertion
func (assertions *assertionsBlockTemplate) decode() (assertion.Assertion, error) {

	bodyAssertions := make([]assertion.BodyAssertion, cap(assertions.Body), len(assertions.Body))

	// Decode all yaml nodes to body assertions
	for i, node := range assertions.Body {
		ass, err := yamlNodeToBodyAssertion(node)
		if err != nil {
			return nil, fmt.Errorf("fail to parse assertion : %v", err)
		}
		bodyAssertions[i] = ass
	}

	return assertion.MakeRequestAssertion(bodyAssertions, assertions.Code, assertions.Headers), nil

}

// yamlNodeToBodyAssertion determines the BodyAssertion type using
// cast functions, returns an error if the node does not match with any
func yamlNodeToBodyAssertion(node yaml.Node) (assertion.BodyAssertion, error) {
	var dst bodyAssertionTemplate

	var err error
	// Regex assertion
	dst, err = castTemplateNodeToBodyAssertion[regexAssertionTemplate](node, regexAssertionTemplate{})
	if err == nil {
		return dst.toAssertion(), nil
	}

	// JsonPath assertion
	dst, err = castTemplateNodeToBodyAssertion[jsonPathAssertionTemplate](node, jsonPathAssertionTemplate{})
	if err == nil {
		return dst.toAssertion(), nil
	}

	return nil, fmt.Errorf("invalid assertion")
}

type bodyAssertionTemplate interface {
	toAssertion() assertion.BodyAssertion
}

type regexAssertionTemplate struct {
	Regex    string `yaml:"regex" validate:"required"`
	Variable string `yaml:"variable" validate:"required"`
}

// toAssertion transforms a regexAssertionTemplate into an assertion.BodyAssertion
func (tpl regexAssertionTemplate) toAssertion() assertion.BodyAssertion {
	return assertion.NewRegexAssertion(tpl.Regex, tpl.Variable)
}

type jsonPathAssertionTemplate struct {
	JsonPath string `yaml:"jsonpath" validate:"required"`
	Variable string `yaml:"variable" validate:"required"`
}

// toAssertion transforms a jsonPathAssertionTemplate into an assertion.BodyAssertion
func (tpl jsonPathAssertionTemplate) toAssertion() assertion.BodyAssertion {
	return assertion.NewJsonPathAssertion(tpl.JsonPath, tpl.Variable)
}
