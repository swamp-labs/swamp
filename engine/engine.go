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
}
