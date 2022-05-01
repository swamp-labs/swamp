package config

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestDecodeRegexBodyAssertion(t *testing.T) {
	yamlString := `
assertions:
  regex: '(\d{4,})'
  variable: myVar
`
	yNode := struct {
		Assertions yaml.Node
	}{}
	err := yaml.Unmarshal([]byte(yamlString), &yNode)
	if err != nil {
		t.Errorf("Yaml to yml.Node casting failed : %v", err)
	}

	assert, err := yamlNodeToBodyAssertion(yNode.Assertions)
	if err != nil {
		t.Errorf("Fail to parse node to assertion : %v", err)
	}

	if assert.Type() != "REGEX" {
		t.Errorf("Invalid assertion, expected REGEX, got %s", assert.Type())
	}
}

func TestDecodeJsonPathBodyAssertion(t *testing.T) {
	yamlString := `
assertions:
  jsonpath: $.id
  variable: myVar
`
	yNode := struct {
		Assertions yaml.Node
	}{}
	err := yaml.Unmarshal([]byte(yamlString), &yNode)
	if err != nil {
		t.Errorf("Yaml to yml.Node casting failed : %v", err)
	}

	assert, err := yamlNodeToBodyAssertion(yNode.Assertions)
	if err != nil {
		t.Errorf("Fail to parse node to assertion : %v", err)
	}

	if assert.Type() != "JSONPATH" {
		t.Errorf("Invalid assertion, expected JSONPATH, got %s", assert.Type())
	}
}
