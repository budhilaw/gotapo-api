package handlers

import (
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// PTZHandler handles PTZ (Pan-Tilt-Zoom) operations
type PTZHandler struct{}

// NewPTZHandler creates a new PTZ handler
func NewPTZHandler() *PTZHandler {
	return &PTZHandler{}
}

// MoveRequest represents a move to coordinates request
type MoveRequest struct {
	XCoord string `json:"x_coord"`
	YCoord string `json:"y_coord"`
}

// StepRequest represents a directional step request
type StepRequest struct {
	Direction int `json:"direction"` // 0=right, 90=up, 180=left, 270=down
}

// Move moves the camera to specific coordinates
// POST /api/cameras/:ip/ptz/move
func (h *PTZHandler) Move(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req MoveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"motor": map[string]interface{}{
			"move": map[string]string{
				"x_coord": req.XCoord,
				"y_coord": req.YCoord,
			},
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

// Step moves the camera in a direction
// POST /api/cameras/:ip/ptz/step
func (h *PTZHandler) Step(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req StepRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	// Validate direction (0-359)
	if req.Direction < 0 || req.Direction > 359 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_direction",
			"message": "Direction must be between 0 and 359",
		})
	}

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"motor": map[string]interface{}{
			"movestep": map[string]string{
				"direction": string(rune('0'+req.Direction/100)) + string(rune('0'+(req.Direction%100)/10)) + string(rune('0'+req.Direction%10)),
			},
		},
	}

	// Properly format direction as string
	payload["motor"].(map[string]interface{})["movestep"].(map[string]string)["direction"] = formatDirection(req.Direction)

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

// formatDirection formats direction as string
func formatDirection(d int) string {
	return string([]byte{
		byte('0' + d/100),
		byte('0' + (d%100)/10),
		byte('0' + d%10),
	})
}

// Calibrate starts motor calibration
// POST /api/cameras/:ip/ptz/calibrate
func (h *PTZHandler) Calibrate(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"motor": map[string]interface{}{
			"manual_cali": "",
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

// GetCapability gets motor capability info
// GET /api/cameras/:ip/ptz/capability
func (h *PTZHandler) GetCapability(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "get",
		"motor": map[string]interface{}{
			"name": []string{"capability"},
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

// StartCruise starts cruise/patrol mode
// POST /api/cameras/:ip/ptz/cruise/start
func (h *PTZHandler) StartCruise(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"motor": map[string]interface{}{
			"cruise": map[string]string{
				"coord": "0",
			},
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

// StopCruise stops cruise/patrol mode
// POST /api/cameras/:ip/ptz/cruise/stop
func (h *PTZHandler) StopCruise(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"motor": map[string]interface{}{
			"cruise_stop": map[string]interface{}{},
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
