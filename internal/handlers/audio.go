package handlers

import (
	"strconv"

	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/budhilaw/gotapo-api/internal/tapo"
	"github.com/gofiber/fiber/v2"
)

// AudioHandler handles audio operations
type AudioHandler struct{}

// NewAudioHandler creates a new audio handler
func NewAudioHandler() *AudioHandler {
	return &AudioHandler{}
}

// SetSpeakerRequest represents a speaker volume request
type SetSpeakerRequest struct {
	Volume int `json:"volume"` // 0-100
}

// SetMicrophoneRequest represents a microphone configuration request
type SetMicrophoneRequest struct {
	Volume int  `json:"volume,omitempty"` // 0-100
	Mute   bool `json:"mute,omitempty"`
}

// GetConfig gets audio configuration
// GET /api/cameras/:ip/audio
func (h *AudioHandler) GetConfig(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "get",
		"audio_config": map[string]interface{}{
			"name": []string{"microphone", "speaker"},
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

// SetSpeaker sets speaker volume
// PUT /api/cameras/:ip/audio/speaker
func (h *AudioHandler) SetSpeaker(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetSpeakerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	// Validate volume range
	if req.Volume < 0 || req.Volume > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_volume",
			"message": "Volume must be between 0 and 100",
		})
	}

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "set",
		"audio_config": map[string]interface{}{
			"speaker": map[string]string{
				"volume": strconv.Itoa(req.Volume),
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

// SetMicrophone sets microphone configuration
// PUT /api/cameras/:ip/audio/microphone
func (h *AudioHandler) SetMicrophone(c *fiber.Ctx) error {
	cameraIP := c.Params("ip")
	username, password := middleware.GetTapoCredentials(c)

	var req SetMicrophoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "Invalid request body",
		})
	}

	micConfig := make(map[string]string)

	if req.Volume > 0 {
		if req.Volume > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_volume",
				"message": "Volume must be between 0 and 100",
			})
		}
		micConfig["volume"] = strconv.Itoa(req.Volume)
	}

	muteValue := "off"
	if req.Mute {
		muteValue = "on"
	}
	micConfig["mute"] = muteValue

	client := tapo.NewClient(cameraIP, username, password)

	payload := map[string]interface{}{
		"method": "set",
		"audio_config": map[string]interface{}{
			"microphone": micConfig,
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
