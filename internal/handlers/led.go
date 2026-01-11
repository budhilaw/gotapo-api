package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// LEDHandler handles LED operations
type LEDHandler struct{}

// NewLEDHandler creates a new LED handler
func NewLEDHandler() *LEDHandler {
	return &LEDHandler{}
}

// SetLEDRequest represents an LED configuration request
type SetLEDRequest struct {
	Enabled bool `json:"enabled"`
}

// GetStatus gets LED status
// GET /api/cameras/:ip/led
func (h *LEDHandler) GetStatus(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getLedStatus", map[string]interface{}{
		"led": map[string]interface{}{
			"name": []string{"config"},
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "execution_failed",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  result,
	})
}

// SetStatus sets LED enabled state
// PUT /api/cameras/:ip/led
func (h *LEDHandler) SetStatus(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetLEDRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	enabled := "off"
	if req.Enabled {
		enabled = "on"
	}

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("setLedStatus", map[string]interface{}{
		"led": map[string]interface{}{
			"config": map[string]string{
				"enabled": enabled,
			},
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "execution_failed",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"result":  result,
	})
}
