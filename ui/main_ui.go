package ui

import (
	"fmt"
	"image"
	"image/color"
	"paint-drawer-pro/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)


type MainUI struct {
	Window          fyne.Window
	Container       *fyne.Container
	Canvas          *canvas.Raster
	ToolsContainer  *fyne.Container
	StatusLabel     *widget.Label
	CurrentToolText *widget.Label
	State           models.DrawingState
}


func NewMainUI(window fyne.Window) *MainUI {
	ui := &MainUI{
		Window: window,
		State: models.DrawingState{
			Shapes:         []models.Shape{},
			CurrentAction:  "draw",
			AntiAliasing:   true,
			PenType:        "brush", 
			BrushThickness: 3,       
		},
	}

	
	ui.Canvas = canvas.NewRaster(func(w, h int) image.Image {
		return ui.renderCanvas(w, h)
	})
	
	
	ui.StatusLabel = widget.NewLabel("Ready")
	ui.CurrentToolText = widget.NewLabel("Current tool: Line")

	
	lineBtn := widget.NewButton("Line", func() {
		ui.State.CurrentAction = "line"
		ui.CurrentToolText.SetText("Current tool: Line")
		ui.StatusLabel.SetText("Line tool selected")
	})

	circleBtn := widget.NewButton("Circle", func() {
		ui.State.CurrentAction = "circle"
		ui.CurrentToolText.SetText("Current tool: Circle")
		ui.StatusLabel.SetText("Circle tool selected")
	})

	polygonBtn := widget.NewButton("Polygon", func() {
		ui.State.CurrentAction = "polygon"
		ui.CurrentToolText.SetText("Current tool: Polygon")
		ui.StatusLabel.SetText("Polygon tool selected")
	})

	selectBtn := widget.NewButton("Select", func() {
		ui.State.CurrentAction = "select"
		ui.CurrentToolText.SetText("Current tool: Select")
		ui.StatusLabel.SetText("Select tool active")
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
	
	
	ui.ToolsContainer = container.NewVBox(
		widget.NewLabel("Drawing Tools:"),
		lineBtn,
		circleBtn,
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
	
	return img
}