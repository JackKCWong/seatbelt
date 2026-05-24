package scanner

import (
	"regexp"
)

type PrivateKeyDetector struct {
	pattern *regexp.Regexp
}

func NewPrivateKeyDetector() *PrivateKeyDetector {
	return &PrivateKeyDetector{
		pattern: regexp.MustCompile(`-----BEGIN PRIVATE KEY-----`),
	}
}

func (d *PrivateKeyDetector) Name() string {
	return "private_key"
}

func (d *PrivateKeyDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding
	if d.pattern.MatchString(prompt) {
		findings = append(findings, Finding{
			Type:     "private_key",
			Value:    "-----BEGIN PRIVATE KEY-----",
			Location: "prompt",
		})
	}
	return findings, nil
}