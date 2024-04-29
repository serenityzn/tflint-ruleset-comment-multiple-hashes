package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

// VariablesRule checks whether ...
type VariablesRule struct {
	tflint.DefaultRule
}

// NewVariablesRule returns a new rule
func NewVariablesRule() *VariablesRule {
	return &VariablesRule{}
}

// Name returns the rule name
func (r *VariablesRule) Name() string {
	return "vars_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *VariablesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *VariablesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *VariablesRule) Link() string {
	return ""
}

// Check checks whether comments has more than one hash symbol
func (r *VariablesRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	myFiles, err := runner.GetFiles()
	if err != nil {
		return err
	}

	err = checkVars(myFiles, runner, r)
	if err != nil {
		return err
	}

	return nil
}

func checkVars(files map[string]*hcl.File, runner tflint.Runner, r *VariablesRule) error {
	//var hclRange hcl.Range
	for fileName, value := range files {
		err := checkVariables(runner, fileName, value, r)
		if err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("FileName is ... %s. (variables.tf- skipped", fileName))
	}

	return nil
}

func checkVariables(runner tflint.Runner, fileNameWithPath string, value *hcl.File, r *VariablesRule) error {
	paths := strings.Split(fileNameWithPath, "/")
	fileName := paths[len(paths)-1]
	logger.Debug(fmt.Sprintf("(VariablesCheck) Full fileName is ... %s. Processing. Terraform file is ... %s", fileNameWithPath, fileName))
	lines := strings.Split(string(value.Bytes), "\n")
	for index, line := range lines {
		if strings.Contains(line, "variable \"") && fileName != "variables.tf" {
			hclRange := value.Body.MissingItemRange()
			hclRange.Start.Line = index + 1
			err := runner.EmitIssue(r, fmt.Sprintf("Variable declaration in wrong file. "+
				"(All variables should be declared in variables.tf) [%s]", line), hclRange)
			if err != nil {
				return err
			}
		}
		if strings.Contains(line, "data \"") && fileName != "data.tf" {
			hclRange := value.Body.MissingItemRange()
			hclRange.Start.Line = index + 1
			err := runner.EmitIssue(r, fmt.Sprintf("Data declaration in wrong file. "+
				"(All datas should be declared in data.tf) [%s]", line), hclRange)
			if err != nil {
				return err
			}
		}
		if strings.Contains(line, "output \"") && fileName != "outputs.tf" {
			hclRange := value.Body.MissingItemRange()
			hclRange.Start.Line = index + 1
			err := runner.EmitIssue(r, fmt.Sprintf("Output declaration in wrong file. "+
				"(All outputs should be declared in outputs.tf) [%s]", line), hclRange)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
