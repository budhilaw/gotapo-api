package crypto

import (
	"testing"
)

func TestMD5Hash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"password", "5F4DCC3B5AA765D61D8327DEB882CF99"},
		{"admin", "21232F297A57A5A743894A0E4A801FC3"},
		{"", "D41D8CD98F00B204E9800998ECF8427E"},
		{"test123", "CC03E747A6AFBBCBF8BE7668ACFEBEE5"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MD5Hash(tt.input)
			if result != tt.expected {
				t.Errorf("MD5Hash(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSHA256Hash(t *testing.T) {
	// Test that output is uppercase hex with correct length
	result := SHA256Hash("test")
	if len(result) != 64 {
		t.Errorf("SHA256Hash should return 64 character hex string, got %d", len(result))
	}

	// Test empty string
	emptyResult := SHA256Hash("")
	if emptyResult != "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855" {
		t.Errorf("SHA256Hash('') = %q, unexpected result", emptyResult)
	}

	// Test known hash
	passwordResult := SHA256Hash("password")
	if len(passwordResult) != 64 {
		t.Errorf("SHA256Hash('password') should be 64 chars, got %d", len(passwordResult))
	}
}

func TestSHA256HashBytes(t *testing.T) {
	result := SHA256HashBytes("test")
	if len(result) != 32 {
		t.Errorf("SHA256HashBytes should return 32 bytes, got %d", len(result))
	}
}

func TestGenerateCnonce(t *testing.T) {
	cnonce := GenerateCnonce()

	// Check length
	if len(cnonce) != 8 {
		t.Errorf("GenerateCnonce() should return 8 characters, got %d", len(cnonce))
	}

	// Check that it's uppercase hex
	for _, c := range cnonce {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			t.Errorf("GenerateCnonce() contains non-hex character: %c", c)
		}
	}

	// Check randomness (two calls should produce different results)
	cnonce2 := GenerateCnonce()
	if cnonce == cnonce2 {
		t.Error("GenerateCnonce() should produce random values")
	}
}
