package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"paint-drawer-pro/algorithms"
	"paint-drawer-pro/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type MainUI struct {
	Window          fyne.Window
	Container       *fyne.Container
	Canvas          *canvas.Raster
	ToolsContainer  *fyne.Container
	StatusLabel     *widget.Label
	CurrentToolText *widget.Label
	PillLengthSlider *widget.Slider
	PillLengthLabel  *widget.Label
	PillLengthContainer *fyne.Container
	State           models.DrawingState
}

func NewMainUI(window fyne.Window) *MainUI {
	ui := &MainUI{
		Window: window,
		State: models.DrawingState{
			Shapes:         []models.Shape{},
			CurrentAction:  "line",
			AntiAliasing:   true,
			PenType:        "brush",
			BrushThickness: 3,
			CurrentColor:   color.RGBA{0, 0, 0, 255}, 
			FillEnabled:    false,
			FillColor:      color.RGBA{255, 255, 255, 255},
			UseImageFill:   false,
		},
	}

	
	ui.PillLengthLabel = widget.NewLabel("Pill Length:")
	pillLengthValue := widget.NewLabel("100")
	ui.PillLengthSlider = widget.NewSlider(50, 600)
	ui.PillLengthSlider.SetValue(100)
	ui.PillLengthSlider.Step = 1
	ui.PillLengthSlider.OnChanged = func(value float64) {
		length := int(value)
		pillLengthValue.SetText(fmt.Sprintf("%d", length))
		ui.StatusLabel.SetText(fmt.Sprintf("Pill length set to %d", length))
		
		
		ui.updatePillLength(length)
	}
	ui.PillLengthContainer = container.NewBorder(
		nil, nil, ui.PillLengthLabel, pillLengthValue, ui.PillLengthSlider,
	)
	
	ui.PillLengthContainer.Hide()

	ui.Canvas = canvas.NewRaster(func(w, h int) image.Image {
		return ui.renderCanvas(w, h)
	})

	// Set both min size and a fixed size for the canvas to ensure proper positioning
	ui.Canvas.SetMinSize(fyne.NewSize(600, 500))
	// Make the canvas resize with the window
	ui.Canvas.Resize(fyne.NewSize(600, 500))

	
	ui.StatusLabel = widget.NewLabel("Ready")
	ui.CurrentToolText = widget.NewLabel("Current tool: Line")

	
	lineBtn := widget.NewButton("Line", func() {
		ui.State.CurrentAction = "line"
		ui.CurrentToolText.SetText("Current tool: Line")
		ui.StatusLabel.SetText("Line tool selected")
		ui.PillLengthContainer.Hide() 
	})

	circleBtn := widget.NewButton("Circle", func() {
		ui.State.CurrentAction = "circle"
		ui.CurrentToolText.SetText("Current tool: Circle")
		ui.StatusLabel.SetText("Circle tool selected")
		ui.PillLengthContainer.Hide() 
	})

	polygonBtn := widget.NewButton("Polygon", func() {
		ui.State.CurrentAction = "polygon"
		ui.CurrentToolText.SetText("Current tool: Polygon")
		ui.StatusLabel.SetText("Polygon tool selected")
		ui.PillLengthContainer.Hide() 
	})

	rectangleBtn := widget.NewButton("Rectangle", func() {
		ui.State.CurrentAction = "rectangle"
		ui.CurrentToolText.SetText("Current tool: Rectangle")
		ui.StatusLabel.SetText("Rectangle tool selected")
		ui.PillLengthContainer.Hide() 
	})

	pillBtn := widget.NewButton("Pill", func() {
		ui.State.CurrentAction = "pill"
		ui.CurrentToolText.SetText("Current tool: Pill")
		ui.StatusLabel.SetText("Pill tool selected: Click to place first end, then set radius, then place second end")
		
		
		ui.PillLengthContainer.Show()
		
		ui.PillLengthSlider.SetValue(100)
	})

	selectBtn := widget.NewButton("Select", func() {
		ui.State.CurrentAction = "select"
		ui.CurrentToolText.SetText("Current tool: Select")
		ui.StatusLabel.SetText("Select tool active")
		
		
		ui.PillLengthContainer.Hide()
	})

	clearBtn := widget.NewButton("Clear All", func() {
		ui.State.Shapes = []models.Shape{}
		ui.Canvas.Refresh()
		ui.StatusLabel.SetText("Canvas cleared")
	})
	
	// File save button
	saveBtn := widget.NewButton("Save", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			if writer == nil {
				return // User cancelled
			}
			
			// Close the writer as we'll use our own file handler
			writer.Close()
			
			// Get the file path
			filePath := writer.URI().Path()
			
			// Save shapes to the file
			err = ui.SaveShapesToFile(filePath)
			if err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			
			ui.StatusLabel.SetText("Drawing saved to file")
		}, ui.Window)
		
		fd.SetFileName("drawing.json")
		fd.Show()
	})
	
	// File load button
	loadBtn := widget.NewButton("Load", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			if reader == nil {
				return // User cancelled
			}
			
			// Close the reader as we'll use our own file handler
			reader.Close()
			
			// Get the file path
			filePath := reader.URI().Path()
			
			// Load shapes from the file
			err = ui.LoadShapesFromFile(filePath)
			if err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			
			ui.StatusLabel.SetText("Drawing loaded from file")
		}, ui.Window)
		
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fd.Show()
	})
	
	// Clipping button
	clipBtn := widget.NewButton("Clip Polygon", func() {
		// Enable only when we have a selected shape that's a polygon
		if ui.State.SelectedShape == nil {
			dialog.ShowInformation("Clipping", "Please select a polygon to clip first.", ui.Window)
			return
		}
		
		selectedPoly, isPolygon := ui.State.SelectedShape.(*models.Polygon)
		if !isPolygon {
			dialog.ShowInformation("Clipping", "Only polygons can be clipped. Please select a polygon.", ui.Window)
			return
		}
		
		// Convert to algorithm points for simplification
		vertices := selectedPoly.GetVertices()
		algVertices := make([]algorithms.Point, len(vertices))
		for i, v := range vertices {
			algVertices[i] = algorithms.Point{X: v.X, Y: v.Y}
		}
		
		// Simplify the polygon to remove possible duplicate points
		simplified := algorithms.SimplifyPolygon(algVertices, 2.0)
		
		// Display simplified vertex count
		originalCount := len(vertices)
		simplifiedCount := len(simplified)
		
		// Check if the simplified polygon is convex
		if !algorithms.IsPolygonConvex(simplified) {
			message := fmt.Sprintf("Only convex polygons can be used for clipping. Selected polygon has %d vertices (simplified from %d).", 
				simplifiedCount, originalCount)
			dialog.ShowInformation("Clipping", message, ui.Window)
			return
		}
		
		ui.State.CurrentAction = "clipping"
		ui.CurrentToolText.SetText("Current tool: Clipping")
		ui.StatusLabel.SetText("Clipping mode active. Select a polygon to clip against " + 
			"the current selection.")
	})

	aaCheck := widget.NewCheck("Anti-aliasing", func(checked bool) {
		ui.State.AntiAliasing = checked
		ui.Canvas.Refresh()
		if checked {
			ui.StatusLabel.SetText("Anti-aliasing enabled")
		} else {
			ui.StatusLabel.SetText("Anti-aliasing disabled")
		}
	})
	aaCheck.SetChecked(ui.State.AntiAliasing)

	penTypeLabel := widget.NewLabel("Pen Type:")
	regularPenRadio := widget.NewRadioGroup([]string{"Regular Pen", "Brush"}, func(selected string) {
		if selected == "Regular Pen" {
			ui.State.PenType = "regular"
			ui.StatusLabel.SetText("Regular Pen selected")
		} else {
			ui.State.PenType = "brush"
			ui.StatusLabel.SetText("Brush selected")
		}
		ui.Canvas.Refresh()
	})
	regularPenRadio.SetSelected("Brush")

	thicknessLabel := widget.NewLabel("Brush Thickness:")
	thicknessValue := widget.NewLabel("3")
	thicknessSlider := widget.NewSlider(1, 10)
	thicknessSlider.SetValue(3)
	thicknessSlider.Step = 1
	thicknessSlider.OnChanged = func(value float64) {
		thickness := int(value)
		ui.State.BrushThickness = thickness
		thicknessValue.SetText(fmt.Sprintf("%d", thickness))
		ui.StatusLabel.SetText(fmt.Sprintf("Brush thickness set to %d", thickness))
	}
	thicknessContainer := container.NewBorder(
		nil, nil, thicknessLabel, thicknessValue, thicknessSlider,
	)

	// Fill controls
	fillCheck := widget.NewCheck("Fill Shapes", func(checked bool) {
		ui.State.FillEnabled = checked
		ui.StatusLabel.SetText(fmt.Sprintf("Fill %s", map[bool]string{true: "enabled", false: "disabled"}[checked]))
	})
	
	fillColorBtn := widget.NewButton("Fill Color", func() {
		// Create RGB sliders for fill color
		rSlider := widget.NewSlider(0, 255)
		gSlider := widget.NewSlider(0, 255)
		bSlider := widget.NewSlider(0, 255)

		// Set initial values to current fill color or default
		fillColor := ui.State.FillColor
		if fillColor == nil {
			fillColor = color.RGBA{255, 255, 255, 255}
		}
		
		r, g, b, _ := fillColor.RGBA()
		rSlider.Value = float64(uint8(r))
		gSlider.Value = float64(uint8(g))
		bSlider.Value = float64(uint8(b))

		// Color preview
		preview := canvas.NewRectangle(fillColor)
		preview.SetMinSize(fyne.NewSize(100, 60))

		// Labels
		rLabel := widget.NewLabel(fmt.Sprintf("R: %d", uint8(r)))
		gLabel := widget.NewLabel(fmt.Sprintf("G: %d", uint8(g)))
		bLabel := widget.NewLabel(fmt.Sprintf("B: %d", uint8(b)))

		// Update function
		updateFillColor := func() {
			r := uint8(rSlider.Value)
			g := uint8(gSlider.Value)
			b := uint8(bSlider.Value)
			newColor := color.RGBA{r, g, b, 255}

			preview.FillColor = newColor
			preview.Refresh()

			rLabel.SetText(fmt.Sprintf("R: %d", r))
			gLabel.SetText(fmt.Sprintf("G: %d", g))
			bLabel.SetText(fmt.Sprintf("B: %d", b))
			}

		// Slider change handlers
		rSlider.OnChanged = func(value float64) {
			updateFillColor()
		}
		gSlider.OnChanged = func(value float64) {
			updateFillColor()
		}
		bSlider.OnChanged = func(value float64) {
			updateFillColor()
		}

		// Dialog content
		content := container.NewVBox(
			preview,
			widget.NewSeparator(),
			container.NewHBox(widget.NewLabel("Red:"), rLabel),
			rSlider,
			container.NewHBox(widget.NewLabel("Green:"), gLabel),
			gSlider,
			container.NewHBox(widget.NewLabel("Blue:"), bLabel),
			bSlider,
			)

		// Create and show dialog
		customDialog := dialog.NewCustom("Choose Fill Color", "Apply", content, ui.Window)
		customDialog.SetOnClosed(func() {
			newColor := color.RGBA{
				R: uint8(rSlider.Value),
				G: uint8(gSlider.Value),
				B: uint8(bSlider.Value),
				A: 255,
			}
			ui.State.FillColor = newColor
			ui.StatusLabel.SetText("Fill color updated")
			
			// If a shape is selected, apply the fill color
			if ui.State.SelectedShape != nil {
				switch s := ui.State.SelectedShape.(type) {
				case *models.Polygon:
					s.SetFillColor(newColor)
				case *models.Rectangle:
					s.SetFillColor(newColor)
				}
				ui.Canvas.Refresh()
			}
		})
		customDialog.Show()
	})
	
	loadImageBtn := widget.NewButton("Load Fill Image", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			if reader == nil {
				return
			}
			
			imgData, _, err := image.Decode(reader)
			if err != nil {
				dialog.ShowError(err, ui.Window)
				return
			}
			
			reader.Close()
			
			bounds := imgData.Bounds()
			width, height := bounds.Max.X, bounds.Max.Y
			
			fillImage := make([][]color.Color, height)
			for y := 0; y < height; y++ {
				fillImage[y] = make([]color.Color, width)
				for x := 0; x < width; x++ {
					fillImage[y][x] = imgData.At(x, y)
				}
			}
			
			ui.State.FillImage = fillImage
			ui.State.UseImageFill = true
			ui.StatusLabel.SetText("Fill image loaded")
			
			// Apply to selected shape if any
			if ui.State.SelectedShape != nil {
				switch s := ui.State.SelectedShape.(type) {
				case *models.Polygon:
					s.SetFillImage(fillImage)
				case *models.Rectangle:
					s.SetFillImage(fillImage)
				}
				ui.Canvas.Refresh()
			}
		}, ui.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
	
	fillContainer := container.NewVBox(
		fillCheck,
		container.NewHBox(fillColorBtn, loadImageBtn),
	)
	
	ui.ToolsContainer = container.NewVBox(
		widget.NewLabel("Drawing Tools:"),
		lineBtn,
		circleBtn,
		pillBtn,
		polygonBtn,
		rectangleBtn,
		selectBtn,
		widget.NewSeparator(),
		clearBtn,
		saveBtn, // Add save button to the tools container
		loadBtn, // Add load button to the tools container
		clipBtn,
		widget.NewSeparator(),
		aaCheck,
		widget.NewSeparator(),
		penTypeLabel,
		regularPenRadio,
		widget.NewSeparator(),
		thicknessContainer,
		widget.NewSeparator(),
		ui.PillLengthContainer, 
		widget.NewSeparator(),
		fillContainer,
		widget.NewSeparator(),
		widget.NewSeparator(),
		ui.CurrentToolText,
	)

	
	statusBar := container.NewHBox(
		ui.StatusLabel,
	)

	ui.Container = container.NewBorder(
		nil, statusBar, ui.ToolsContainer, nil,
		ui.Canvas,
	)

	return ui
}

func (ui *MainUI) renderCanvas(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.White)
		}
	}

	
	for _, shape := range ui.State.Shapes {
		canvas := make([][]color.Color, h)
		for j := range canvas {
			canvas[j] = make([]color.Color, w)
			for i := 0; i < w; i++ {
				canvas[j][i] = img.At(i, j)
			}
		}

		shape.Draw(canvas, ui.State.AntiAliasing)

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if canvas[y][x] != nil {
					img.Set(x, y, canvas[y][x])
				}
			}
		}
	}

	
	if ui.State.CurrentShape != nil {
		canvas := make([][]color.Color, h)
		for j := range canvas {
			canvas[j] = make([]color.Color, w)
			for i := 0; i < w; i++ {
				canvas[j][i] = img.At(i, j)
			}
		}

		ui.State.CurrentShape.Draw(canvas, ui.State.AntiAliasing)

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if canvas[y][x] != nil {
					img.Set(x, y, canvas[y][x])
				}
			}
		}
	}

	
	if ui.State.CurrentAction == "select" && ui.State.SelectedShape != nil {
		
		controlPoints := ui.State.SelectedShape.GetControlPoints()
		
		
		indicatorColor := color.RGBA{0, 119, 255, 255} 
		
		
		canvas := make([][]color.Color, h)
		for j := range canvas {
			canvas[j] = make([]color.Color, w)
		}
		
		// Special handling for rectangles to show resize handles
		if rect, isRect := ui.State.SelectedShape.(*models.Rectangle); isRect {
			drawRectangleSelectionHandles(canvas, rect, indicatorColor)
		} else {
			// Regular selection indicators for other shapes
			for _, point := range controlPoints {
				drawSelectionIndicator(canvas, point.X, point.Y, 5, indicatorColor)
			}
		}
		
		
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if canvas[y][x] != nil {
					img.Set(x, y, canvas[y][x])
				}
			}
		}
	}

	return img
}


func drawSelectionIndicator(canvas [][]color.Color, x, y, size int, c color.Color) {
	halfSize := size / 2
	
	
	if len(canvas) == 0 || len(canvas[0]) == 0 {
		return
	}
	
	
	for dy := -halfSize; dy <= halfSize; dy++ {
		for dx := -halfSize; dx <= halfSize; dx++ {
			
			px := x + dx
			py := y + dy
			
			
			if py >= 0 && py < len(canvas) && px >= 0 && px < len(canvas[0]) {
				
				if dx == 0 && dy == 0 {
					canvas[py][px] = c
				} else if dx == -halfSize || dx == halfSize || dy == -halfSize || dy == halfSize {
					canvas[py][px] = c
				}
			}
		}
	}
}

// drawRectangleSelectionHandles draws handles at the corners of a selected rectangle
func drawRectangleSelectionHandles(canvas [][]color.Color, rect *models.Rectangle, c color.Color) {
	// Draw corner handles with larger size for easier grabbing
	points := rect.GetControlPoints()
	for _, point := range points {
		drawSelectionIndicator(canvas, point.X, point.Y, 8, c)
	}
}


func (ui *MainUI) updatePillLength(length int) {
	
	if ui.State.CurrentAction == "select" && ui.State.SelectedShape != nil {
		if pill, ok := ui.State.SelectedShape.(*models.Pill); ok {
			
			dx := pill.End.X - pill.Start.X
			dy := pill.End.Y - pill.Start.Y
			currentLength := math.Sqrt(float64(dx*dx + dy*dy))
			
			
			if currentLength <= 0 {
				return
			}
			
			
			dirX := float64(dx) / currentLength
			dirY := float64(dy) / currentLength
			
			
			pill.End.X = pill.Start.X + int(dirX * float64(length))
			pill.End.Y = pill.Start.Y + int(dirY * float64(length))
			
			ui.Canvas.Refresh()
		}
	} else if ui.State.CurrentShape != nil {
		
		if pill, ok := ui.State.CurrentShape.(*models.Pill); ok && pill.Step >= 2 {
			
			dx := pill.End.X - pill.Start.X
			dy := pill.End.Y - pill.Start.Y
			currentLength := math.Sqrt(float64(dx*dx + dy*dy))
			
			
			if currentLength <= 0 {
				return
			}
			
			
			dirX := float64(dx) / currentLength
			dirY := float64(dy) / currentLength
			
			
			pill.End.X = pill.Start.X + int(dirX * float64(length))
			pill.End.Y = pill.Start.Y + int(dirY * float64(length))
			
			ui.Canvas.Refresh()
		}
	}
}