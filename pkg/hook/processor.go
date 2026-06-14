package hook

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/JackKCWong/seatbelt/pkg/scanner"
)

type HookInput struct {
	Timestamp     string `json:"timestamp"`
	CWD           string `json:"cwd"`
	SessionID     string `json:"sessionId"`
	HookEventName string `json:"hookEventName"`
	Prompt        string `json:"prompt,omitempty"`
}

type HookOutput struct {
	Continue      bool   `json:"continue"`
	StopReason    string `json:"stopReason,omitempty"`
	SystemMessage string `json:"systemMessage,omitempty"`
}

func ProcessInput(input HookInput, registry *scanner.Registry) HookOutput {
	findings, err := registry.Detect(input.Prompt)
	if err != nil {
		return HookOutput{
			Continue:      false,
			StopReason:    "Detection error",
			SystemMessage: fmt.Sprintf("Error scanning prompt: %v", err),
		}
	}

	if len(findings) > 0 {
		var locations []string
		seen := make(map[string]bool)
		for _, f := range findings {
			if !seen[f.Location] {
				seen[f.Location] = true
				locations = append(locations, f.Location)
			}
		}
		return HookOutput{
			Continue:      false,
			StopReason:    "Security policy violation",
			SystemMessage: fmt.Sprintf("Secret from %s in prompt detected", strings.Join(locations, ", ")),
		}
	}

	return HookOutput{Continue: true}
}

func ReadInput() (HookInput, error) {
	var input HookInput
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&input); err != nil {
		return HookInput{}, fmt.Errorf("decoding stdin: %w", err)
	}
	return input, nil
}

func WriteOutput(output HookOutput) error {
	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("encoding stdout: %w", err)
	}
	return nil
}
