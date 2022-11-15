package main

import (
	"flag"
	cnf "github.com/swamp-labs/swamp/engine/config"
	"github.com/swamp-labs/swamp/engine/logger"
	"os"
)

func main() {

	var (
		configFile = flag.String("config.file", os.Getenv("CONFIG_FILE"), " Simulation configuration file name.")
	)

	flag.Parse()
	config, err := cnf.Parse(*configFile)
	if err != nil {
		logger.Engine.Error("Error while reading YAML file: ", err)
		os.Exit(3)
	}

	//logger.EngineLogger.Info(config)
	groups := config.GetGroups()

	for _, table := range groups {
		sessionVar := make(map[string]string)

		for _, r := range table {
			_, err := r.Execute(sessionVar)
			if err != nil {
				//logger.EngineLogger.Info(err)
			}
		}
	}
}
