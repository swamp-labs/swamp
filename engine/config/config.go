package config

import (
	"fmt"
	"github.com/swamp-labs/swamp/engine/simulation"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// Parse read a YAML configuration file to extract a simulation.
func Parse(configFile string) (simulation.Simulation, error) {
	log.Println("Parsing configuration file...")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to read config file %s : %v", configFile, err)
	}

	config := SimulationTemplate{}
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, fmt.Errorf("fail to parse config file : %v", err)
	}
	sim, err := config.decode()
	if err != nil {
		return nil, fmt.Errorf("fail to decode template : %v", err)
	}
	return sim, nil
}
