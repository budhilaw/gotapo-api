package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// SystemHandler handles system operations
type SystemHandler struct{}

// NewSystemHandler creates a new system handler
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// Reboot reboots the camera
// POST /api/cameras/:ip/reboot
func (h *SystemHandler) Reboot(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("rebootDevice", map[string]interface{}{
		"system": map[string]string{
			"reboot": "null",
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
		"message": "Camera reboot initiated",
		"result":  result,
	})
}

// GetFirmwareInfo checks for firmware updates
// GET /api/cameras/:ip/firmware
func (h *SystemHandler) GetFirmwareInfo(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	// Use multiple requests to check firmware and get upgrade info
	request := tapo.MultipleRequest{
		Method: "multipleRequest",
		Params: tapo.MultipleReqParams{
			Requests: []tapo.SingleRequest{
				{
					Method: "checkFirmwareVersionByCloud",
					Params: map[string]interface{}{
						"cloud_config": map[string]string{
							"check_fw_version": "null",
						},
					},
				},
				{
					Method: "getCloudConfig",
					Params: map[string]interface{}{
						"cloud_config": map[string]interface{}{
							"name": []string{"upgrade_info"},
						},
					},
				},
			},
		},
	}

	result, err := client.ExecuteDirect(request)
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

// StartFirmwareUpgrade starts firmware upgrade
// POST /api/cameras/:ip/firmware/upgrade
func (h *SystemHandler) StartFirmwareUpgrade(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"cloud_config": map[string]string{
			"fw_download": "null",
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
		"message": "Firmware upgrade started",
		"result":  result,
	})
}
