package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// ImageHandler handles image and video settings operations
type ImageHandler struct{}

// NewImageHandler creates a new image handler
func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}

// SetFlipRequest represents an image flip request
type SetFlipRequest struct {
	FlipType string `json:"flip_type"` // "off", "center", etc.
}

// SetNightModeRequest represents a night mode request
type SetNightModeRequest struct {
	Mode string `json:"mode"` // "auto", "on" (night), "off" (day)
}

// GetSettings gets common image settings
// GET /api/cameras/:ip/image
func (h *ImageHandler) GetSettings(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getLdc", map[string]interface{}{
		"image": map[string]interface{}{
			"name": []string{"common", "switch"},
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

// SetFlip sets image flip mode
// PUT /api/cameras/:ip/image/flip
func (h *ImageHandler) SetFlip(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetFlipRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.FlipType == "" {
		req.FlipType = "off"
	}

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("setLdc", map[string]interface{}{
		"image": map[string]interface{}{
			"switch": map[string]string{
				"flip_type": req.FlipType,
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

// SetNightMode sets day/night mode
// PUT /api/cameras/:ip/image/nightmode
func (h *ImageHandler) SetNightMode(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetNightModeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	// Validate mode
	validModes := map[string]bool{"auto": true, "on": true, "off": true}
	if !validModes[req.Mode] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_mode",
			"message": "Mode must be 'auto', 'on' (night), or 'off' (day)",
		})
	}

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("setLdc", map[string]interface{}{
		"image": map[string]interface{}{
			"common": map[string]string{
				"inf_type": req.Mode,
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
