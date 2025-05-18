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
	IsMoving     bool      
	MoveStartX   int       
	MoveStartY   int       
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
	
	
	if h.UI.State.CurrentAction == "select" && ev.Button == desktop.MouseButtonPrimary {
		
		h.UI.State.SelectedShape = nil
		
		h.UI.PillLengthContainer.Hide()
		
		
		for i := len(h.UI.State.Shapes) - 1; i >= 0; i-- {
			shape := h.UI.State.Shapes[i]
			if shape.Contains(adjustedPoint) {
				h.UI.State.SelectedShape = shape
				h.IsMoving = true
				h.MoveStartX = adjustedPoint.X
				h.MoveStartY = adjustedPoint.Y
				
				
				if pill, isPill := shape.(*models.Pill); isPill {
					dx := pill.End.X - pill.Start.X
					dy := pill.End.Y - pill.Start.Y
					length := math.Sqrt(float64(dx*dx + dy*dy))
					h.UI.PillLengthSlider.SetValue(length)
					h.UI.PillLengthContainer.Show()
					h.UI.StatusLabel.SetText("Pill selected. Use slider to adjust length or drag to move.")
				} else {
					h.UI.StatusLabel.SetText("Shape selected. Drag to move.")
				}
				
				h.UI.Canvas.Refresh()
				return
			}
		}
		
		h.UI.StatusLabel.SetText("No shape selected.")
		return
	}
	
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
				
				
				h.UI.PillLengthContainer.Show()
				h.UI.PillLengthSlider.SetValue(float64(pill.Radius * 4)) 
				
				
				if dx != 0 || dy != 0 {
					length := float64(h.UI.PillLengthSlider.Value)
					distance := math.Sqrt(float64(dx*dx + dy*dy))
					dirX := float64(dx) / distance
					dirY := float64(dy) / distance
					pill.End.X = pill.Start.X + int(dirX * length)
					pill.End.Y = pill.Start.Y + int(dirY * length)
				} else {
					
					pill.End.X = pill.Start.X + int(h.UI.PillLengthSlider.Value)
					pill.End.Y = pill.Start.Y
				}
				
				h.UI.StatusLabel.SetText("Pill radius set. Use slider to adjust length, click to finalize.")
				h.UI.Canvas.Refresh()
				return
			} else if pill.Step == 2 {
				
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
		
	case "rectangle":
		rectangle := models.NewRectangle(
			h.StartPoint,
			h.StartPoint, // Initially both corners are the same
			h.UI.State.CurrentColor,
			h.UI.State.BrushThickness,
		)
		
		// Apply fill settings if enabled
		if h.UI.State.FillEnabled {
			if h.UI.State.UseImageFill && h.UI.State.FillImage != nil {
				rectangle.SetFillImage(h.UI.State.FillImage)
			} else if h.UI.State.FillColor != nil {
				rectangle.SetFillColor(h.UI.State.FillColor)
			}
		}
		
		h.UI.State.CurrentShape = rectangle
		h.UI.StatusLabel.SetText("Drawing rectangle... Release to complete")
	}
}



func (h *MouseHandler) MouseUp(ev *desktop.MouseEvent) {
	
	if h.IsMoving && h.UI.State.SelectedShape != nil {
		h.IsMoving = false
		h.UI.StatusLabel.SetText("Shape moved.")
		h.UI.Canvas.Refresh()
		return
	}
	
	// Handle clipping action
	if h.UI.State.CurrentAction == "clipping" && !h.IsDrawing {
		adjustedPoint := h.adjustMousePosition(ev.PointEvent)
		
		// Find which shape was clicked
		for i := len(h.UI.State.Shapes) - 1; i >= 0; i-- {
			shape := h.UI.State.Shapes[i]
			if shape.Contains(adjustedPoint) {
				// Check if it's a polygon
				polygon, isPolygon := shape.(*models.Polygon)
				if !isPolygon {
					h.UI.StatusLabel.SetText("Clipping only works with polygons. Please select a polygon.")
					return
				}
				
				// Get the selected polygon (clipper)
				selectedPoly, _ := h.UI.State.SelectedShape.(*models.Polygon)
				
				// Check if the polygon to be clipped is convex
				if !polygon.IsConvex() {
					h.UI.StatusLabel.SetText("Only convex polygons can be clipped.")
					return
				}
				
				// Perform clipping using our utility function
				clippedVertices := ClipPolygon(polygon.GetVertices(), selectedPoly.GetVertices())
				
				// Create new polygon with clipped vertices
				if len(clippedVertices) >= 3 {
					clippedPoly := models.NewPolygon(clippedVertices, polygon.GetColor(), polygon.Thickness)
					
					// Copy fill properties
					if polygon.IsFilled {
						if polygon.UseImage {
							clippedPoly.SetFillImage(polygon.FillImage)
						} else {
							clippedPoly.SetFillColor(polygon.FillColor)
						}
					}
					
					// Add the clipped polygon to the shapes
					h.UI.State.Shapes = append(h.UI.State.Shapes, clippedPoly)
					h.UI.Canvas.Refresh()
					h.UI.StatusLabel.SetText("Polygon clipped successfully.")
				} else {
					h.UI.StatusLabel.SetText("Clipping result is not a valid polygon.")
				}
				
				return
			}
		}
		
		return
	}

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
		case "rectangle":
			h.UI.StatusLabel.SetText("Rectangle added")
		}
	}
}




func (h *MouseHandler) MouseMoved(ev *desktop.MouseEvent) {
	h.CurrentPoint = h.adjustMousePosition(ev.PointEvent)
	
	
	if h.IsMoving && h.UI.State.SelectedShape != nil {
		
		deltaX := h.CurrentPoint.X - h.MoveStartX
		deltaY := h.CurrentPoint.Y - h.MoveStartY
		
		
		if deltaX != 0 || deltaY != 0 {
			h.UI.State.SelectedShape.Move(deltaX, deltaY)
			
			
			h.MoveStartX = h.CurrentPoint.X
			h.MoveStartY = h.CurrentPoint.Y
			
			h.UI.Canvas.Refresh()
		}
		return
	}
	
	
	if !h.IsDrawing || h.UI.State.CurrentShape == nil {
		return
	}
	
	switch shape := h.UI.State.CurrentShape.(type) {
	case *models.Line:
		shape.End = h.CurrentPoint
		h.UI.Canvas.Refresh()
		
	case *models.Circle:
		dx := h.CurrentPoint.X - shape.Center.X
		dy := h.CurrentPoint.Y - shape.Center.Y
		shape.Radius = int(math.Sqrt(float64(dx*dx + dy*dy)))
		h.UI.Canvas.Refresh()
		
	case *models.Rectangle:
		// Update the bottom-right corner as the mouse moves
		shape.BottomRight = h.CurrentPoint
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
		
		// Apply fill settings if enabled
		if h.UI.State.FillEnabled {
			if h.UI.State.UseImageFill && h.UI.State.FillImage != nil {
				poly.SetFillImage(h.UI.State.FillImage)
			} else if h.UI.State.FillColor != nil {
				poly.SetFillColor(h.UI.State.FillColor)
			}
		}
		
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
	
	if h.IsMoving && h.UI.State.SelectedShape != nil {
		h.IsMoving = false
		h.UI.StatusLabel.SetText("Shape moved.")
		h.UI.Canvas.Refresh()
	}
}


func (h *MouseHandler) Tapped(ev *fyne.PointEvent) {
	h.MouseDown(&desktop.MouseEvent{
		PointEvent: *ev,
		Button:     desktop.MouseButtonPrimary,
	})
	
	
	if h.UI.State.CurrentAction == "polygon" || h.UI.State.CurrentAction == "pill" {
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