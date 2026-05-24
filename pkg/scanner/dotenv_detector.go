package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type DotEnvDetector struct {
	envPaths []string
}

func NewDotEnvDetector() *DotEnvDetector {
	return &DotEnvDetector{
		envPaths: findEnvFiles(),
	}
}

func NewDotEnvDetectorWithPaths(paths []string) *DotEnvDetector {
	return &DotEnvDetector{
		envPaths: paths,
	}
}

func (d *DotEnvDetector) Name() string {
	return "dotenv"
}

func (d *DotEnvDetector) Detect(prompt string) ([]Finding, error) {
	var findings []Finding

	for _, path := range d.envPaths {
		secrets := readSecretsFromEnvFile(path)
		for name, value := range secrets {
			if isSensitiveName(name) && value != "" && strings.Contains(prompt, value) {
				findings = append(findings, Finding{
					Type:     "dotenv",
					Value:    name + "=" + value,
					Location: path,
				})
			}
		}
	}

	if len(findings) == 0 && containsEnvReference(prompt) {
		findings = append(findings, Finding{
			Type:     "dotenv",
			Value:    ".env",
			Location: "file reference",
		})
	}

	return findings, nil
}

func findEnvFiles() []string {
	cwd, _ := os.Getwd()
	if cwd != "" {
		return []string{filepath.Join(cwd, ".env")}
	}
	return nil
}

func readSecretsFromEnvFile(path string) map[string]string {
	secrets := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return secrets
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"")
			secrets[key] = value
		}
	}

	return secrets
}

func isSensitiveName(name string) bool {
	sensitivePatterns := []string{
		"TOKEN", "KEY", "SECRET", "PASSWORD", "PASS", "CREDENTIAL",
		"API_KEY", "AUTH", "PRIVATE", "ACCESS",
	}
	upper := strings.ToUpper(name)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

func containsEnvReference(prompt string) bool {
	return strings.Contains(strings.ToLower(prompt), ".env")
}