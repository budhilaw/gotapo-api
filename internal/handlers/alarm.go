package handlers

import (
	"strconv"

	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// AlarmHandler handles alarm operations
type AlarmHandler struct{}

// NewAlarmHandler creates a new alarm handler
func NewAlarmHandler() *AlarmHandler {
	return &AlarmHandler{}
}

// SetAlarmRequest represents an alarm configuration request
type SetAlarmRequest struct {
	Enabled   bool     `json:"enabled"`
	AlarmType string   `json:"alarm_type,omitempty"` // "0" for default
	LightType string   `json:"light_type,omitempty"` // "0" for default
	AlarmMode []string `json:"alarm_mode,omitempty"` // ["sound", "light"]
}

// GetAlarm gets alarm status
// GET /api/cameras/:ip/alarm
func (h *AlarmHandler) GetAlarm(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	result, err := client.Execute("getLastAlarmInfo", map[string]interface{}{
		"msg_alarm": map[string]interface{}{
			"name": []string{"chn1_msg_alarm_info"},
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

// SetAlarm sets alarm configuration
// PUT /api/cameras/:ip/alarm
func (h *AlarmHandler) SetAlarm(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetAlarmRequest
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

	alarmConfig := map[string]interface{}{
		"enabled": enabled,
	}

	if req.AlarmType != "" {
		alarmConfig["alarm_type"] = req.AlarmType
	} else {
		alarmConfig["alarm_type"] = "0"
	}

	if req.LightType != "" {
		alarmConfig["light_type"] = req.LightType
	} else {
		alarmConfig["light_type"] = "0"
	}

	if len(req.AlarmMode) > 0 {
		alarmConfig["alarm_mode"] = req.AlarmMode
	} else {
		alarmConfig["alarm_mode"] = []string{"sound", "light"}
	}

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "set",
		"msg_alarm": map[string]interface{}{
			"chn1_msg_alarm_info": alarmConfig,
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

// TriggerAlarm starts manual alarm
// POST /api/cameras/:ip/alarm/trigger
func (h *AlarmHandler) TriggerAlarm(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"msg_alarm": map[string]interface{}{
			"manual_msg_alarm": map[string]string{
				"action": "start",
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

// StopAlarm stops manual alarm
// DELETE /api/cameras/:ip/alarm/trigger
func (h *AlarmHandler) StopAlarm(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "do",
		"msg_alarm": map[string]interface{}{
			"manual_msg_alarm": map[string]string{
				"action": "stop",
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

// Helper to use strconv package
var _ = strconv.Itoa
