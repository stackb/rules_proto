package main

import (
	"flag"
	"log"

	"github.com/stackb/rules_proto/tools/protorule"
)

var (
	config = flag.String("rule_json", "", "The JSON configuration file")
)

func main() {
	flag.Parse()

	rule, err := protorule.FromJSONFile(*config)
	if err != nil {
		log.Fatalf("protorule: %v", err)
	}

	if err := protorule.Generate(rule); err != nil {
		log.Fatalf("protorule: %v", err)
	}
}
