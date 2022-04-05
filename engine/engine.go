package main

import (
	"flag"
	"fmt"
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
		fmt.Println("key:", k)
		fmt.Println("value", t[k])
		table := t[k]
		sessionVar := make(map[string]string)

		for _, r := range table {
			v, err := r.Execute(sessionVar)
			if err != nil {
				log.Println(err)
			}
			log.Println(r, v)

		}
		log.Println("sessionVar:", sessionVar)

	}
}
