package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// PrivacyHandler handles privacy and security operations
type PrivacyHandler struct{}

// NewPrivacyHandler creates a new privacy handler
func NewPrivacyHandler() *PrivacyHandler {
	return &PrivacyHandler{}
}

// SetPrivacyRequest represents a privacy mode request
type SetPrivacyRequest struct {
	Enabled bool `json:"enabled"`
}

// SetEncryptionRequest represents a media encryption request
type SetEncryptionRequest struct {
	Enabled bool `json:"enabled"`
}

// GetPrivacy gets lens mask status
// GET /api/cameras/:ip/privacy
func (h *PrivacyHandler) GetPrivacy(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getLensMaskConfig", map[string]interface{}{
		"lens_mask": map[string]interface{}{
			"name": []string{"lens_mask_info"},
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

// SetPrivacy sets lens mask (privacy mode)
// PUT /api/cameras/:ip/privacy
func (h *PrivacyHandler) SetPrivacy(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetPrivacyRequest
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

	result, err := client.Execute("setLensMaskConfig", map[string]interface{}{
		"lens_mask": map[string]interface{}{
			"lens_mask_info": map[string]string{
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

// GetEncryption gets media encryption status
// GET /api/cameras/:ip/encryption
func (h *PrivacyHandler) GetEncryption(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getMediaEncrypt", map[string]interface{}{
		"cet": map[string]interface{}{
			"name": []string{"media_encrypt"},
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

// SetEncryption sets media encryption
// PUT /api/cameras/:ip/encryption
func (h *PrivacyHandler) SetEncryption(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetEncryptionRequest
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

	result, err := client.Execute("setMediaEncrypt", map[string]interface{}{
		"cet": map[string]interface{}{
			"media_encrypt": map[string]string{
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
