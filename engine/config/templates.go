package config

import (
	"github.com/swamp-labs/swamp/engine/templateString"
	"gopkg.in/yaml.v3"
)

// assertionsBlockTemplate defines the yaml template used to represent a request
// assertions.
type assertionsBlockTemplate struct {
	// Body is a list of nodes that contains all assertions needed to validate a body
	Body []yaml.Node `yaml:"body"`
	// Code is the list of allowed status code for a request
	Code []int `yaml:"code"`

	Headers []map[string][]string `yaml:"headers"`
}

// requestTemplate defines the yaml template used to define a request
type requestTemplate struct {
	Name            string                        `yaml:"name"`
	Method          string                        `yaml:"method"`
	Protocol        string                        `yaml:"protocol"`
	Headers         []map[string]string           `yaml:"headers"`
	URL             templateString.TemplateString `yaml:"url"`
	Body            templateString.TemplateString `yaml:"body"`
	QueryParameters map[string]string             `yaml:"query_parameters"`
	Assertions      assertionsBlockTemplate       `yaml:"assertions"`
}

// simulationTemplate wraps configuration file content
type simulationTemplate struct {
	Tasks map[string]taskTemplate `yaml:"tasks"`
}

// taskTemplate defines a list of actions to execute, can be compared to a user scenario
type taskTemplate struct {
	// Requests lists to http request to execute in plan
	Requests []requestTemplate `yaml:"requests"`
	// Volume defines an execution plan.
	Volume []map[string]int `yaml:"volume"`
}
