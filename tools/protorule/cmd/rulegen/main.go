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

	rule, err := protorule.NewProtoRuleFromJSONFile(*config)
	if err != nil {
		log.Fatalf("rulegen: %v", err)
	}

	if err := rule.Generate(); err != nil {
		log.Fatalf("rulegen: %v", err)
	}
}
