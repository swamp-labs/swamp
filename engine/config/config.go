package config

import (
	"fmt"
	"io/ioutil"
	"log"

	s "github.com/swamp-labs/swamp/engine/simulation"
	"gopkg.in/yaml.v3"
)

// Parse read a YAML configuration file to extract a simulation
func Parse(configFile string) (s.Simulation, error) {
	log.Println("Parsing configuration file...")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("fail to read config file %s : %v", configFile, err)
	}

	config := simulationTemplate{}
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
