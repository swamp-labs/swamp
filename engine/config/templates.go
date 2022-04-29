package config

import "gopkg.in/yaml.v3"

type assertionsBlockTemplate struct {
	Body    []yaml.Node           `yaml:"body"`
	Code    []int                 `yaml:"code"`
	Headers []map[string][]string `yaml:"headers"`
}

type requestTemplate struct {
	Name            string                  `yaml:"name"`
	Method          string                  `yaml:"method"`
	Protocol        string                  `yaml:"protocol"`
	Headers         []map[string]string     `yaml:"headers"`
	URL             string                  `yaml:"url"`
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
