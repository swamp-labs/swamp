package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/swamp-labs/swamp/engine/assertion"
	"gopkg.in/yaml.v3"
)

// castRegexAssertionTemplateNode uses the playground validator function to
// determine if the yaml node matchs with the regexAssertionTemplate structure
func castRegexAssertionTemplateNode(node yaml.Node) (*regexAssertionTemplate, error) {
	dst := regexAssertionTemplate{}
	err := node.Decode(&dst)
	if err != nil {
		return nil, fmt.Errorf("fail to cast node to regex assertion template : %v", err)
	}
	validate := validator.New()
	err = validate.Struct(dst)
	if err != nil {
		return nil, fmt.Errorf("can not validate regex assertion template : %v", err)
	}
	return &dst, nil
}

// castJsonPathAssertionTemplateNode uses the playground validator to
// determine if the yaml node matchs with the jsonPathAssertionTemplate structure
func castJsonPathAssertionTemplateNode(node yaml.Node) (*jsonPathAssertionTemplate, error) {
	dst := jsonPathAssertionTemplate{}
	err := node.Decode(&dst)
	if err != nil {
		return nil, fmt.Errorf("fail to cast node to jsonpath assertion template : %v", err)
	}
	validate := validator.New()
	err = validate.Struct(dst)
	if err != nil {
		return nil, fmt.Errorf("can not validate jsonpath assertion template : %v", err)
	}
	return &dst, nil
}

func (assertions *assertionsBlockTemplate) decode() (assertion.Assertion, error) {

	bodyAssertions := make([]assertion.BodyAssertion, cap(assertions.Body), len(assertions.Body))

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
	dst, err = castRegexAssertionTemplateNode(node)
	if err == nil {
		return dst.toAssertion(), nil
	}

	// JsonPath assertion
	dst, err = castJsonPathAssertionTemplateNode(node)
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

func (tpl regexAssertionTemplate) toAssertion() assertion.BodyAssertion {
	return assertion.NewRegexAssertion(tpl.Regex, tpl.Variable)
}

type jsonPathAssertionTemplate struct {
	JsonPath string `yaml:"jsonpath" validate:"required"`
	Variable string `yaml:"variable" validate:"required"`
}

func (tpl jsonPathAssertionTemplate) toAssertion() assertion.BodyAssertion {
	return assertion.NewJsonPathAssertion(tpl.JsonPath, tpl.Variable)
}
