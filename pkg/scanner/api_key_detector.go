package scanner

import (
	"regexp"
)

type APIKeyDetector struct {
	patterns []*regexp.Regexp
}

func NewAPIKeyDetector() *APIKeyDetector {
	return &APIKeyDetector{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(AIzaSy[a-zA-Z0-9_-]{20,})`),
			regexp.MustCompile(`(?i)(sk_live_[a-zA-Z0-9]{20,})`),
			regexp.MustCompile(`(?i)(xox[baprs]-[a-zA-Z0-9_-]{10,})`),
			regexp.MustCompile(`(?i)(ghp_[a-zA-Z0-9]{20,})`),
			regexp.MustCompile(`(?i)(glpat-[a-zA-Z0-9_-]{10,})`),
		},
	}
}

func (d *APIKeyDetector) Name() string {
	return "api_key"
}

func (d *APIKeyDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding
	for _, pattern := range d.patterns {
		matches := pattern.FindAllString(prompt, -1)
		for _, match := range matches {
			findings = append(findings, Finding{
				Type:     "api_key",
				Value:    match,
				Location: "prompt",
			})
		}
	}
	return findings, nil
}