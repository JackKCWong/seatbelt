package scanner

import (
	"os"
	"testing"
)

func TestDotEnvDetector(t *testing.T) {
	t.Run(".env file with secret patterns in prompt", func(t *testing.T) {
		tmpDir := t.TempDir()
		envFile := tmpDir + "/.env"
		err := os.WriteFile(envFile, []byte("STRIPE_TOKEN=sk_live_abc123\nAPI_KEY=key_12345\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		d := NewDotEnvDetectorWithPaths([]string{envFile})

		prompt := "Process payment with sk_live_abc123 and use key_12345"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) == 0 {
			t.Error("Expected to find .env secrets in prompt")
		}
	})

	t.Run("prompt references .env file", func(t *testing.T) {
		d := NewDotEnvDetector()

		prompt := "Check .env for configuration"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) == 0 {
			t.Error("Expected to find .env reference in prompt")
		}
	})

	t.Run("no .env secrets in prompt", func(t *testing.T) {
		tmpDir := t.TempDir()
		envFile := tmpDir + "/.env"
		err := os.WriteFile(envFile, []byte("STRIPE_TOKEN=sk_live_abc123\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create .env file: %v", err)
		}

		d := NewDotEnvDetectorWithPaths([]string{envFile})

		prompt := "Just a normal prompt without any secrets"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) > 0 {
			t.Errorf("Expected no findings, got %d", len(findings))
		}
	})
}