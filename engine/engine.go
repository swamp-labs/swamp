package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	var (
		configFile = flag.String("config.file", os.Getenv("CONFIG_FILE"), " Simulation configuration file name.")
	)

	flag.Parse()
	config, err := ConfigReader(*configFile)
	if err != nil {
		log.Fatal("Error while reading YAML file: ", err)
	}
	log.Println(config)
	t := config.Simulation

	for k := range t {
		table := t[k]
		sessionVar := make(map[string]string)

		for _, r := range table {
			_, err := r.Execute(sessionVar)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
