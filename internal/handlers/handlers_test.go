package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// mockAuthMiddleware adds mock credentials to context
func mockAuthMiddleware(c *fiber.Ctx) error {
	c.Locals("tapo_username", "test_user")
	c.Locals("tapo_password", "test_pass")
	return c.Next()
}

func TestPTZHandler_Step_InvalidBody(t *testing.T) {
	app := fiber.New()
	handler := NewPTZHandler()

	app.Post("/cameras/:ip/ptz/step", mockAuthMiddleware, handler.Step)

	req := httptest.NewRequest("POST", "/cameras/192.168.1.100/ptz/step",
		bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestPTZHandler_Step_InvalidDirection(t *testing.T) {
	app := fiber.New()
	handler := NewPTZHandler()

	app.Post("/cameras/:ip/ptz/step", mockAuthMiddleware, handler.Step)

	body, _ := json.Marshal(StepRequest{Direction: 400}) // > 359
	req := httptest.NewRequest("POST", "/cameras/192.168.1.100/ptz/step",
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid direction, got %d", resp.StatusCode)
	}
}

func TestPresetsHandler_Create_EmptyName(t *testing.T) {
	app := fiber.New()
	handler := NewPresetsHandler()

	app.Post("/cameras/:ip/presets", mockAuthMiddleware, handler.Create)

	body, _ := json.Marshal(CreatePresetRequest{Name: ""})
	req := httptest.NewRequest("POST", "/cameras/192.168.1.100/presets",
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400 for empty name, got %d", resp.StatusCode)
	}
}

func TestImageHandler_SetNightMode_InvalidMode(t *testing.T) {
	app := fiber.New()
	handler := NewImageHandler()

	app.Put("/cameras/:ip/image/nightmode", mockAuthMiddleware, handler.SetNightMode)

	body, _ := json.Marshal(SetNightModeRequest{Mode: "invalid"})
	req := httptest.NewRequest("PUT", "/cameras/192.168.1.100/image/nightmode",
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid mode, got %d", resp.StatusCode)
	}
}

func TestAudioHandler_SetSpeaker_InvalidVolume(t *testing.T) {
	app := fiber.New()
	handler := NewAudioHandler()

	app.Put("/cameras/:ip/audio/speaker", mockAuthMiddleware, handler.SetSpeaker)

	body, _ := json.Marshal(SetSpeakerRequest{Volume: 150}) // > 100
	req := httptest.NewRequest("PUT", "/cameras/192.168.1.100/audio/speaker",
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid volume, got %d", resp.StatusCode)
	}

	// Check response body
	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	if result["error"] != "invalid_volume" {
		t.Errorf("Expected error=invalid_volume, got %v", result["error"])
	}
}
