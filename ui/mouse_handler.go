package ui

import (
	"image/color"
	"math"
	"paint-drawer-pro/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)


type MouseHandler struct {
	widget.BaseWidget
	UI           *MainUI
	StartPoint   models.Point
	CurrentPoint models.Point
	IsDrawing    bool
	LastPoint    models.Point
	PolyPoints   []models.Point
}


func NewMouseHandler(ui *MainUI) *MouseHandler {
	handler := &MouseHandler{
		UI: ui,
	}
	handler.ExtendBaseWidget(handler)
	return handler
}


func (h *MouseHandler) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(color.RGBA{255, 255, 255, 0}) 
	return widget.NewSimpleRenderer(background)
}


func (h *MouseHandler) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}




func (h *MouseHandler) MouseDown(ev *desktop.MouseEvent) {
	adjustedPoint := h.adjustMousePosition(ev.PointEvent)
	h.StartPoint = adjustedPoint
	h.CurrentPoint = adjustedPoint
	
	if h.UI.State.CurrentAction == "polygon" && ev.Button == desktop.MouseButtonPrimary {
		
		if len(h.PolyPoints) == 0 {
			h.PolyPoints = append(h.PolyPoints, h.StartPoint)
			h.UI.StatusLabel.SetText("Creating polygon... Click to add points, press Enter to finish")
		} else {
			h.PolyPoints = append(h.PolyPoints, h.StartPoint)
			if len(h.PolyPoints) >= 3 {
				poly := models.NewPolygon(h.PolyPoints, h.UI.State.CurrentColor, 1)
				h.UI.State.CurrentShape = poly
				h.UI.Canvas.Refresh()
			}
			h.UI.StatusLabel.SetText("Added point to polygon. Click for more points, press Enter to finish")
		}
		return
	}
	
	
	if h.UI.State.CurrentAction == "pill" && ev.Button == desktop.MouseButtonPrimary {
		
		if pill, isPill := h.UI.State.CurrentShape.(*models.Pill); isPill {
			if pill.Step == 1 {
				
				dx := adjustedPoint.X - pill.Start.X
				dy := adjustedPoint.Y - pill.Start.Y
				pill.Radius = int(math.Sqrt(float64(dx*dx + dy*dy)))
				pill.Step = 2
				h.UI.StatusLabel.SetText("Pill radius set. Click to place the second end.")
				h.UI.Canvas.Refresh()
				return
			} else if pill.Step == 2 {
				
				pill.End = adjustedPoint
				pill.Step = 3
				
				
				h.UI.State.Shapes = append(h.UI.State.Shapes, pill)
				h.UI.State.CurrentShape = nil
				h.UI.StatusLabel.SetText("Pill added")
				h.UI.Canvas.Refresh()
				return
			}
		} else {
			
			pill := models.NewPill(adjustedPoint, 5, h.UI.State.CurrentColor)
			h.UI.State.CurrentShape = pill
			h.UI.StatusLabel.SetText("Pill started. Click to set radius.")
			h.UI.Canvas.Refresh()
			return
		}
	}
	
	h.IsDrawing = true
	
	switch h.UI.State.CurrentAction {
	case "line":
		thickness := 1
		if h.UI.State.PenType == "brush" {
			thickness = h.UI.State.BrushThickness 
		}
		
		line := models.NewLine(
			h.StartPoint,
			h.StartPoint, 
			h.UI.State.CurrentColor, 
			thickness,
			h.UI.State.PenType,
		)
		h.UI.State.CurrentShape = line
		h.UI.StatusLabel.SetText("Drawing line... Release to complete")
		
	case "circle":
		circle := models.NewCircle(
			h.StartPoint,
			1, 
			h.UI.State.CurrentColor, 
		)
		h.UI.State.CurrentShape = circle
		h.UI.StatusLabel.SetText("Drawing circle... Release to complete")
	}
}



func (h *MouseHandler) MouseUp(ev *desktop.MouseEvent) {
	if !h.IsDrawing || h.UI.State.CurrentAction == "polygon" {
		return
	}
	
	h.CurrentPoint = h.adjustMousePosition(ev.PointEvent)
	
	
	if h.UI.State.CurrentShape != nil {
		
		h.UI.State.Shapes = append(h.UI.State.Shapes, h.UI.State.CurrentShape)
		h.UI.State.CurrentShape = nil
		h.IsDrawing = false
		h.UI.Canvas.Refresh()
		
		switch h.UI.State.CurrentAction {
		case "line":
			h.UI.StatusLabel.SetText("Line added")
		case "circle":
			h.UI.StatusLabel.SetText("Circle added")
		}
	}
}




func (h *MouseHandler) MouseMoved(ev *desktop.MouseEvent) {
	if !h.IsDrawing || h.UI.State.CurrentShape == nil {
		return
	}
	
	h.CurrentPoint = h.adjustMousePosition(ev.PointEvent)
	
	switch shape := h.UI.State.CurrentShape.(type) {
	case *models.Line:
		shape.End = h.CurrentPoint
		h.UI.Canvas.Refresh()
		
	case *models.Circle:
		
		dx := h.CurrentPoint.X - shape.Center.X
		dy := h.CurrentPoint.Y - shape.Center.Y
		shape.Radius = int(math.Sqrt(float64(dx*dx + dy*dy)))
		h.UI.Canvas.Refresh()
		
	case *models.Pill:
		
		if shape.Step == 1 {
			
			return
		} else if shape.Step == 2 {
			
			shape.End = h.CurrentPoint
			h.UI.Canvas.Refresh()
		}
	}
}


func (h *MouseHandler) KeyDown(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn && h.UI.State.CurrentAction == "polygon" && len(h.PolyPoints) >= 3 {
		
		poly := models.NewPolygon(h.PolyPoints, h.UI.State.CurrentColor, 1)
		h.UI.State.Shapes = append(h.UI.State.Shapes, poly)
		h.PolyPoints = nil
		h.UI.State.CurrentShape = nil
		h.UI.Canvas.Refresh()
		h.UI.StatusLabel.SetText("Polygon added")
	} else if ev.Name == fyne.KeyEscape {
		
		h.UI.State.CurrentShape = nil
		h.IsDrawing = false
		h.PolyPoints = nil
		h.UI.Canvas.Refresh()
		h.UI.StatusLabel.SetText("Drawing canceled")
	}
}


func (h *MouseHandler) Dragged(ev *fyne.DragEvent) {
	h.MouseMoved(&desktop.MouseEvent{
		PointEvent: fyne.PointEvent{
			AbsolutePosition: ev.AbsolutePosition,
			Position:         ev.Position,
		},
	})
}


func (h *MouseHandler) DragEnd() {
	
}


func (h *MouseHandler) Tapped(ev *fyne.PointEvent) {
	h.MouseDown(&desktop.MouseEvent{
		PointEvent: *ev,
		Button:     desktop.MouseButtonPrimary,
	})
	
	
	if h.UI.State.CurrentAction == "polygon" {
		return
	}
	
	
	h.MouseUp(&desktop.MouseEvent{
		PointEvent: *ev,
		Button:     desktop.MouseButtonPrimary,
	})
}


func (h *MouseHandler) TappedSecondary(ev *fyne.PointEvent) {
	
	h.UI.State.CurrentShape = nil
	h.IsDrawing = false
	h.UI.Canvas.Refresh()
	h.UI.StatusLabel.SetText("Drawing canceled")
}



func (h *MouseHandler) adjustMousePosition(ev fyne.PointEvent) models.Point {
	
	
	x := int(ev.Position.X)
	y := int(ev.Position.Y)
	
	
	canvasSize := h.UI.Canvas.Size()
	maxX := int(canvasSize.Width) - 1
	maxY := int(canvasSize.Height) - 1
	
	
	if x < 0 {
		x = 0
	} else if x > maxX {
		x = maxX
	}
	
	if y < 0 {
		y = 0
	} else if y > maxY {
		y = maxY
	}
	
	return models.Point{X: x, Y: y}
}