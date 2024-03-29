package main

import (
	"flag"
	"fmt"
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

	logger.Engine.Info(fmt.Sprintf("%+v", config))
	tasks := config.GetTasks()

	for _, task := range tasks {
		sessionVar := make(map[string]string)

		for _, r := range task.GetRequest() {
			_, err := r.Execute(sessionVar)
			if err != nil {
				logger.Engine.Info(err.Error())
			}
		}
	}
}
