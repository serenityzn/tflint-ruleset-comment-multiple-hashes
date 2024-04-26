package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

// AwsInstanceExampleTypeRule checks whether ...
type AwsInstanceExampleTypeRule struct {
	tflint.DefaultRule
}

// NewAwsInstanceExampleTypeRule returns a new rule
func NewAwsInstanceExampleTypeRule() *AwsInstanceExampleTypeRule {
	return &AwsInstanceExampleTypeRule{}
}

// Name returns the rule name
func (r *AwsInstanceExampleTypeRule) Name() string {
	return ""
}

// Enabled returns whether the rule is enabled by default
func (r *AwsInstanceExampleTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsInstanceExampleTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsInstanceExampleTypeRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *AwsInstanceExampleTypeRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	myFiles, err := runner.GetFiles()
	if err != nil {
		return err
	}

	err = checkHashes(myFiles, runner, r)
	if err != nil {
		return err
	}

	return nil
}

func checkHashes(files map[string]*hcl.File, runner tflint.Runner, r *AwsInstanceExampleTypeRule) error {
	var hclRange hcl.Range
	for _, value := range files {
		lines := strings.Split(string(value.Bytes), "\n")
		for index, line := range lines {
			if strings.Count(line, "#") > 1 {
				hclRange = value.Body.MissingItemRange()
				hclRange.Start.Line = index + 1
				err := runner.EmitIssue(r, fmt.Sprintf("Multiple hashes in one line. [%s]", line), hclRange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
