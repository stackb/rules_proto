package main

import (
	"flag"
	"log"
)

var config = flag.String("config_json", "", "The JSON configuration file")

func main() {
	flag.Parse()

	if *config == "" {
		log.Fatalln("error: --config_json is required")
	}

	c, err := fromJSON(*config)
	if err != nil {
		log.Fatalf("depsgen: %v", err)
	}

	if err := generate(c); err != nil {
		log.Fatalf("depsgen: %v", err)
	}
}
