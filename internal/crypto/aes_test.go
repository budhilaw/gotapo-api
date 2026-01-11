package crypto

import (
	"bytes"
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key := []byte("1234567890123456") // 16 bytes
	iv := []byte("abcdefghijklmnop")  // 16 bytes
	plaintext := []byte("Hello, Tapo Camera!")

	// Encrypt
	ciphertext, err := AESEncrypt(plaintext, key, iv)
	if err != nil {
		t.Fatalf("AESEncrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := AESDecrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatalf("AESDecrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original. Got %q, want %q", decrypted, plaintext)
	}
}

func TestAESEncryptDecryptEmptyInput(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("abcdefghijklmnop")
	plaintext := []byte("")

	ciphertext, err := AESEncrypt(plaintext, key, iv)
	if err != nil {
		t.Fatalf("AESEncrypt failed for empty input: %v", err)
	}

	decrypted, err := AESDecrypt(ciphertext, key, iv)
	if err != nil {
		t.Fatalf("AESDecrypt failed for empty input: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted empty text doesn't match. Got %q, want %q", decrypted, plaintext)
	}
}

func TestAESEncryptInvalidKeyLength(t *testing.T) {
	key := []byte("short") // Too short
	iv := []byte("abcdefghijklmnop")
	plaintext := []byte("test")

	_, err := AESEncrypt(plaintext, key, iv)
	if err == nil {
		t.Error("AESEncrypt should fail with invalid key length")
	}
}

func TestAESEncryptInvalidIVLength(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("short") // Too short
	plaintext := []byte("test")

	_, err := AESEncrypt(plaintext, key, iv)
	if err == nil {
		t.Error("AESEncrypt should fail with invalid IV length")
	}
}

func TestAESDecryptInvalidCiphertext(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("abcdefghijklmnop")
	ciphertext := []byte("not a multiple") // Not a multiple of block size

	_, err := AESDecrypt(ciphertext, key, iv)
	if err == nil {
		t.Error("AESDecrypt should fail with invalid ciphertext length")
	}
}

func TestPKCS7Padding(t *testing.T) {
	tests := []struct {
		inputLen    int
		expectedLen int
	}{
		{0, 16},  // Empty input gets full block of padding
		{1, 16},  // 1 byte + 15 padding = 16
		{15, 16}, // 15 bytes + 1 padding = 16
		{16, 32}, // 16 bytes + 16 padding = 32
		{17, 32}, // 17 bytes + 15 padding = 32
	}

	for _, tt := range tests {
		input := bytes.Repeat([]byte("a"), tt.inputLen)
		padded := pkcs7Pad(input, 16)
		if len(padded) != tt.expectedLen {
			t.Errorf("pkcs7Pad(len=%d) = len %d, want %d", tt.inputLen, len(padded), tt.expectedLen)
		}
	}
}

func TestPKCS7UnpadInvalidPadding(t *testing.T) {
	_, err := pkcs7Unpad([]byte{})
	if err == nil {
		t.Error("pkcs7Unpad should fail with empty data")
	}

	// Invalid padding value (0)
	_, err = pkcs7Unpad([]byte{0})
	if err == nil {
		t.Error("pkcs7Unpad should fail with padding value 0")
	}

	// Padding value larger than data
	_, err = pkcs7Unpad([]byte{20})
	if err == nil {
		t.Error("pkcs7Unpad should fail with padding larger than data")
	}
}
