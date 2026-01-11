package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// RecordingHandler handles recording and storage operations
type RecordingHandler struct{}

// NewRecordingHandler creates a new recording handler
func NewRecordingHandler() *RecordingHandler {
	return &RecordingHandler{}
}

// GetRecordPlan gets recording plan configuration
// GET /api/cameras/:ip/recording/plan
func (h *RecordingHandler) GetRecordPlan(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getRecordPlan", map[string]interface{}{
		"record_plan": map[string]interface{}{
			"name": []string{"chn1_channel"},
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

// GetStorageStatus gets SD card status
// GET /api/cameras/:ip/storage
func (h *RecordingHandler) GetStorageStatus(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getSdCardStatus", map[string]interface{}{
		"harddisk_manage": map[string]interface{}{
			"table": []string{"hd_info"},
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

// FormatStorage formats SD card
// POST /api/cameras/:ip/storage/format
func (h *RecordingHandler) FormatStorage(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("formatSdCard", map[string]interface{}{
		"harddisk_manage": map[string]string{
			"format_hd": "1",
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
		"message": "SD card format initiated",
		"result":  result,
	})
}
