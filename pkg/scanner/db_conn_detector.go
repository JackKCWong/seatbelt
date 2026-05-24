package scanner

import (
	"regexp"
)

type DBConnectionDetector struct {
	patterns []*regexp.Regexp
}

func NewDBConnectionDetector() *DBConnectionDetector {
	return &DBConnectionDetector{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(mongodb\+srv://[^\s]+)`),
			regexp.MustCompile(`(?i)(postgres://[^\s]+)`),
			regexp.MustCompile(`(?i)(mysql://[^\s]+)`),
			regexp.MustCompile(`(?i)(redis://[^\s]+)`),
			regexp.MustCompile(`(?i)(sqlserver://[^\s]+)`),
		},
	}
}

func (d *DBConnectionDetector) Name() string {
	return "db_connection"
}

func (d *DBConnectionDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding
	for _, pattern := range d.patterns {
		matches := pattern.FindAllString(prompt, -1)
		for _, match := range matches {
			findings = append(findings, Finding{
				Type:     "db_connection",
				Value:    match,
				Location: "prompt",
			})
		}
	}
	return findings, nil
}