package config

import "gopkg.in/yaml.v3"

type AssertionsBlockTemplate struct {
	Body    []yaml.Node           `yaml:"body"`
	Code    []int                 `yaml:"code"`
	Headers []map[string][]string `yaml:"headers"`
}

type RequestTemplate struct {
	Name            string                  `yaml:"name"`
	Method          string                  `yaml:"method"`
	Protocol        string                  `yaml:"protocol"`
	Headers         []map[string]string     `yaml:"headers"`
	URL             string                  `yaml:"url"`
	Body            string                  `yaml:"body"`
	QueryParameters map[string]string       `yaml:"query_parameters"`
	Assertions      AssertionsBlockTemplate `yaml:"assertions"`
}

// VolumeTemplate defines an execution plan.
type VolumeTemplate struct {
	RequestGroup string `yaml:"request_group"`
	//	Execution    []string `yaml:"execution"`
}

// SimulationTemplate wraps configuration file content
type SimulationTemplate struct {
	Groups  map[string][]RequestTemplate `yaml:"groups"`
	Volumes []VolumeTemplate             `yaml:"volume"`
}
