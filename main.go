package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-template/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "turyachka",
			Version: "0.1.3",
			Rules: []tflint.Rule{
				rules.NewCommentMultipleHashesRule(),
				rules.NewVariablesRule(),
			},
		},
	})
}
