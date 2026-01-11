package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// DetectionHandler handles detection configuration operations
type DetectionHandler struct{}

// NewDetectionHandler creates a new detection handler
func NewDetectionHandler() *DetectionHandler {
	return &DetectionHandler{}
}

// SetMotionDetectionRequest represents a motion detection configuration request
type SetMotionDetectionRequest struct {
	Enabled     bool `json:"enabled"`
	Sensitivity int  `json:"sensitivity,omitempty"` // 0-100
}

// SetPersonDetectionRequest represents a person detection configuration request
type SetPersonDetectionRequest struct {
	Enabled     bool `json:"enabled"`
	Sensitivity int  `json:"sensitivity,omitempty"` // 0-100
}

// GetMotionDetection gets motion detection configuration
// GET /api/cameras/:ip/detection/motion
func (h *DetectionHandler) GetMotionDetection(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getDetectionConfig", map[string]interface{}{
		"motion_detection": map[string]interface{}{
			"name": []string{"motion_det"},
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

// SetMotionDetection sets motion detection configuration
// PUT /api/cameras/:ip/detection/motion
func (h *DetectionHandler) SetMotionDetection(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetMotionDetectionRequest
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

	params := map[string]interface{}{
		"motion_detection": map[string]interface{}{
			"motion_det": map[string]interface{}{
				"enabled": enabled,
			},
		},
	}

	// Add sensitivity if specified
	if req.Sensitivity > 0 {
		params["motion_detection"].(map[string]interface{})["motion_det"].(map[string]interface{})["digital_sensitivity"] = formatNumber(req.Sensitivity)
	}

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("setDetectionConfig", params)
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

// GetPersonDetection gets person detection configuration
// GET /api/cameras/:ip/detection/person
func (h *DetectionHandler) GetPersonDetection(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getPersonDetectionConfig", map[string]interface{}{
		"people_detection": map[string]interface{}{
			"name": []string{"detection"},
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

// SetPersonDetection sets person detection configuration
// PUT /api/cameras/:ip/detection/person
func (h *DetectionHandler) SetPersonDetection(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetPersonDetectionRequest
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

	params := map[string]interface{}{
		"people_detection": map[string]interface{}{
			"detection": map[string]interface{}{
				"enabled": enabled,
			},
		},
	}

	// Add sensitivity if specified
	if req.Sensitivity > 0 {
		params["people_detection"].(map[string]interface{})["detection"].(map[string]interface{})["sensitivity"] = formatNumber(req.Sensitivity)
	}

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("setPersonDetectionConfig", params)
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

// formatNumber converts an int to a string
func formatNumber(n int) string {
	return string([]byte{
		byte('0' + n/100),
		byte('0' + (n%100)/10),
		byte('0' + n%10),
	})
}
