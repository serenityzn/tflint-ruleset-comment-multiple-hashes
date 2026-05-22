package rules

import (
	"github.com/hashicorp/hcl/v2"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_LocalsRuleType(t *testing.T) {
	tests := []struct {
		Name     string
		FileName string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "proper locals placement",
			FileName: "locals.tf",
			Content: `locals {
  env = "prod"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:     "proper locals placement with long path",
			FileName: "/tmp/terraform/modules/testmodule/locals.tf",
			Content: `locals {
  env = "prod"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:     "wrong locals placement in main.tf",
			FileName: "main.tf",
			Content: `locals {
  env = "prod"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewLocalsRule(),
					Message: "Locals declaration in wrong file. (All locals should be declared in locals.tf) [locals {]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name:     "wrong locals placement in variables.tf",
			FileName: "variables.tf",
			Content: `locals {
  env = "prod"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewLocalsRule(),
					Message: "Locals declaration in wrong file. (All locals should be declared in locals.tf) [locals {]",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
	}

	rule := NewLocalsRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{test.FileName: test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
