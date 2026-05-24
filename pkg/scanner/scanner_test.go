package scanner

import (
	"testing"
)

func TestAPIKeyDetector(t *testing.T) {
	d := NewAPIKeyDetector()
	tests := []struct {
		name     string
		prompt   string
		wantFind bool
	}{
		{"AIzaSy key", "Use AIzaSyABC123xyzDEFGHIJKLMNOPQRSTUV for API", true},
		{"Stripe live key", "sk_live_abc123XYZ456789012345678901234", true},
		{"GitHub token", "ghp_abcdefghijklmnopqrstuvwxyz1234567890", true},
		{"No secret", "Just a normal prompt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, _ := d.Detect(tt.prompt)
			found := len(findings) > 0
			if found != tt.wantFind {
				t.Errorf("Detect() found=%v, want %v", found, tt.wantFind)
			}
		})
	}
}

func TestDBConnectionDetector(t *testing.T) {
	d := NewDBConnectionDetector()
	tests := []struct {
		name     string
		prompt   string
		wantFind bool
	}{
		{"MongoDB", "mongodb+srv://user:pass@cluster.example.com", true},
		{"Postgres", "postgres://user:pass@localhost:5432/db", true},
		{"MySQL", "mysql://user:pass@localhost:3306/db", true},
		{"No secret", "SELECT * FROM users", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, _ := d.Detect(tt.prompt)
			found := len(findings) > 0
			if found != tt.wantFind {
				t.Errorf("Detect() found=%v, want %v", found, tt.wantFind)
			}
		})
	}
}

func TestPrivateKeyDetector(t *testing.T) {
	d := NewPrivateKeyDetector()
	tests := []struct {
		name     string
		prompt   string
		wantFind bool
	}{
		{"Private key", "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----", true},
		{"No key", "Just some text", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, _ := d.Detect(tt.prompt)
			found := len(findings) > 0
			if found != tt.wantFind {
				t.Errorf("Detect() found=%v, want %v", found, tt.wantFind)
			}
		})
	}
}

func TestPasswordDetector(t *testing.T) {
	d := NewPasswordDetector()
	tests := []struct {
		name     string
		prompt   string
		wantFind bool
	}{
		{"password=", "password=supersecret", true},
		{"pwd=", "pwd=secret123", true},
		{"passwd=", "passwd=my-password", true},
		{"No password", "user: admin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, _ := d.Detect(tt.prompt)
			found := len(findings) > 0
			if found != tt.wantFind {
				t.Errorf("Detect() found=%v, want %v", found, tt.wantFind)
			}
		})
	}
}

func TestDotEnvDetector(t *testing.T) {
	d := NewDotEnvDetector()
	tests := []struct {
		name     string
		prompt   string
		wantFind bool
	}{
		{".env file", "Check .env for configuration", true},
		{"No file", "Just some text", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings, _ := d.Detect(tt.prompt)
			found := len(findings) > 0
			if found != tt.wantFind {
				t.Errorf("Detect() found=%v, want %v", found, tt.wantFind)
			}
		})
	}
}

func TestRegistry(t *testing.T) {
	r := NewRegistry(
		NewAPIKeyDetector(),
		NewPasswordDetector(),
	)

	findings, err := r.Detect("password=secret and AIzaSyABC123xyzDEFGHIJKLMNOPQRSTUV")
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}
	if len(findings) != 2 {
		t.Errorf("Detect() found %d, want 2", len(findings))
	}
}