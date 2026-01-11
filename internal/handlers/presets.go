package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// PresetsHandler handles preset operations
type PresetsHandler struct{}

// NewPresetsHandler creates a new presets handler
func NewPresetsHandler() *PresetsHandler {
	return &PresetsHandler{}
}

// CreatePresetRequest represents a create preset request
type CreatePresetRequest struct {
	Name string `json:"name"`
}

// List gets all presets
// GET /api/cameras/:ip/presets
func (h *PresetsHandler) List(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getPresetConfig", map[string]interface{}{
		"preset": map[string]interface{}{
			"name": []string{"preset"},
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

// Create saves current position as a preset
// POST /api/cameras/:ip/presets
func (h *PresetsHandler) Create(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req CreatePresetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_name",
			"message": "Preset name is required",
		})
	}

	client := tapo.NewClient(cameraIP, username, password)

	// Note: The API has a typo "addMotorPostion" - this is intentional
	result, err := client.Execute("addMotorPostion", map[string]interface{}{
		"preset": map[string]interface{}{
			"set_preset": map[string]string{
				"name":     req.Name,
				"save_ptz": "1",
			},
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "execution_failed",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"result":  result,
	})
}

// Goto moves camera to a preset position
// POST /api/cameras/:ip/presets/:id/goto
func (h *PresetsHandler) Goto(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	presetID := c.Params("id")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("motorMoveToPreset", map[string]interface{}{
		"preset": map[string]interface{}{
			"goto_preset": map[string]string{
				"id": presetID,
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

// Delete removes a preset
// DELETE /api/cameras/:ip/presets/:id
func (h *PresetsHandler) Delete(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	presetID := c.Params("id")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("deletePreset", map[string]interface{}{
		"preset": map[string]interface{}{
			"remove_preset": map[string][]string{
				"id": {presetID},
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
