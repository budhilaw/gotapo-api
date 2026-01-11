# Tapo Camera HTTP API Reference

This document describes the raw HTTP API protocol used by TP-Link Tapo cameras. Use this reference to implement your own library in any programming language.

---

## Table of Contents

1. [Overview](#overview)
2. [Connection Setup](#connection-setup)
3. [Authentication](#authentication)
4. [Making Requests](#making-requests)
5. [Encryption](#encryption)
6. [API Commands Reference](#api-commands-reference)
   - [PTZ (Pan-Tilt-Zoom)](#ptz-pan-tilt-zoom)
   - [Presets](#presets)
   - [Device Info](#device-info)
   - [Privacy & Security](#privacy--security)
   - [Detection & Alarms](#detection--alarms)
   - [Image & Video Settings](#image--video-settings)
   - [LED & Lights](#led--lights)
   - [Audio](#audio)
   - [Recording](#recording)
   - [System](#system)
7. [Error Codes](#error-codes)

---

## Overview

The Tapo camera API is a JSON-based HTTP API running on the camera's local network. All requests are made via HTTPS POST to the camera's IP address.

### Key Concepts

| Concept | Description |
|---------|-------------|
| `stok` | Session token obtained after authentication |
| `cnonce` | Client-generated random nonce (8 hex characters) |
| `nonce` | Server-generated nonce returned during auth |
| `lsk` | AES encryption key (16 bytes) |
| `ivb` | AES initialization vector (16 bytes) |
| `seq` | Sequence number for encrypted requests |

---

## Connection Setup

### Base URL
```
https://{camera_ip}:443
```

### Default Ports
| Port | Purpose |
|------|---------|
| 443 | Control API (HTTPS) |
| 8800 | Video streaming |

### Required HTTP Headers
```http
Host: {camera_ip}:443
Referer: https://{camera_ip}
Accept: application/json
Accept-Encoding: gzip, deflate
User-Agent: Tapo CameraClient Android
Connection: close
requestByApp: true
Content-Type: application/json; charset=UTF-8
```

---

## Authentication

Tapo cameras support two authentication modes: **Secure (encrypt_type 3)** and **Legacy (insecure)**.

### Step 1: Detect Connection Type

```http
POST https://{camera_ip}
Content-Type: application/json

{
  "method": "login",
  "params": {
    "encrypt_type": "3",
    "username": "{username}"
  }
}
```

**Response indicating secure connection:**
```json
{
  "error_code": -40413,
  "result": {
    "data": {
      "encrypt_type": ["3"]
    }
  }
}
```

---

### Step 2a: Secure Authentication (encrypt_type 3)

#### Phase 1: Request Server Nonce

```json
{
  "method": "login",
  "params": {
    "cnonce": "{RANDOM_8_HEX_UPPERCASE}",
    "encrypt_type": "3",
    "username": "{username}"
  }
}
```

**Response:**
```json
{
  "error_code": 0,
  "result": {
    "data": {
      "nonce": "{server_nonce}",
      "device_confirm": "{validation_hash}"
    }
  }
}
```

#### Phase 2: Validate Device & Generate Keys

**Password Hashing:**
```
hashedMD5Password = MD5(password).toUpperCase()
hashedSHA256Password = SHA256(password).toUpperCase()
```

**Validate device_confirm:**
```
// Try SHA256 first, then MD5
expectedConfirm = SHA256(cnonce + hashedPassword + nonce).toUpperCase() + nonce + cnonce
if (device_confirm == expectedConfirm) {
    // Use this hash method
}
```

**Generate Encryption Keys:**
```
hashedKey = SHA256(cnonce + hashedPassword + nonce).toUpperCase()
lsk = SHA256("lsk" + cnonce + nonce + hashedKey)[0:16]  // First 16 bytes
ivb = SHA256("ivb" + cnonce + nonce + hashedKey)[0:16]  // First 16 bytes
```

#### Phase 3: Complete Login

```json
{
  "method": "login",
  "params": {
    "cnonce": "{cnonce}",
    "encrypt_type": "3",
    "digest_passwd": "{digestPasswd}{cnonce}{nonce}",
    "username": "{username}"
  }
}
```

Where `digestPasswd = SHA256(hashedPassword + cnonce + nonce).toUpperCase()`

**Success Response:**
```json
{
  "error_code": 0,
  "result": {
    "stok": "{session_token}",
    "start_seq": 1234,
    "user_group": "root"
  }
}
```

---

### Step 2b: Legacy Authentication (Insecure)

```json
{
  "method": "login",
  "params": {
    "hashed": true,
    "password": "{MD5(password).toUpperCase()}",
    "username": "{username}"
  }
}
```

**Response:**
```json
{
  "error_code": 0,
  "result": {
    "stok": "{session_token}"
  }
}
```

---

## Making Requests

### Request URL Format
```
POST https://{camera_ip}/stok={stok}/ds
```

### Simple Request (Legacy/Insecure)

```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {
        "method": "{method_name}",
        "params": { ... }
      }
    ]
  }
}
```

### Direct Request (for simple commands)

```json
{
  "method": "{get|set|do}",
  "{module}": { ... }
}
```

---

## Encryption

### Secure Request Format

For secure connections, encrypt the payload using AES-128-CBC:

```json
{
  "method": "securePassthrough",
  "params": {
    "request": "{base64(AES_CBC_encrypt(payload, lsk, ivb))}"
  }
}
```

### Additional Headers for Secure Requests

```http
Seq: {seq}
Tapo_tag: {tag}
```

**Tag Calculation:**
```
tag1 = SHA256(hashedPassword + cnonce).toUpperCase()
tag = SHA256(tag1 + JSON.stringify(request) + seq.toString()).toUpperCase()
```

> **Note:** Increment `seq` after each request.

### Decrypting Responses

```
encryptedResponse = base64_decode(response.result.response)
decryptedJSON = AES_CBC_decrypt(encryptedResponse, lsk, ivb)
// Remove PKCS7 padding
```

---

## API Commands Reference

### Method Types

| Method | Purpose | Example |
|--------|---------|---------|
| `get` | Read settings/status | Get privacy mode |
| `set` | Modify settings | Set LED enabled |
| `do` | Execute an action | Move motor |
| `multipleRequest` | Batch commands | Multiple operations |

---

### PTZ (Pan-Tilt-Zoom)

#### Move Motor to Coordinates
```json
{"method": "do", "motor": {"move": {"x_coord": "10", "y_coord": "5"}}}
```
- `x_coord`: Horizontal position (string)
- `y_coord`: Vertical position (string)

#### Move Motor by Direction
```json
{"method": "do", "motor": {"movestep": {"direction": "90"}}}
```
- `direction`: Angle 0-359 (string)
  - `0` = Right (clockwise)
  - `90` = Up (vertical)
  - `180` = Left (counter-clockwise)
  - `270` = Down (horizontal)

#### Calibrate Motor
```json
{"method": "do", "motor": {"manual_cali": ""}}
```

#### Get Motor Capability
```json
{"method": "get", "motor": {"name": ["capability"]}}
```

#### Start Cruise/Patrol
```json
{"method": "do", "motor": {"cruise": {"coord": "0"}}}
```

#### Stop Cruise
```json
{"method": "do", "motor": {"cruise_stop": {}}}
```

---

### Presets

#### Get All Presets
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getPresetConfig", "params": {"preset": {"name": ["preset"]}}}
    ]
  }
}
```

**Response:**
```json
{
  "result": {
    "preset": {
      "preset": {
        "id": ["1", "2", "3"],
        "name": ["Living Room", "Kitchen", "Front Door"]
      }
    }
  }
}
```

#### Go to Preset
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "motorMoveToPreset", "params": {"preset": {"goto_preset": {"id": "1"}}}}
    ]
  }
}
```

#### Save Current Position as Preset
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "addMotorPostion", "params": {"preset": {"set_preset": {"name": "MyPreset", "save_ptz": "1"}}}}
    ]
  }
}
```
> Note: The typo `addMotorPostion` is in the actual API.

#### Delete Preset
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "deletePreset", "params": {"preset": {"remove_preset": {"id": ["1"]}}}}
    ]
  }
}
```

---

### Device Info

#### Get Basic Info
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getDeviceInfo", "params": {"device_info": {"name": ["basic_info"]}}}
    ]
  }
}
```

#### Get Time
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getClockStatus", "params": {"system": {"name": "clock_status"}}}
    ]
  }
}
```

#### Get Module Specifications
```json
{"method": "get", "function": {"name": ["module_spec"]}}
```

---

### Privacy & Security

#### Get Privacy Mode (Lens Mask)
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getLensMaskConfig", "params": {"lens_mask": {"name": ["lens_mask_info"]}}}
    ]
  }
}
```

#### Set Privacy Mode
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setLensMaskConfig", "params": {"lens_mask": {"lens_mask_info": {"enabled": "on"}}}}
    ]
  }
}
```
- Values: `"on"` or `"off"`

#### Get Media Encryption
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getMediaEncrypt", "params": {"cet": {"name": ["media_encrypt"]}}}
    ]
  }
}
```

#### Set Media Encryption
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setMediaEncrypt", "params": {"cet": {"media_encrypt": {"enabled": "on"}}}}
    ]
  }
}
```

---

### Detection & Alarms

#### Get Motion Detection
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getDetectionConfig", "params": {"motion_detection": {"name": ["motion_det"]}}}
    ]
  }
}
```

#### Set Motion Detection
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setDetectionConfig", "params": {"motion_detection": {"motion_det": {"enabled": "on", "digital_sensitivity": "50"}}}}
    ]
  }
}
```

#### Get Person Detection
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getPersonDetectionConfig", "params": {"people_detection": {"name": ["detection"]}}}
    ]
  }
}
```

#### Set Person Detection
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setPersonDetectionConfig", "params": {"people_detection": {"detection": {"enabled": "on", "sensitivity": "50"}}}}
    ]
  }
}
```

#### Get Alarm Status
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getLastAlarmInfo", "params": {"msg_alarm": {"name": ["chn1_msg_alarm_info"]}}}
    ]
  }
}
```

#### Set Alarm
```json
{
  "method": "set",
  "msg_alarm": {
    "chn1_msg_alarm_info": {
      "enabled": "on",
      "alarm_type": "0",
      "light_type": "0",
      "alarm_mode": ["sound", "light"]
    }
  }
}
```

#### Start Manual Alarm
```json
{
  "method": "do",
  "msg_alarm": {
    "manual_msg_alarm": {
      "action": "start"
    }
  }
}
```

#### Stop Manual Alarm
```json
{
  "method": "do",
  "msg_alarm": {
    "manual_msg_alarm": {
      "action": "stop"
    }
  }
}
```

---

### Image & Video Settings

#### Get Common Image Settings
```json
{"method": "get", "image": {"name": "common"}}
```

#### Get/Set Image Flip
```json
// Get
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getLdc", "params": {"image": {"name": ["switch"]}}}
    ]
  }
}

// Set vertical flip
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setLdc", "params": {"image": {"switch": {"flip_type": "center"}}}}
    ]
  }
}
```
- `flip_type`: `"off"`, `"center"`, etc.

#### Get Day/Night Mode
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getLdc", "params": {"image": {"name": ["common", "switch"]}}}
    ]
  }
}
```

#### Set Day/Night Mode
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setLdc", "params": {"image": {"common": {"inf_type": "auto"}}}}
    ]
  }
}
```
- `inf_type`: `"auto"`, `"on"` (night), `"off"` (day)

---

### LED & Lights

#### Get LED Status
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getLedStatus", "params": {"led": {"name": ["config"]}}}
    ]
  }
}
```

#### Set LED Enabled
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "setLedStatus", "params": {"led": {"config": {"enabled": "on"}}}}
    ]
  }
}
```

---

### Audio

#### Get Audio Config
```json
{
  "method": "get",
  "audio_config": {
    "name": ["microphone", "speaker"]
  }
}
```

#### Set Speaker Volume
```json
{"method": "set", "audio_config": {"speaker": {"volume": "80"}}}
```

#### Set Microphone
```json
{"method": "set", "audio_config": {"microphone": {"volume": "50", "mute": "off"}}}
```

---

### Recording

#### Get Record Plan
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getRecordPlan", "params": {"record_plan": {"name": ["chn1_channel"]}}}
    ]
  }
}
```

#### Get SD Card Status
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "getSdCardStatus", "params": {"harddisk_manage": {"table": ["hd_info"]}}}
    ]
  }
}
```

#### Format SD Card
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "formatSdCard", "params": {"harddisk_manage": {"format_hd": "1"}}}
    ]
  }
}
```

---

### System

#### Reboot Camera
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "rebootDevice", "params": {"system": {"reboot": "null"}}}
    ]
  }
}
```

#### Check Firmware Update
```json
{
  "method": "multipleRequest",
  "params": {
    "requests": [
      {"method": "checkFirmwareVersionByCloud", "params": {"cloud_config": {"check_fw_version": "null"}}},
      {"method": "getCloudConfig", "params": {"cloud_config": {"name": ["upgrade_info"]}}}
    ]
  }
}
```

#### Start Firmware Upgrade
```json
{"method": "do", "cloud_config": {"fw_download": "null"}}
```

---

## Error Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `-40401` | Invalid/Expired Token (re-authenticate) |
| `-40404` | Temporary suspension (rate limited) |
| `-40411` | Invalid authentication data |
| `-40413` | Login required / Secure connection handshake |
| `-64303` | Cruise in progress (stop cruise first) |
| `-1` | General error |

---

## Response Format

### Success Response
```json
{
  "error_code": 0,
  "result": {
    "responses": [
      {
        "method": "getDeviceInfo",
        "result": { ... },
        "error_code": 0
      }
    ]
  }
}
```

### Error Response
```json
{
  "error_code": -40401,
  "result": {
    "err_msg": "Invalid stok value"
  }
}
```

---

## Implementation Tips

1. **Always use HTTPS** with SSL/TLS verification disabled (self-signed cert)
2. **Generate random cnonce** - 8 uppercase hex characters
3. **Handle session expiry** - Re-authenticate on error `-40401`
4. **Increment seq** after each encrypted request
5. **Use multipleRequest** for batching multiple commands
6. **Handle rate limiting** - Check `sec_left` in error responses
7. **Password hashing** - Try SHA256 first, fall back to MD5

---

## Credits

Based on reverse engineering from:
- [pytapo](https://github.com/JurajNyiri/pytapo)
- [HomeAssistant-Tapo-Control](https://github.com/JurajNyiri/HomeAssistant-Tapo-Control)
- Research by Dale Pavey (NCC Group), likaci, Tim Zhang, and others
