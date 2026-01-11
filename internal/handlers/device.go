package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// DeviceHandler handles device information operations
type DeviceHandler struct{}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

// GetInfo gets device basic information
// GET /api/cameras/:ip/info
func (h *DeviceHandler) GetInfo(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getDeviceInfo", map[string]interface{}{
		"device_info": map[string]interface{}{
			"name": []string{"basic_info"},
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

// GetTime gets device clock status
// GET /api/cameras/:ip/time
func (h *DeviceHandler) GetTime(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getClockStatus", map[string]interface{}{
		"system": map[string]string{
			"name": "clock_status",
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

// GetSpecs gets module specifications
// GET /api/cameras/:ip/specs
func (h *DeviceHandler) GetSpecs(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "get",
		"function": map[string]interface{}{
			"name": []string{"module_spec"},
		},
	}

	result, err := client.ExecuteDirect(payload)
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
