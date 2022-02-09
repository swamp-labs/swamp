package main

import (
	"flag"
	"log"
	"os"
)

func main() {

	var (
		configFile = flag.String("config.file", os.Getenv("CONFIG_FILE"), "Synthetic manager configuration file name.")
	)

	flag.Parse()
	config := ConfigReader(*configFile)
	log.Println(config)
}
