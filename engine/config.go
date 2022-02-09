package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Simulation wraps configuration file content
type Simulation struct {
	Groups  []Group
	Volumes []Volume `yaml:"volumes"`
}

// Group defines a list of requests
type Group struct {
	Requests []Request
}

// Request struct defines all parameters for http requests to execute
type Request struct {
	Name       string
	Verb       string
	URL        string
	Body       string
	Parameters []string
	Assertions []string
	SaveVal    []string
}

// Volume defines an execution plan.
type Volume struct {
	RequestGroup int      `yaml:"requestGroup"`
	Execution    []string `yaml:"execution"`
}

// ConfigReader read a YAML configuration file to extract a simulation.
func ConfigReader(configFile string) Simulation {
	log.Println("Parsing configuration file...")

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error while reading YAML file: ", err)
	}
	var config Simulation
	yaml.Unmarshal(data, &config)
	return config
}
