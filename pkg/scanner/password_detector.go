package scanner

import (
	"regexp"
)

type PasswordDetector struct {
	patterns []*regexp.Regexp
}

func NewPasswordDetector() *PasswordDetector {
	return &PasswordDetector{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(password\s*[=:]\s*[^\s]+)`),
			regexp.MustCompile(`(?i)(pwd\s*[=:]\s*[^\s]+)`),
			regexp.MustCompile(`(?i)(passwd\s*[=:]\s*[^\s]+)`),
		},
	}
}

func (d *PasswordDetector) Name() string {
	return "password"
}

func (d *PasswordDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding
	for _, pattern := range d.patterns {
		matches := pattern.FindAllString(prompt, -1)
		for _, match := range matches {
			findings = append(findings, Finding{
				Type:     "password",
				Value:    match,
				Location: "prompt",
			})
		}
	}
	return findings, nil
}