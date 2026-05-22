package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_CommentMultipleHashesRuleType(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "Valid one hash comments",
			Content: `#Valid comment
resource "aws_instance" "test" { # valid comment
  ami = "ami-12345"
  instance_type = "asss.micro"
}
     # Also valid comment
`,
			Expected: helper.Issues{},
		},
		{
			Name: "Valid one hash comment",
			Content: `#Valid comment
resource "aws_instance" "test" {
  ami = "ami-12345"
  instance_type = "asss.micro"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "issue found Line1",
			Content: `######
resource "aws_instance" "test" {
  ami = "ami-12345"
  instance_type = "asss.micro"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [######]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "Valid hashes inside string literals",
			Content: `resource "example" "test" {
  expression = "(http.host eq \"api-survey.prod.oddity.com\" and all(http.request.headers[\"api-key\"][*] ne \"+-e-idyR}&r6IwY9WEsHT=]7gNMPZ#8n!fZK}p1<9HRfAlHORN*@P3Adf@#xV<E0YGVO$DUY@eYC+C&q\"))"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "Valid inline comment with single hash",
			Content: `resource "aws_instance" "test" { # valid comment
  ami = "ami-#12345"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "issue found comment with decorative hashes",
			Content: `# Comment  ######### --- ######
`,
			Expected: helper.Issues{
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [# Comment  ######### --- ######]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "issue found double hash comment",
			Content: `##
`,
			Expected: helper.Issues{
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [##]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "issue found Multiple lines (1,4,10)",
			Content: `######
resource "aws_instance" "test" {
  ami = "ami-12345"
  instance_type = "asss.micro"
}
######
  ######


#########
`,
			Expected: helper.Issues{
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [######]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [######]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [  ######]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
				{
					Rule:    NewCommentMultipleHashesRule(),
					Message: "Multiple hash symbols (#) in one line. [#########]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
	}

	rule := NewCommentMultipleHashesRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
