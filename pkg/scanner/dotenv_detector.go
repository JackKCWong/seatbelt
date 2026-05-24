package scanner

import (
	"regexp"
)

type DotEnvDetector struct {
	pattern *regexp.Regexp
}

func NewDotEnvDetector() *DotEnvDetector {
	return &DotEnvDetector{
		pattern: regexp.MustCompile(`(?i)\.env`),
	}
}

func (d *DotEnvDetector) Name() string {
	return "dotenv"
}

func (d *DotEnvDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding
	if d.pattern.MatchString(prompt) {
		findings = append(findings, Finding{
			Type:     "dotenv",
			Value:    ".env",
			Location: "file reference",
		})
	}
	return findings, nil
}