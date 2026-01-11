package tapo

import (
	"testing"

	"github.com/budhilaw/gotapo-api/internal/crypto"
)

func TestNewClient(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")

	if client.Host != "192.168.1.100" {
		t.Errorf("Expected host 192.168.1.100, got %s", client.Host)
	}
	if client.Username != "admin" {
		t.Errorf("Expected username admin, got %s", client.Username)
	}
	if client.Password != "password" {
		t.Errorf("Expected password password, got %s", client.Password)
	}
	if client.Timeout.Seconds() != 10 {
		t.Errorf("Expected timeout 10s, got %v", client.Timeout)
	}
}

func TestClient_IsAuthenticated(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")

	if client.IsAuthenticated() {
		t.Error("New client should not be authenticated")
	}

	client.stok = "test_token"
	if !client.IsAuthenticated() {
		t.Error("Client with stok should be authenticated")
	}
}

func TestClient_GetSessionToken(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")

	if client.GetSessionToken() != "" {
		t.Error("New client should have empty token")
	}

	client.stok = "test_token"
	if client.GetSessionToken() != "test_token" {
		t.Error("Should return the stok value")
	}
}

func TestClient_getBaseURL(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")

	expected := "https://192.168.1.100:443"
	if client.getBaseURL() != expected {
		t.Errorf("Expected %s, got %s", expected, client.getBaseURL())
	}
}

func TestClient_getRequestURL(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")
	client.stok = "abc123"

	expected := "https://192.168.1.100:443/stok=abc123/ds"
	if client.getRequestURL() != expected {
		t.Errorf("Expected %s, got %s", expected, client.getRequestURL())
	}
}

func TestClient_validateDeviceConfirm(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")
	client.cnonce = "ABCD1234"
	client.nonce = "server_nonce"

	hashedPass := crypto.SHA256Hash("password")

	// Generate expected confirm
	expectedConfirm := crypto.SHA256Hash(client.cnonce+hashedPass+client.nonce) + client.nonce + client.cnonce

	if !client.validateDeviceConfirm(hashedPass, expectedConfirm) {
		t.Error("Should validate correct device confirm")
	}

	if client.validateDeviceConfirm(hashedPass, "wrong_confirm") {
		t.Error("Should not validate incorrect device confirm")
	}
}

func TestClient_generateEncryptionKeys(t *testing.T) {
	client := NewClient("192.168.1.100", "admin", "password")
	client.cnonce = "ABCD1234"
	client.nonce = "server_nonce"
	client.hashedPass = crypto.SHA256Hash("password")

	client.generateEncryptionKeys()

	if len(client.lsk) != 16 {
		t.Errorf("lsk should be 16 bytes, got %d", len(client.lsk))
	}
	if len(client.ivb) != 16 {
		t.Errorf("ivb should be 16 bytes, got %d", len(client.ivb))
	}
}

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{ErrorCodeSuccess, "Success"},
		{ErrorCodeInvalidToken, "Invalid or expired token"},
		{ErrorCodeRateLimited, "Rate limited - temporary suspension"},
		{ErrorCodeInvalidAuth, "Invalid authentication data"},
		{ErrorCodeLoginRequired, "Login required"},
		{ErrorCodeCruiseInProgress, "Cruise in progress - stop cruise first"},
		{ErrorCodeGeneral, "General error"},
		{-99999, "Unknown error"},
	}

	for _, tt := range tests {
		result := ErrorMessage(tt.code)
		if result != tt.expected {
			t.Errorf("ErrorMessage(%d) = %q, want %q", tt.code, result, tt.expected)
		}
	}
}

func TestTapoError(t *testing.T) {
	err := NewTapoError(-40401, "Invalid token")

	if err.Code != -40401 {
		t.Errorf("Expected code -40401, got %d", err.Code)
	}
	if err.Message != "Invalid token" {
		t.Errorf("Expected message 'Invalid token', got %s", err.Message)
	}
	if err.Error() != "Invalid token" {
		t.Errorf("Error() should return message")
	}
}
