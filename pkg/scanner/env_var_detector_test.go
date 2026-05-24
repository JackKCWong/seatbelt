package scanner

import (
	"testing"
)

func TestEnvVarDetector(t *testing.T) {
	t.Run("sensitive env var value in prompt", func(t *testing.T) {
		t.Setenv("MY_API_TOKEN", "abc123xyz-secret-token")
		t.Setenv("SECRET_KEY", "sk_12345")
		t.Setenv("DATABASE_PASSWORD", "super_secret_db_pass")

		d := NewEnvVarDetector()

		prompt := "Use the API token abc123xyz-secret-token for authentication"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) == 0 {
			t.Error("Expected to find MY_API_TOKEN value in prompt")
		}
	})

	t.Run("non-sensitive env var value in prompt", func(t *testing.T) {
		t.Setenv("HOME", "/Users/test")
		t.Setenv("USER", "testuser")

		d := NewEnvVarDetector()

		prompt := "Home is /Users/test"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) > 0 {
			t.Errorf("Expected no findings for non-sensitive env vars, got %d", len(findings))
		}
	})

	t.Run("no matching env var values", func(t *testing.T) {
		t.Setenv("API_KEY", "key_12345")

		d := NewEnvVarDetector()

		prompt := "Just a normal prompt without any secret values"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) > 0 {
			t.Errorf("Expected no findings, got %d", len(findings))
		}
	})

	t.Run("sensitive env var name detection", func(t *testing.T) {
		t.Setenv("STRIPE_TOKEN", "sk_live_abc123")

		d := NewEnvVarDetector()

		prompt := "Process payment with sk_live_abc123"
		findings, err := d.Detect(prompt)
		if err != nil {
			t.Fatalf("Detect() error = %v", err)
		}
		if len(findings) == 0 {
			t.Error("Expected to find STRIPE_TOKEN value in prompt")
		}
	})
}