# Tapo Camera API

A Go Fiber REST API wrapper for TP-Link Tapo cameras. Control your Tapo cameras via HTTP endpoints.

## Features

- **Full PTZ Control** - Move, step, calibrate, cruise mode
- **Preset Management** - Create, list, goto, delete presets
- **Device Info** - Basic info, time, module specs
- **Privacy Controls** - Lens mask, media encryption
- **Detection Settings** - Motion and person detection
- **Alarm Control** - Configure and trigger alarms
- **Image Settings** - Flip, day/night mode
- **LED Control** - Status indicator toggle
- **Audio Settings** - Speaker and microphone volume
- **Recording** - Record plan, SD card management
- **System** - Reboot, firmware updates

## Installation

```bash
# Clone the repository
git clone https://github.com/budhilaw/gotapo-api.git
cd gotapo-api

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `3000` | Server port |
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `LOG_LEVEL` | `info` | Logging level |

## API Usage

All camera endpoints require authentication headers:

```
X-Tapo-Username: your_tapo_username
X-Tapo-Password: your_tapo_password
```

### Examples

#### Get Device Info
```bash
curl -X GET "http://localhost:3000/api/cameras/192.168.1.100/info" \
  -H "X-Tapo-Username: admin" \
  -H "X-Tapo-Password: yourpassword"
```

#### Move Camera (PTZ)
```bash
curl -X POST "http://localhost:3000/api/cameras/192.168.1.100/ptz/step" \
  -H "Content-Type: application/json" \
  -H "X-Tapo-Username: admin" \
  -H "X-Tapo-Password: yourpassword" \
  -d '{"direction": 90}'
```

#### Toggle Privacy Mode
```bash
curl -X PUT "http://localhost:3000/api/cameras/192.168.1.100/privacy" \
  -H "Content-Type: application/json" \
  -H "X-Tapo-Username: admin" \
  -H "X-Tapo-Password: yourpassword" \
  -d '{"enabled": true}'
```

#### Toggle LED
```bash
curl -X PUT "http://localhost:3000/api/cameras/192.168.1.100/led" \
  -H "Content-Type: application/json" \
  -H "X-Tapo-Username: admin" \
  -H "X-Tapo-Password: yourpassword" \
  -d '{"enabled": false}'
```

## API Endpoints

### PTZ
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/cameras/:ip/ptz/move` | Move to coordinates |
| POST | `/api/cameras/:ip/ptz/step` | Move by direction |
| POST | `/api/cameras/:ip/ptz/calibrate` | Calibrate motor |
| GET | `/api/cameras/:ip/ptz/capability` | Get motor capability |
| POST | `/api/cameras/:ip/ptz/cruise/start` | Start cruise |
| POST | `/api/cameras/:ip/ptz/cruise/stop` | Stop cruise |

### Presets
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/presets` | List presets |
| POST | `/api/cameras/:ip/presets` | Create preset |
| POST | `/api/cameras/:ip/presets/:id/goto` | Go to preset |
| DELETE | `/api/cameras/:ip/presets/:id` | Delete preset |

### Device
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/info` | Get device info |
| GET | `/api/cameras/:ip/time` | Get camera time |
| GET | `/api/cameras/:ip/specs` | Get module specs |

### Privacy & Security
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/privacy` | Get privacy mode |
| PUT | `/api/cameras/:ip/privacy` | Set privacy mode |
| GET | `/api/cameras/:ip/encryption` | Get encryption |
| PUT | `/api/cameras/:ip/encryption` | Set encryption |

### Detection
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/detection/motion` | Get motion detection |
| PUT | `/api/cameras/:ip/detection/motion` | Set motion detection |
| GET | `/api/cameras/:ip/detection/person` | Get person detection |
| PUT | `/api/cameras/:ip/detection/person` | Set person detection |

### Alarm
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/alarm` | Get alarm config |
| PUT | `/api/cameras/:ip/alarm` | Set alarm config |
| POST | `/api/cameras/:ip/alarm/trigger` | Start alarm |
| DELETE | `/api/cameras/:ip/alarm/trigger` | Stop alarm |

### Image
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/image` | Get image settings |
| PUT | `/api/cameras/:ip/image/flip` | Set image flip |
| PUT | `/api/cameras/:ip/image/nightmode` | Set night mode |

### LED
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/led` | Get LED status |
| PUT | `/api/cameras/:ip/led` | Set LED status |

### Audio
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/audio` | Get audio config |
| PUT | `/api/cameras/:ip/audio/speaker` | Set speaker volume |
| PUT | `/api/cameras/:ip/audio/microphone` | Set microphone |

### Recording & Storage
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/cameras/:ip/recording/plan` | Get record plan |
| GET | `/api/cameras/:ip/storage` | Get SD card status |
| POST | `/api/cameras/:ip/storage/format` | Format SD card |

### System
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/cameras/:ip/reboot` | Reboot camera |
| GET | `/api/cameras/:ip/firmware` | Check firmware |
| POST | `/api/cameras/:ip/firmware/upgrade` | Start upgrade |

## Running Tests

```bash
go test ./... -v
```

## Credits

This project is based on the protocol reverse-engineering from:
- [pytapo](https://github.com/JurajNyiri/pytapo) - Python library for Tapo cameras by Juraj Nyiri

## License

MIT
