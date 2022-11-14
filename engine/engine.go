package main

import (
	"flag"
	cnf "github.com/swamp-labs/swamp/engine/config"
	"log"
	"os"
)

func main() {

	var (
		configFile = flag.String("config.file", os.Getenv("CONFIG_FILE"), " Simulation configuration file name.")
	)

	flag.Parse()
	config, err := cnf.Parse(*configFile)
	if err != nil {
		log.Fatal("Error while reading YAML file: ", err)
	}

	log.Println(config)
	tasks := config.GetTasks()

	for _, task := range tasks {
		sessionVar := make(map[string]string)

		for _, r := range task.GetRequest() {
			_, err := r.Execute(sessionVar)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
