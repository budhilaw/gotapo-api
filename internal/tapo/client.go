package tapo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/budhilaw/gotapo-api/internal/crypto"
)

// NewClient creates a new Tapo camera client
func NewClient(host, username, password string) *Client {
	return &Client{
		Host:     host,
		Username: username,
		Password: password,
		Timeout:  10 * time.Second,
	}
}

// getHTTPClient returns an HTTP client with TLS verification disabled
func (c *Client) getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Tapo cameras use self-signed certs
			},
		},
	}
}

// getBaseURL returns the base URL for the camera
func (c *Client) getBaseURL() string {
	return fmt.Sprintf("https://%s:443", c.Host)
}

// getRequestURL returns the URL for authenticated requests
func (c *Client) getRequestURL() string {
	return fmt.Sprintf("%s/stok=%s/ds", c.getBaseURL(), c.stok)
}

// IsAuthenticated returns true if the client has a valid session
func (c *Client) IsAuthenticated() bool {
	return c.stok != ""
}

// GetSessionToken returns the current session token
func (c *Client) GetSessionToken() string {
	return c.stok
}

// Authenticate performs the full authentication flow
func (c *Client) Authenticate() error {
	// Try secure authentication first
	isSecure, err := c.detectConnectionType()
	if err != nil {
		return fmt.Errorf("failed to detect connection type: %w", err)
	}

	c.isSecure = isSecure

	if isSecure {
		return c.secureAuthenticate()
	}
	return c.legacyAuthenticate()
}

// detectConnectionType checks if the camera supports secure authentication
func (c *Client) detectConnectionType() (bool, error) {
	req := LoginRequest{
		Method: "login",
		Params: LoginParams{
			EncryptType: "3",
			Username:    c.Username,
		},
	}

	resp, err := c.makeRawRequest(req)
	if err != nil {
		return false, err
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		return false, fmt.Errorf("failed to parse response: %w", err)
	}

	// Error code -40413 with encrypt_type in result indicates secure connection
	if loginResp.ErrorCode == ErrorCodeLoginRequired {
		if loginResp.Result.Data != nil && len(loginResp.Result.Data.EncryptType) > 0 {
			for _, t := range loginResp.Result.Data.EncryptType {
				if t == "3" {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// secureAuthenticate performs the 3-phase secure authentication
func (c *Client) secureAuthenticate() error {
	// Phase 1: Request server nonce
	c.cnonce = crypto.GenerateCnonce()

	req := LoginRequest{
		Method: "login",
		Params: LoginParams{
			Cnonce:      c.cnonce,
			EncryptType: "3",
			Username:    c.Username,
		},
	}

	resp, err := c.makeRawRequest(req)
	if err != nil {
		return fmt.Errorf("phase 1 failed: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		return fmt.Errorf("failed to parse phase 1 response: %w", err)
	}

	if loginResp.ErrorCode != 0 || loginResp.Result.Data == nil {
		return NewTapoError(loginResp.ErrorCode, "failed to get server nonce")
	}

	c.nonce = loginResp.Result.Data.Nonce
	deviceConfirm := loginResp.Result.Data.DeviceConfirm

	// Phase 2: Validate device and generate keys
	// Try SHA256 first, then MD5
	hashedPassSHA := crypto.SHA256Hash(c.Password)
	hashedPassMD5 := crypto.MD5Hash(c.Password)

	var hashedPass string
	if c.validateDeviceConfirm(hashedPassSHA, deviceConfirm) {
		hashedPass = hashedPassSHA
	} else if c.validateDeviceConfirm(hashedPassMD5, deviceConfirm) {
		hashedPass = hashedPassMD5
	} else {
		return NewTapoError(ErrorCodeInvalidAuth, "device validation failed - check password")
	}

	c.hashedPass = hashedPass
	c.generateEncryptionKeys()

	// Phase 3: Complete login
	digestPasswd := crypto.SHA256Hash(hashedPass + c.cnonce + c.nonce)

	req = LoginRequest{
		Method: "login",
		Params: LoginParams{
			Cnonce:       c.cnonce,
			EncryptType:  "3",
			DigestPasswd: digestPasswd + c.cnonce + c.nonce,
			Username:     c.Username,
		},
	}

	resp, err = c.makeRawRequest(req)
	if err != nil {
		return fmt.Errorf("phase 3 failed: %w", err)
	}

	if err := json.Unmarshal(resp, &loginResp); err != nil {
		return fmt.Errorf("failed to parse phase 3 response: %w", err)
	}

	if loginResp.ErrorCode != 0 {
		return NewTapoError(loginResp.ErrorCode, ErrorMessage(loginResp.ErrorCode))
	}

	c.stok = loginResp.Result.Stok
	c.seq = loginResp.Result.StartSeq

	return nil
}

// validateDeviceConfirm validates the device_confirm value
func (c *Client) validateDeviceConfirm(hashedPass, deviceConfirm string) bool {
	expected := crypto.SHA256Hash(c.cnonce+hashedPass+c.nonce) + c.nonce + c.cnonce
	return deviceConfirm == expected
}

// generateEncryptionKeys generates the AES encryption keys (lsk and ivb)
func (c *Client) generateEncryptionKeys() {
	hashedKey := crypto.SHA256Hash(c.cnonce + c.hashedPass + c.nonce)

	// lsk = first 16 bytes of SHA256("lsk" + cnonce + nonce + hashedKey)
	lskFull := crypto.SHA256HashBytes("lsk" + c.cnonce + c.nonce + hashedKey)
	c.lsk = lskFull[:16]

	// ivb = first 16 bytes of SHA256("ivb" + cnonce + nonce + hashedKey)
	ivbFull := crypto.SHA256HashBytes("ivb" + c.cnonce + c.nonce + hashedKey)
	c.ivb = ivbFull[:16]
}

// legacyAuthenticate performs legacy MD5-based authentication
func (c *Client) legacyAuthenticate() error {
	hashedPass := crypto.MD5Hash(c.Password)

	req := LoginRequest{
		Method: "login",
		Params: LoginParams{
			Hashed:   true,
			Password: hashedPass,
			Username: c.Username,
		},
	}

	resp, err := c.makeRawRequest(req)
	if err != nil {
		return fmt.Errorf("legacy auth failed: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if loginResp.ErrorCode != 0 {
		return NewTapoError(loginResp.ErrorCode, ErrorMessage(loginResp.ErrorCode))
	}

	c.stok = loginResp.Result.Stok
	c.hashedPass = hashedPass

	return nil
}
