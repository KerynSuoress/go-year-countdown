# Year of Synchronicity - Countdown Timer

A Go-based desktop countdown timer that displays time remaining until the end of the year.

## Features
- Real-time countdown (months, weeks, days, hours, minutes)
- Custom rounded cards with navy blue theme

## Requirements
- Go 1.16+ (for embed functionality)
- Windows (for GUI executable)

## Quick Start

### Run Directly
```bash
go run main.go
```

### OR Create Executable
```bash
# Standard build
go build -o year-countdown.exe .

# Windows GUI build (no console window)
go build -ldflags -H=windowsgui -o year-countdown.exe .
```

### Windows Startup (Optional)
1. Press `Win + R`, type `shell:startup`, press Enter
2. Copy `year-countdown.exe` to this folder
3. Timer will start with Windows

To disable: Delete the file from startup folder

## Building with Icon (Advanced)
```bash
# Generate Windows resources (swap for your information)
goversioninfo -64 versioninfo.json

# Build with icon and metadata
go build -ldflags -H=windowsgui -o year-countdown.exe .
```

## Dependencies
- [Fyne v2](https://fyne.io/) - Cross-platform GUI toolkit