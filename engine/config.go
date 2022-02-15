package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Simulation wraps configuration file content
type Simulation struct {
	Groups  []Group  `yaml:"groups"`
	Volumes []Volume `yaml:"volumes"`
}

// Group defines a list of requests
type Group struct {
	Requests []Request `yaml:"requests"`
}

// Request struct defines all parameters for http requests to execute
type Request struct {
	Name       string   `yaml:"name"`
	Verb       string   `yaml:"verb"`
	URL        string   `yaml:"URL"`
	Body       string   `yaml:"body"`
	Parameters []string `yaml:"parameters"`
	Assertions []string `yaml:"assertions"`
	SaveVal    []string `yaml:"save_val"`
}

// Volume defines an execution plan.
type Volume struct {
	RequestGroup int      `yaml:"request_group"`
	Execution    []string `yaml:"execution"`
}

// ConfigReader read a YAML configuration file to extract a simulation.
func ConfigReader(configFile string) (Simulation, error) {
	log.Println("Parsing configuration file...")
	data, err := ioutil.ReadFile(configFile)
	var config Simulation
	yaml.Unmarshal(data, &config)
	return config, err
}
