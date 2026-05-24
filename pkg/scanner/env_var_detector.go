package scanner

import (
	"os"
	"strings"
)

type EnvVarDetector struct {
	sensitiveNames []string
}

func NewEnvVarDetector() *EnvVarDetector {
	return &EnvVarDetector{
		sensitiveNames: []string{
			"TOKEN", "KEY", "SECRET", "PASSWORD", "PASS", "CREDENTIAL",
			"API_KEY", "AUTH", "PRIVATE", "ACCESS",
		},
	}
}

func (d *EnvVarDetector) Name() string {
	return "env_var"
}

func (d *EnvVarDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding

	for _, envName := range d.getSensitiveEnvVars() {
		envValue := os.Getenv(envName)
		if envValue != "" && strings.Contains(prompt, envValue) {
			findings = append(findings, Finding{
				Type:     "env_var",
				Value:    envName + "=" + envValue,
				Location: "environment variable " + envName,
			})
		}
	}

	return findings, nil
}

func (d *EnvVarDetector) getSensitiveEnvVars() []string {
	var sensitive []string
	for _, envName := range os.Environ() {
		name := strings.SplitN(envName, "=", 2)[0]
		if d.isSensitiveName(name) {
			sensitive = append(sensitive, name)
		}
	}
	return sensitive
}

func (d *EnvVarDetector) isSensitiveName(name string) bool {
	upper := strings.ToUpper(name)
	for _, pattern := range d.sensitiveNames {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}