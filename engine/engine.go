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
	config.Run()
	//_ := config.GetTasks()

}
