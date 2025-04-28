package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"paint-drawer-pro/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	
	ui.Canvas.SetMinSize(fyne.NewSize(600, 500))

	
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

	
	colorLabel := widget.NewLabel("Color:")

	
	colorPreview := canvas.NewRectangle(ui.State.CurrentColor)
	colorPreview.SetMinSize(fyne.NewSize(30, 20))

	
	blackColor := color.RGBA{0, 0, 0, 255}
	redColor := color.RGBA{255, 0, 0, 255}
	greenColor := color.RGBA{0, 255, 0, 255}
	blueColor := color.RGBA{0, 0, 255, 255}

	
	blackBtn := widget.NewButton("", func() {
		ui.State.CurrentColor = blackColor
		colorPreview.FillColor = ui.State.CurrentColor
		colorPreview.Refresh()
		ui.StatusLabel.SetText("Black color selected")
	})
	blackBtn.Importance = widget.LowImportance
	
	blackRect := canvas.NewRectangle(blackColor)
	blackRect.SetMinSize(fyne.NewSize(20, 20))
	blackBtnContainer := container.NewHBox(blackRect, blackBtn)

	
	redBtn := widget.NewButton("", func() {
		ui.State.CurrentColor = redColor
		colorPreview.FillColor = ui.State.CurrentColor
		colorPreview.Refresh()
		ui.StatusLabel.SetText("Red color selected")
	})
	redBtn.Importance = widget.LowImportance
	
	redRect := canvas.NewRectangle(redColor)
	redRect.SetMinSize(fyne.NewSize(20, 20))
	redBtnContainer := container.NewHBox(redRect, redBtn)

	
	greenBtn := widget.NewButton("", func() {
		ui.State.CurrentColor = greenColor
		colorPreview.FillColor = ui.State.CurrentColor
		colorPreview.Refresh()
		ui.StatusLabel.SetText("Green color selected")
	})
	greenBtn.Importance = widget.LowImportance
	
	greenRect := canvas.NewRectangle(greenColor)
	greenRect.SetMinSize(fyne.NewSize(20, 20))
	greenBtnContainer := container.NewHBox(greenRect, greenBtn)

	
	blueBtn := widget.NewButton("", func() {
		ui.State.CurrentColor = blueColor
		colorPreview.FillColor = ui.State.CurrentColor
		colorPreview.Refresh()
		ui.StatusLabel.SetText("Blue color selected")
	})
	blueBtn.Importance = widget.LowImportance
	
	blueRect := canvas.NewRectangle(blueColor)
	blueRect.SetMinSize(fyne.NewSize(20, 20))
	blueBtnContainer := container.NewHBox(blueRect, blueBtn)

	
	customColorBtn := widget.NewButton("Custom...", func() {
		
		rSlider := widget.NewSlider(0, 255)
		gSlider := widget.NewSlider(0, 255)
		bSlider := widget.NewSlider(0, 255)

		
		rSlider.Value = float64(ui.State.CurrentColor.R)
		gSlider.Value = float64(ui.State.CurrentColor.G)
		bSlider.Value = float64(ui.State.CurrentColor.B)

		
		preview := canvas.NewRectangle(ui.State.CurrentColor)
		preview.SetMinSize(fyne.NewSize(100, 60))

		
		rLabel := widget.NewLabel(fmt.Sprintf("R: %d", ui.State.CurrentColor.R))
		gLabel := widget.NewLabel(fmt.Sprintf("G: %d", ui.State.CurrentColor.G))
		bLabel := widget.NewLabel(fmt.Sprintf("B: %d", ui.State.CurrentColor.B))

		
		updateColor := func() {
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

		
		rSlider.OnChanged = func(value float64) {
			updateColor()
		}
		gSlider.OnChanged = func(value float64) {
			updateColor()
		}
		bSlider.OnChanged = func(value float64) {
			updateColor()
		}

		
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

		
		dialog := dialog.NewCustom("Select Color", "Apply", content, window)
		dialog.SetOnClosed(func() {
			r := uint8(rSlider.Value)
			g := uint8(gSlider.Value)
			b := uint8(bSlider.Value)
			ui.State.CurrentColor = color.RGBA{r, g, b, 255}
			colorPreview.FillColor = ui.State.CurrentColor
			colorPreview.Refresh()
			ui.StatusLabel.SetText(fmt.Sprintf("Custom color selected (R:%d, G:%d, B:%d)", r, g, b))
		})
		dialog.Show()
	})

	
	colorButtons := container.NewGridWithColumns(2,
		blackBtnContainer, redBtnContainer, 
		greenBtnContainer, blueBtnContainer)

	colorContainer := container.NewVBox(
		container.NewHBox(colorLabel, colorPreview),
		colorButtons,
		customColorBtn,
	)

	
	ui.ToolsContainer = container.NewVBox(
		widget.NewLabel("Drawing Tools:"),
		lineBtn,
		circleBtn,
		pillBtn,
		polygonBtn,
		selectBtn,
		widget.NewSeparator(),
		clearBtn,
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
		colorLabel,
		colorContainer,
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


type buttonBackgroundLayout struct{}

func (d *buttonBackgroundLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}

	
	objects[0].Resize(size)
	objects[0].Move(fyne.NewPos(0, 0))

	
	objects[1].Resize(size)
	objects[1].Move(fyne.NewPos(0, 0))
}

func (d *buttonBackgroundLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(20, 20)
	}

	
	return objects[1].MinSize()
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
		
		
		for _, point := range controlPoints {
			drawSelectionIndicator(canvas, point.X, point.Y, 5, indicatorColor)
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