package main

import (
	"flag"
	"log"

	"github.com/stackb/rules_proto/tools/protorule"
)

var (
	config = flag.String("language_json", "", "The JSON configuration file")
)

func main() {
	flag.Parse()

	lang, err := protorule.NewProtoLanguageFromJSONFile(*config)
	if err != nil {
		log.Fatalf("langgen: %v", err)
	}

	if err := lang.Generate(); err != nil {
		log.Fatalf("langgen: %v", err)
	}
}
