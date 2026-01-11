package tapo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/budhilaw/gotapo-api/internal/crypto"
)

// defaultHeaders returns the default HTTP headers for Tapo API requests
func (c *Client) defaultHeaders() map[string]string {
	return map[string]string{
		"Host":            fmt.Sprintf("%s:443", c.Host),
		"Referer":         fmt.Sprintf("https://%s", c.Host),
		"Accept":          "application/json",
		"Accept-Encoding": "gzip, deflate",
		"User-Agent":      "Tapo CameraClient Android",
		"Connection":      "close",
		"requestByApp":    "true",
		"Content-Type":    "application/json; charset=UTF-8",
	}
}

// makeRawRequest makes a raw HTTP POST request to the camera
func (c *Client) makeRawRequest(payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.getBaseURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.defaultHeaders() {
		req.Header.Set(key, value)
	}

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

// Execute sends a command to the camera
func (c *Client) Execute(method string, params interface{}) (map[string]interface{}, error) {
	if !c.IsAuthenticated() {
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	request := SingleRequest{
		Method: method,
		Params: params,
	}

	multiReq := MultipleRequest{
		Method: "multipleRequest",
		Params: MultipleReqParams{
			Requests: []SingleRequest{request},
		},
	}

	if c.isSecure {
		return c.executeSecure(multiReq)
	}
	return c.executePlain(multiReq)
}

// ExecuteDirect sends a direct command without multipleRequest wrapper
func (c *Client) ExecuteDirect(payload interface{}) (map[string]interface{}, error) {
	if !c.IsAuthenticated() {
		if err := c.Authenticate(); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	if c.isSecure {
		return c.executeSecure(payload)
	}
	return c.executePlain(payload)
}

// executePlain sends an unencrypted request
func (c *Client) executePlain(payload interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.getRequestURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.defaultHeaders() {
		req.Header.Set(key, value)
	}

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResp.ErrorCode != 0 {
		// Check if we need to re-authenticate
		if apiResp.ErrorCode == ErrorCodeInvalidToken {
			c.stok = ""
			return nil, NewTapoError(apiResp.ErrorCode, "session expired - please re-authenticate")
		}
		return nil, NewTapoError(apiResp.ErrorCode, ErrorMessage(apiResp.ErrorCode))
	}

	return apiResp.Result, nil
}

// executeSecure sends an encrypted request
func (c *Client) executeSecure(payload interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Encrypt the payload
	encrypted, err := crypto.AESEncrypt(jsonData, c.lsk, c.ivb)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt request: %w", err)
	}

	// Create secure request
	secureReq := SecureRequest{
		Method: "securePassthrough",
		Params: SecureParams{
			Request: base64.StdEncoding.EncodeToString(encrypted),
		},
	}

	secureJSON, err := json.Marshal(secureReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal secure request: %w", err)
	}

	req, err := http.NewRequest("POST", c.getRequestURL(), bytes.NewBuffer(secureJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.defaultHeaders() {
		req.Header.Set(key, value)
	}

	// Add secure headers
	req.Header.Set("Seq", strconv.Itoa(c.seq))
	req.Header.Set("Tapo_tag", c.calculateTag(string(secureJSON)))

	// Increment sequence number
	c.seq++

	client := c.getHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var secureResp SecureResponse
	if err := json.Unmarshal(body, &secureResp); err != nil {
		return nil, fmt.Errorf("failed to parse secure response: %w", err)
	}

	if secureResp.ErrorCode != 0 {
		if secureResp.ErrorCode == ErrorCodeInvalidToken {
			c.stok = ""
			return nil, NewTapoError(secureResp.ErrorCode, "session expired - please re-authenticate")
		}
		return nil, NewTapoError(secureResp.ErrorCode, ErrorMessage(secureResp.ErrorCode))
	}

	// Decrypt response
	encryptedResp, err := base64.StdEncoding.DecodeString(secureResp.Result.Response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	decrypted, err := crypto.AESDecrypt(encryptedResp, c.lsk, c.ivb)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt response: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(decrypted, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted response: %w", err)
	}

	if apiResp.ErrorCode != 0 {
		return nil, NewTapoError(apiResp.ErrorCode, ErrorMessage(apiResp.ErrorCode))
	}

	return apiResp.Result, nil
}

// calculateTag calculates the Tapo_tag header value
func (c *Client) calculateTag(requestJSON string) string {
	tag1 := crypto.SHA256Hash(c.hashedPass + c.cnonce)
	tag := crypto.SHA256Hash(tag1 + requestJSON + strconv.Itoa(c.seq))
	return tag
}
