package main

import (
	"io/ioutil"
	"log"

	"github.com/swamp-labs/swamp/engine/httpreq"
	"gopkg.in/yaml.v3"
)

// Simulation wraps configuration file content
type Simulation struct {
	Simulation map[string][]httpreq.Request
	Volumes    []Volume `yaml:"volumes"`
}

// Volume defines an execution plan.
type Volume struct {
	RequestGroup int `yaml:"request_group"`
	//	Execution    []string `yaml:"execution"`
}

// ConfigReader read a YAML configuration file to extract a simulation.
func ConfigReader(configFile string) (Simulation, error) {
	log.Println("Parsing configuration file...")
	data, err := ioutil.ReadFile(configFile)
	var config Simulation
	yaml.Unmarshal(data, &config)
	return config, err
}
