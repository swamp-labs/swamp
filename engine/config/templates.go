package config

import "gopkg.in/yaml.v3"

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
	Name            string                  `yaml:"name"`
	Method          string                  `yaml:"method"`
	Protocol        string                  `yaml:"protocol"`
	Headers         []map[string]string     `yaml:"headers"`
	URL             yaml.Node               `yaml:"url"`
	Body            string                  `yaml:"body"`
	QueryParameters map[string]string       `yaml:"query_parameters"`
	Assertions      assertionsBlockTemplate `yaml:"assertions"`
}

// volumeTemplate defines an execution plan.
type volumeTemplate struct {
	RequestGroup string `yaml:"request_group"`
	//	Execution    []string `yaml:"execution"`
}

// simulationTemplate wraps configuration file content
type simulationTemplate struct {
	Groups  map[string][]requestTemplate `yaml:"groups"`
	Volumes []volumeTemplate             `yaml:"volume"`
}
