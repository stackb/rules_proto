package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/stackb/rules_proto/tools/protorule"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:      filepath.Base(os.Args[0]),
		HelpName:  filepath.Base(os.Args[0]),
		Usage:     "Bazel proto rule generator",
		UsageText: filepath.Base(os.Args[0]) + " FLAGS",
		Action: func(c *cli.Context) error {
			rule, err := protorule.FromJSONFile(c.String("rule_json"))
			if err != nil {
				return cli.NewExitError(fmt.Errorf("could not generate rule: %w", err), 1)
			}
			if err := protorule.ParseRuleTemplates(rule); err != nil {
				return cli.NewExitError(fmt.Errorf("could not generate rule: %w", err), 1)
			}
			if err := protorule.Generate(rule); err != nil {
				return cli.NewExitError(fmt.Errorf("could not generate rule: %w", err), 1)
			}
			return nil
		},
		Writer: os.Stdout,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "rule_json",
				Usage: "Config JSON file for the rule",
			},
		},
	}

	app.Run(os.Args)
}
