package config

import (
	"fmt"
	"github.com/swamp-labs/swamp/engine/logger"
	s "github.com/swamp-labs/swamp/engine/simulation"
	"gopkg.in/yaml.v3"
	"os"
)

// Parse read a YAML configuration file to extract a simulation
func Parse(configFile string) (s.Simulation, error) {
	logger.Engine.Info("Parsing configuration file...")
	data, err := os.ReadFile(configFile)
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
