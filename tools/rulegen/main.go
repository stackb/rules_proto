package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:      filepath.Base(os.Args[0]),
		HelpName:  filepath.Base(os.Args[0]),
		Usage:     "Bazel proto rule generator",
		UsageText: filepath.Base(os.Args[0]) + " FLAGS",
		Action: func(c *cli.Context) error {
			rule, err := makeRule(c)
			if err != nil {
				return cli.NewExitError(fmt.Errorf("could not generate rule: %w", err), 1)
			}
			if err := generate(rule); err != nil {
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

func makeRule(c *cli.Context) (*Rule, error) {
	ruleConfig := c.String("rule_json")

	data, err := ioutil.ReadFile(ruleConfig)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	var rule Rule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	tpl, err := template.ParseFiles(rule.ImplementationTmpl)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}
	rule.Implementation = tpl

	tpl, err = template.ParseFiles(rule.WorkspaceExampleTmpl)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}
	rule.WorkspaceExample = tpl

	tpl, err = template.ParseFiles(rule.BuildExampleTmpl)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}
	rule.BuildExample = tpl

	tpl, err = template.ParseFiles(rule.TestTmpl)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}
	rule.Test = tpl

	return &rule, nil
}
