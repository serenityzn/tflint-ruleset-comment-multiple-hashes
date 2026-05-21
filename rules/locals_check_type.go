package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

// LocalsRule checks whether locals blocks are declared in locals.tf
type LocalsRule struct {
	tflint.DefaultRule
}

// NewLocalsRule returns a new rule
func NewLocalsRule() *LocalsRule {
	return &LocalsRule{}
}

// Name returns the rule name
func (r *LocalsRule) Name() string {
	return "locals_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *LocalsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *LocalsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *LocalsRule) Link() string {
	return ""
}

// Check checks whether locals blocks are declared in locals.tf
func (r *LocalsRule) Check(runner tflint.Runner) error {
	myFiles, err := runner.GetFiles()
	if err != nil {
		return err
	}

	return checkLocals(myFiles, runner, r)
}

func checkLocals(files map[string]*hcl.File, runner tflint.Runner, r *LocalsRule) error {
	for fileName, value := range files {
		err := checkLocalsFile(runner, fileName, value, r)
		if err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("FileName is ... %s. (locals.tf - skipped)", fileName))
	}

	return nil
}

func checkLocalsFile(runner tflint.Runner, fileNameWithPath string, value *hcl.File, r *LocalsRule) error {
	paths := strings.Split(fileNameWithPath, "/")
	fileName := paths[len(paths)-1]
	logger.Debug(fmt.Sprintf("(LocalsCheck) Full fileName is ... %s. Processing. Terraform file is ... %s", fileNameWithPath, fileName))
	lines := strings.Split(string(value.Bytes), "\n")
	for index, line := range lines {
		if strings.Contains(line, "locals {") && fileName != "locals.tf" {
			hclRange := value.Body.MissingItemRange()
			hclRange.Start.Line = index + 1
			err := runner.EmitIssue(r, fmt.Sprintf("Locals declaration in wrong file. "+
				"(All locals should be declared in locals.tf) [%s]", line), hclRange)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
