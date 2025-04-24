# Paint Drawer Pro

A drawing application built with Go and Fyne that implements various computer graphics algorithms for drawing shapes like lines, circles, and polygons.

## Features

- Draw lines using different pen types:
  - Regular pen (1px thickness)
  - Brush with adjustable thickness (1-10px)
- Create circles with precise radius control
- Draw polygons by placing points and connecting them
- Anti-aliasing toggle for smoother drawings
- Shape selection and manipulation
- Clean and intuitive user interface

## Tech Stack

- **Programming Language**: Go 1.24
- **UI Framework**: [Fyne](https://fyne.io/) v2.6.0 - A cross-platform GUI toolkit
- **Graphics Algorithms**:
  - Midpoint Line Algorithm
  - Midpoint Circle Algorithm v2
  - Xiaolin Wu's Line Algorithm (anti-aliasing)
  - Xiaolin Wu's Circle Algorithm (anti-aliasing)
  - Circular brush technique for thick lines

## Project Structure

```
paint-drawer-pro/
├── algorithms/          # Core drawing algorithms
│   └── drawing.go       # Implementation of line and circle algorithms
├── models/              # Data models and shape definitions
│   ├── circle.go        # Circle shape implementation
│   ├── drawing_utils.go # Utilities for drawing shapes
│   ├── line.go          # Line shape implementation with pen types
│   ├── polygon.go       # Polygon shape implementation
│   └── shapes.go        # Shape interfaces and drawing state
├── ui/                  # User interface components
│   ├── main_ui.go       # Main UI layout and controls
│   └── mouse_handler.go # Mouse interaction handling
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── main.go              # Application entry point
```

## Prerequisites

- Go 1.24 or later
- Fyne dependencies:
  - For Linux: `xorg-dev` or equivalent
  - For Windows: GCC and a C compiler
  - For macOS: Xcode command line tools

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd paint-drawer-pro
```

2. Install dependencies:

```bash
go mod download
```

## Running the Application

Run the application using Go:

```bash
go run main.go
```

Or build and run the executable:

```bash
go build -o paint-drawer-pro
./paint-drawer-pro  # On Linux/macOS
paint-drawer-pro.exe  # On Windows
```

## Usage

- Select a drawing tool (Line, Circle, Polygon) from the toolbar
- Toggle between Regular Pen and Brush using the radio buttons
- Adjust brush thickness with the slider (when Brush is selected)
- Click and drag to draw lines and circles
- For polygons:
  - Click to place points
  - Press Enter to complete the polygon
  - Press Escape to cancel
- Toggle anti-aliasing for smoother drawings

## Building for Distribution

Build for the current platform:

```bash
go build
```

Cross-compile for other platforms (examples):

```bash
# For Windows from a non-Windows system
GOOS=windows GOARCH=amd64 go build -o paint-drawer-pro.exe

# For macOS from a non-macOS system
GOOS=darwin GOARCH=amd64 go build -o paint-drawer-pro-mac

# For Linux from a non-Linux system
GOOS=linux GOARCH=amd64 go build -o paint-drawer-pro-linux
```

## License

MIT

## Contributors

Me