package rules

import (
	"github.com/hashicorp/hcl/v2"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_VariablesRuleType(t *testing.T) {
	tests := []struct {
		Name     string
		FileName string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "proper variable placement",
			FileName: "variables.tf",
			Content: `variable "test" { 
  type = string
  default = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:     "proper data placement",
			FileName: "data.tf",
			Content: `data "aws_ec2_host" "test" {
  name = "Test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:     "proper ouptut placement",
			FileName: "outputs.tf",
			Content: `output "name" {
  value = aws_instance.test.ami
}
`,
			Expected: helper.Issues{},
		},
		{
			Name:     "Wrong variable placement",
			FileName: "main.tf",
			Content: `variable "test" { 
  type = string
  default = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewVariablesRule(),
					Message: "Variable declaration in wrong file. (All variables should be declared in variables.tf) [variable \"test\" { ]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name:     "wrong data placement",
			FileName: "main.tf",
			Content: `data "aws_ec2_host" "test" {
  name = "Test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewVariablesRule(),
					Message: "Data declaration in wrong file. (All datas should be declared in data.tf) [data \"aws_ec2_host\" \"test\" {]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name:     "Wrong ouptut placement",
			FileName: "main.tf",
			Content: `output "name" {
  value = aws_instance.test.ami
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewVariablesRule(),
					Message: "Output declaration in wrong file. (All outputs should be declared in outputs.tf) [output \"name\" {]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
	}

	rule := NewVariablesRule()

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
