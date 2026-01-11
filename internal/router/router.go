package router

import (
	"github.com/budhilaw/gotapo-api/internal/handlers"
	"github.com/budhilaw/gotapo-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// Setup configures all routes
func Setup(app *fiber.App) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API v1 routes
	api := app.Group("/api")

	// Camera routes - require credentials
	cameras := api.Group("/cameras/:ip", middleware.TapoCredentials())

	// Initialize handlers
	ptzHandler := handlers.NewPTZHandler()
	presetsHandler := handlers.NewPresetsHandler()
	deviceHandler := handlers.NewDeviceHandler()
	privacyHandler := handlers.NewPrivacyHandler()
	detectionHandler := handlers.NewDetectionHandler()
	alarmHandler := handlers.NewAlarmHandler()
	imageHandler := handlers.NewImageHandler()
	ledHandler := handlers.NewLEDHandler()
	audioHandler := handlers.NewAudioHandler()
	recordingHandler := handlers.NewRecordingHandler()
	systemHandler := handlers.NewSystemHandler()

	// PTZ routes
	ptz := cameras.Group("/ptz")
	ptz.Post("/move", ptzHandler.Move)
	ptz.Post("/step", ptzHandler.Step)
	ptz.Post("/calibrate", ptzHandler.Calibrate)
	ptz.Get("/capability", ptzHandler.GetCapability)
	ptz.Post("/cruise/start", ptzHandler.StartCruise)
	ptz.Post("/cruise/stop", ptzHandler.StopCruise)

	// Presets routes
	presets := cameras.Group("/presets")
	presets.Get("/", presetsHandler.List)
	presets.Post("/", presetsHandler.Create)
	presets.Post("/:id/goto", presetsHandler.Goto)
	presets.Delete("/:id", presetsHandler.Delete)

	// Device info routes
	cameras.Get("/info", deviceHandler.GetInfo)
	cameras.Get("/time", deviceHandler.GetTime)
	cameras.Get("/specs", deviceHandler.GetSpecs)

	// Privacy routes
	cameras.Get("/privacy", privacyHandler.GetPrivacy)
	cameras.Put("/privacy", privacyHandler.SetPrivacy)
	cameras.Get("/encryption", privacyHandler.GetEncryption)
	cameras.Put("/encryption", privacyHandler.SetEncryption)

	// Detection routes
	detection := cameras.Group("/detection")
	detection.Get("/motion", detectionHandler.GetMotionDetection)
	detection.Put("/motion", detectionHandler.SetMotionDetection)
	detection.Get("/person", detectionHandler.GetPersonDetection)
	detection.Put("/person", detectionHandler.SetPersonDetection)

	// Alarm routes
	cameras.Get("/alarm", alarmHandler.GetAlarm)
	cameras.Put("/alarm", alarmHandler.SetAlarm)
	cameras.Post("/alarm/trigger", alarmHandler.TriggerAlarm)
	cameras.Delete("/alarm/trigger", alarmHandler.StopAlarm)

	// Image settings routes
	cameras.Get("/image", imageHandler.GetSettings)
	cameras.Put("/image/flip", imageHandler.SetFlip)
	cameras.Put("/image/nightmode", imageHandler.SetNightMode)

	// LED routes
	cameras.Get("/led", ledHandler.GetStatus)
	cameras.Put("/led", ledHandler.SetStatus)

	// Audio routes
	cameras.Get("/audio", audioHandler.GetConfig)
	cameras.Put("/audio/speaker", audioHandler.SetSpeaker)
	cameras.Put("/audio/microphone", audioHandler.SetMicrophone)

	// Recording routes
	cameras.Get("/recording/plan", recordingHandler.GetRecordPlan)
	cameras.Get("/storage", recordingHandler.GetStorageStatus)
	cameras.Post("/storage/format", recordingHandler.FormatStorage)

	// System routes
	cameras.Post("/reboot", systemHandler.Reboot)
	cameras.Get("/firmware", systemHandler.GetFirmwareInfo)
	cameras.Post("/firmware/upgrade", systemHandler.StartFirmwareUpgrade)
}
