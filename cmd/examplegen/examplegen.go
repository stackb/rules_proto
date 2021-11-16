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
		log.Fatalf("examplegen: %v", err)
	}

	if err := generateTest(c); err != nil {
		log.Fatalf("examplegen test: %v", err)
	}

	if err := generateMarkdown(c); err != nil {
		log.Fatalf("examplegen markdown: %v", err)
	}
}
