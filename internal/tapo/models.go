package tapo

import "time"

// Client represents a Tapo camera client
type Client struct {
	Host     string
	Username string
	Password string

	// Session data
	stok       string
	seq        int
	lsk        []byte // AES key (16 bytes)
	ivb        []byte // AES IV (16 bytes)
	cnonce     string
	nonce      string
	hashedPass string
	isSecure   bool

	// HTTP client timeout
	Timeout time.Duration
}

// LoginRequest represents the login API request
type LoginRequest struct {
	Method string      `json:"method"`
	Params LoginParams `json:"params"`
}

// LoginParams represents login parameters
type LoginParams struct {
	Cnonce       string `json:"cnonce,omitempty"`
	EncryptType  string `json:"encrypt_type,omitempty"`
	Username     string `json:"username"`
	DigestPasswd string `json:"digest_passwd,omitempty"`
	Hashed       bool   `json:"hashed,omitempty"`
	Password     string `json:"password,omitempty"`
}

// LoginResponse represents the login API response
type LoginResponse struct {
	ErrorCode int         `json:"error_code"`
	Result    LoginResult `json:"result"`
}

// LoginResult contains the login result data
type LoginResult struct {
	Data      *LoginData `json:"data,omitempty"`
	Stok      string     `json:"stok,omitempty"`
	StartSeq  int        `json:"start_seq,omitempty"`
	UserGroup string     `json:"user_group,omitempty"`
}

// LoginData contains nonce and device confirm for secure auth
type LoginData struct {
	EncryptType   []string `json:"encrypt_type,omitempty"`
	Nonce         string   `json:"nonce,omitempty"`
	DeviceConfirm string   `json:"device_confirm,omitempty"`
}

// MultipleRequest represents a batch request
type MultipleRequest struct {
	Method string            `json:"method"`
	Params MultipleReqParams `json:"params"`
}

// MultipleReqParams contains the requests array
type MultipleReqParams struct {
	Requests []SingleRequest `json:"requests"`
}

// SingleRequest represents a single API request
type SingleRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// DirectRequest represents a direct method request
type DirectRequest struct {
	Method string      `json:"method"`
	Data   interface{} `json:"-"` // Will be marshaled dynamically
}

// SecureRequest wraps encrypted payload
type SecureRequest struct {
	Method string       `json:"method"`
	Params SecureParams `json:"params"`
}

// SecureParams contains the encrypted request
type SecureParams struct {
	Request string `json:"request"`
}

// APIResponse is a generic API response
type APIResponse struct {
	ErrorCode int                    `json:"error_code"`
	Result    map[string]interface{} `json:"result,omitempty"`
}

// SecureResponse contains encrypted response data
type SecureResponse struct {
	ErrorCode int          `json:"error_code"`
	Result    SecureResult `json:"result,omitempty"`
}

// SecureResult contains the encrypted response
type SecureResult struct {
	Response string `json:"response"`
}

// ==== PTZ Models ====

// MotorMoveRequest for moving to coordinates
type MotorMoveRequest struct {
	Motor MotorMove `json:"motor"`
}

// MotorMove contains move coordinates
type MotorMove struct {
	Move *MoveCoords `json:"move,omitempty"`
	Step *MoveStep   `json:"movestep,omitempty"`
	Cali string      `json:"manual_cali,omitempty"`
}

// MoveCoords represents x/y coordinates
type MoveCoords struct {
	XCoord string `json:"x_coord"`
	YCoord string `json:"y_coord"`
}

// MoveStep represents directional movement
type MoveStep struct {
	Direction string `json:"direction"`
}

// CruiseRequest for cruise control
type CruiseRequest struct {
	Motor CruiseMotor `json:"motor"`
}

// CruiseMotor contains cruise settings
type CruiseMotor struct {
	Cruise     *CruiseCoord `json:"cruise,omitempty"`
	CruiseStop *struct{}    `json:"cruise_stop,omitempty"`
}

// CruiseCoord contains cruise coordinates
type CruiseCoord struct {
	Coord string `json:"coord"`
}

// ==== Preset Models ====

// PresetConfig represents preset configuration
type PresetConfig struct {
	Preset PresetData `json:"preset"`
}

// PresetData contains preset info
type PresetData struct {
	ID   []string `json:"id,omitempty"`
	Name []string `json:"name,omitempty"`
}

// GotoPresetRequest for going to a preset
type GotoPresetRequest struct {
	Preset GotoPreset `json:"preset"`
}

// GotoPreset contains the goto action
type GotoPreset struct {
	GotoPreset GotoPresetID `json:"goto_preset"`
}

// GotoPresetID contains the preset ID
type GotoPresetID struct {
	ID string `json:"id"`
}

// SetPresetRequest for saving a preset
type SetPresetRequest struct {
	Preset SetPreset `json:"preset"`
}

// SetPreset contains preset save data
type SetPreset struct {
	SetPreset SetPresetData `json:"set_preset"`
}

// SetPresetData contains preset name and save flag
type SetPresetData struct {
	Name    string `json:"name"`
	SavePTZ string `json:"save_ptz"`
}

// DeletePresetRequest for deleting a preset
type DeletePresetRequest struct {
	Preset DeletePreset `json:"preset"`
}

// DeletePreset contains preset deletion data
type DeletePreset struct {
	RemovePreset RemovePresetData `json:"remove_preset"`
}

// RemovePresetData contains IDs to remove
type RemovePresetData struct {
	ID []string `json:"id"`
}

// ==== Device Info Models ====

// DeviceInfoRequest for getting device info
type DeviceInfoRequest struct {
	DeviceInfo DeviceInfoParams `json:"device_info"`
}

// DeviceInfoParams contains the info query
type DeviceInfoParams struct {
	Name []string `json:"name"`
}

// DeviceInfo represents camera device information
type DeviceInfo struct {
	BasicInfo BasicInfo `json:"basic_info"`
}

// BasicInfo contains basic device information
type BasicInfo struct {
	DeviceType  string `json:"device_type"`
	DeviceModel string `json:"device_model"`
	DeviceName  string `json:"device_name"`
	DeviceInfo  string `json:"device_info"`
	HwVersion   string `json:"hw_version"`
	SwVersion   string `json:"sw_version"`
	DeviceAlias string `json:"device_alias"`
	Features    string `json:"features"`
	BarcodeMac  string `json:"barcode"`
	MAC         string `json:"mac"`
	DevID       string `json:"dev_id"`
	OEMID       string `json:"oem_id"`
	HwDesc      string `json:"hw_desc"`
}

// ==== Privacy Models ====

// LensMaskConfig represents lens mask settings
type LensMaskConfig struct {
	LensMask LensMaskInfo `json:"lens_mask"`
}

// LensMaskInfo contains lens mask enabled state
type LensMaskInfo struct {
	LensMaskInfo LensMaskEnabled `json:"lens_mask_info"`
}

// LensMaskEnabled contains enabled flag
type LensMaskEnabled struct {
	Enabled string `json:"enabled"`
}

// ==== Detection Models ====

// MotionDetectionConfig represents motion detection settings
type MotionDetectionConfig struct {
	MotionDetection MotionDetInfo `json:"motion_detection"`
}

// MotionDetInfo contains motion detection data
type MotionDetInfo struct {
	MotionDet MotionDetSettings `json:"motion_det"`
}

// MotionDetSettings contains motion detection settings
type MotionDetSettings struct {
	Enabled            string `json:"enabled"`
	DigitalSensitivity string `json:"digital_sensitivity,omitempty"`
}

// PersonDetectionConfig represents person detection settings
type PersonDetectionConfig struct {
	PeopleDetection PersonDetInfo `json:"people_detection"`
}

// PersonDetInfo contains person detection data
type PersonDetInfo struct {
	Detection PersonDetSettings `json:"detection"`
}

// PersonDetSettings contains person detection settings
type PersonDetSettings struct {
	Enabled     string `json:"enabled"`
	Sensitivity string `json:"sensitivity,omitempty"`
}

// ==== Alarm Models ====

// AlarmConfig for alarm settings
type AlarmConfig struct {
	MsgAlarm AlarmInfo `json:"msg_alarm"`
}

// AlarmInfo contains alarm settings
type AlarmInfo struct {
	Chn1AlarmInfo *Chn1AlarmSettings `json:"chn1_msg_alarm_info,omitempty"`
	ManualAlarm   *ManualAlarmAction `json:"manual_msg_alarm,omitempty"`
}

// Chn1AlarmSettings contains channel 1 alarm settings
type Chn1AlarmSettings struct {
	Enabled   string   `json:"enabled"`
	AlarmType string   `json:"alarm_type,omitempty"`
	LightType string   `json:"light_type,omitempty"`
	AlarmMode []string `json:"alarm_mode,omitempty"`
}

// ManualAlarmAction for manual alarm trigger
type ManualAlarmAction struct {
	Action string `json:"action"`
}

// ==== Image Settings Models ====

// ImageConfig for image settings
type ImageConfig struct {
	Image ImageSettings `json:"image"`
}

// ImageSettings contains image configuration
type ImageSettings struct {
	Common *CommonImageSettings `json:"common,omitempty"`
	Switch *SwitchSettings      `json:"switch,omitempty"`
	Name   interface{}          `json:"name,omitempty"`
}

// CommonImageSettings for day/night mode
type CommonImageSettings struct {
	InfType string `json:"inf_type,omitempty"`
}

// SwitchSettings for flip settings
type SwitchSettings struct {
	FlipType string `json:"flip_type,omitempty"`
}

// ==== LED Models ====

// LEDConfig for LED settings
type LEDConfig struct {
	LED LEDSettings `json:"led"`
}

// LEDSettings contains LED configuration
type LEDSettings struct {
	Config LEDEnabled `json:"config"`
}

// LEDEnabled contains enabled state
type LEDEnabled struct {
	Enabled string `json:"enabled"`
}

// ==== Audio Models ====

// AudioConfig for audio settings
type AudioConfig struct {
	AudioConfig AudioSettings `json:"audio_config"`
}

// AudioSettings contains audio configuration
type AudioSettings struct {
	Speaker    *SpeakerSettings    `json:"speaker,omitempty"`
	Microphone *MicrophoneSettings `json:"microphone,omitempty"`
	Name       []string            `json:"name,omitempty"`
}

// SpeakerSettings for speaker volume
type SpeakerSettings struct {
	Volume string `json:"volume"`
}

// MicrophoneSettings for microphone settings
type MicrophoneSettings struct {
	Volume string `json:"volume,omitempty"`
	Mute   string `json:"mute,omitempty"`
}

// ==== Recording Models ====

// RecordPlanConfig for recording settings
type RecordPlanConfig struct {
	RecordPlan RecordPlanInfo `json:"record_plan"`
}

// RecordPlanInfo contains recording plan data
type RecordPlanInfo struct {
	Name []string `json:"name,omitempty"`
}

// SDCardConfig for SD card operations
type SDCardConfig struct {
	HarddiskManage HarddiskInfo `json:"harddisk_manage"`
}

// HarddiskInfo contains SD card data
type HarddiskInfo struct {
	Table    []string `json:"table,omitempty"`
	FormatHD string   `json:"format_hd,omitempty"`
}

// ==== System Models ====

// SystemConfig for system operations
type SystemConfig struct {
	System SystemAction `json:"system"`
}

// SystemAction contains system action data
type SystemAction struct {
	Reboot string `json:"reboot,omitempty"`
	Name   string `json:"name,omitempty"`
}

// FirmwareConfig for firmware operations
type FirmwareConfig struct {
	CloudConfig FirmwareAction `json:"cloud_config"`
}

// FirmwareAction contains firmware action data
type FirmwareAction struct {
	CheckFWVersion string   `json:"check_fw_version,omitempty"`
	FWDownload     string   `json:"fw_download,omitempty"`
	Name           []string `json:"name,omitempty"`
}

// ==== Error Codes ====

const (
	ErrorCodeSuccess          = 0
	ErrorCodeInvalidToken     = -40401
	ErrorCodeRateLimited      = -40404
	ErrorCodeInvalidAuth      = -40411
	ErrorCodeLoginRequired    = -40413
	ErrorCodeCruiseInProgress = -64303
	ErrorCodeGeneral          = -1
)

// TapoError represents a Tapo API error
type TapoError struct {
	Code    int
	Message string
}

func (e *TapoError) Error() string {
	return e.Message
}

// NewTapoError creates a new TapoError
func NewTapoError(code int, message string) *TapoError {
	return &TapoError{Code: code, Message: message}
}

// ErrorMessage returns a human-readable error message for error codes
func ErrorMessage(code int) string {
	switch code {
	case ErrorCodeSuccess:
		return "Success"
	case ErrorCodeInvalidToken:
		return "Invalid or expired token"
	case ErrorCodeRateLimited:
		return "Rate limited - temporary suspension"
	case ErrorCodeInvalidAuth:
		return "Invalid authentication data"
	case ErrorCodeLoginRequired:
		return "Login required"
	case ErrorCodeCruiseInProgress:
		return "Cruise in progress - stop cruise first"
	case ErrorCodeGeneral:
		return "General error"
	default:
		return "Unknown error"
	}
}
