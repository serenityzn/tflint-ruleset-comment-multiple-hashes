package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

// CommentMultipleHashesRule checks whether ...
type CommentMultipleHashesRule struct {
	tflint.DefaultRule
}

// NewCommentMultipleHashesRule returns a new rule
func NewCommentMultipleHashesRule() *CommentMultipleHashesRule {
	return &CommentMultipleHashesRule{}
}

// Name returns the rule name
func (r *CommentMultipleHashesRule) Name() string {
	return "turyachka_detected"
}

// Enabled returns whether the rule is enabled by default
func (r *CommentMultipleHashesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *CommentMultipleHashesRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *CommentMultipleHashesRule) Link() string {
	return ""
}

// Check checks whether comments has more than one hash symbol
func (r *CommentMultipleHashesRule) Check(runner tflint.Runner) error {
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

// commentHashCount returns how many '#' appear in the line's HCL comment (# to EOL).
// '#' inside quoted strings are ignored.
func commentHashCount(line string) int {
	inString := false
	var quote byte
	escaped := false
	commentStart := -1

	for i := 0; i < len(line); i++ {
		c := line[i]
		if commentStart >= 0 {
			break
		}
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == quote {
				inString = false
			}
			continue
		}
		if c == '"' || c == '\'' {
			inString = true
			quote = c
			continue
		}
		if c == '#' {
			commentStart = i
		}
	}

	if commentStart < 0 {
		return 0
	}
	return strings.Count(line[commentStart:], "#")
}

func checkHashes(files map[string]*hcl.File, runner tflint.Runner, r *CommentMultipleHashesRule) error {
	var hclRange hcl.Range
	for _, value := range files {
		lines := strings.Split(string(value.Bytes), "\n")
		for index, line := range lines {
			if commentHashCount(line) > 1 {
				hclRange = value.Body.MissingItemRange()
				hclRange.Start.Line = index + 1
				err := runner.EmitIssue(r, fmt.Sprintf("Multiple hash symbols (#) in one line. [%s]", line), hclRange)
				// Put a log that can be output with `TFLINT_LOG=debug`
				logger.Debug(fmt.Sprintf("[DEBUG] Range - Start.Line = %d Start.Column= %d Start.Byte= %d",
					hclRange.Start.Line, hclRange.Start.Column, hclRange.Start.Byte))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
