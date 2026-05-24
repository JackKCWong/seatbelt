package hook

import (
	"strings"
	"testing"

	"github.com/seatbelt/pkg/scanner"
)

func TestProcessInput(t *testing.T) {
	registry := scanner.NewRegistry(
		scanner.NewAPIKeyDetector(),
		scanner.NewPasswordDetector(),
	)

	tests := []struct {
		name       string
		input      HookInput
		wantCont   bool
		wantReason string
	}{
		{
			name:       "No secrets",
			input:     HookInput{Prompt: "Just a normal prompt"},
			wantCont:  true,
			wantReason: "",
		},
		{
			name:       "API key detected",
			input:     HookInput{Prompt: "Use AIzaSyABC123xyzDEFGHIJKLMNOPQRSTUV for API"},
			wantCont:  false,
			wantReason: "Security policy violation",
		},
		{
			name:       "Password detected",
			input:     HookInput{Prompt: "password=supersecret"},
			wantCont:  false,
			wantReason: "Security policy violation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ProcessInput(tt.input, registry)
			if output.Continue != tt.wantCont {
				t.Errorf("Continue = %v, want %v", output.Continue, tt.wantCont)
			}
			if !tt.wantCont && !strings.Contains(output.StopReason, tt.wantReason) {
				t.Errorf("StopReason = %v, want to contain %v", output.StopReason, tt.wantReason)
			}
		})
	}
}