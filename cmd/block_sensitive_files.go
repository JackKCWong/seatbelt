package cmd

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type ReadFileInput struct {
	FilePath  string `json:"filePath"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}

type HookInput struct {
	ToolName  string         `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

func readIgnorePatterns() ([]string, error) {
	data, err := os.ReadFile(".secretsignore")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}

func defaultSecretPatterns() []string {
	return []string{
		".env",
		".env.local",
		".env.*",
		".npmrc",
		".yarnrc",
		".pypirc",
		"pip.conf",
		"settings.xml",
		".gitcredentials",
		".netrc",
		"credentials",
		".aws/credentials",
		".aws/config",
		"docker/config.json",
		".docker/config.json",
		"secrets.yaml",
		"secrets.yml",
		"secrets.json",
		".vault-token",
		".燃料箱.toml",
		".sops.yaml",
		".sops.yml",
		".sops.json",
	}
}

func getAllPatterns() ([]string, error) {
	patterns := defaultSecretPatterns()
	userPatterns, err := readIgnorePatterns()
	if err != nil {
		return nil, err
	}
	patterns = append(patterns, userPatterns...)
	return patterns, nil
}

func matchesPattern(filePath string, patterns []string) (bool, string) {
	filePath = filepath.ToSlash(filePath)
	for _, pattern := range patterns {
		pattern = filepath.ToSlash(pattern)
		if ok, _ := filepath.Match(pattern, filePath); ok {
			return true, pattern
		}
		dir := filepath.Dir(filePath)
		base := filepath.Base(filePath)
		if ok, _ := filepath.Match(pattern, base); ok {
			return true, pattern
		}
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if ok, _ := filepath.Match(dirPattern, dir); ok {
				return true, pattern
			}
		}
	}
	return false, ""
}

func BlockSensitiveFilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-sensitive-files",
		Short: "Block sensitive files",
		RunE: func(cmd *cobra.Command, args []string) error {
			patterns, err := getAllPatterns()
			if err != nil {
				return err
			}

			input, _ := io.ReadAll(os.Stdin)
			if len(input) == 0 {
				payload := map[string]interface{}{
					"continue": true,
					"hookSpecificOutput": map[string]interface{}{
						"hookEventName":      "PreToolUse",
						"permissionDecision": "allow",
					},
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}

			var hookInput HookInput
			if err := json.Unmarshal(input, &hookInput); err != nil {
				payload := map[string]interface{}{
					"continue": true,
					"hookSpecificOutput": map[string]interface{}{
						"hookEventName":      "PreToolUse",
						"permissionDecision": "allow",
					},
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}

			if hookInput.ToolName == "read_file" {
				var readFileInput ReadFileInput
				if err := json.Unmarshal(hookInput.ToolInput, &readFileInput); err == nil {
					if matched, pattern := matchesPattern(readFileInput.FilePath, patterns); matched {
						payload := map[string]interface{}{
							"continue": true,
							"hookSpecificOutput": map[string]interface{}{
								"hookEventName":            "PreToolUse",
								"permissionDecision":       "deny",
								"permissionDecisionReason": "sensitive file read blocked by policy",
								"additionalContext": map[string]string{
									"file":   readFileInput.FilePath,
									"pattern": pattern,
								},
							},
						}
						enc := json.NewEncoder(os.Stdout)
						enc.SetIndent("", "  ")
						return enc.Encode(payload)
					}
				}
			}

			payload := map[string]interface{}{
				"continue": true,
				"hookSpecificOutput": map[string]interface{}{
					"hookEventName":      "PreToolUse",
					"permissionDecision": "allow",
				},
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(payload)
		},
	}
}