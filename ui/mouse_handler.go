package ui

import (
	"image" // Added import for image
	"image/color"
	"math"
	"paint-drawer-pro/algorithms" // Added import for algorithms
	"paint-drawer-pro/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ResizePoint represents which control point is being dragged
type ResizePoint int

const (
	None ResizePoint = iota
	TopLeft
	TopRight
	BottomRight
	BottomLeft
)

type MouseHandler struct {
	widget.BaseWidget
	UI                *MainUI
	StartPoint        models.Point
	CurrentPoint      models.Point
	IsDrawing         bool
	LastPoint         models.Point
	PolyPoints        []models.Point
	IsMoving          bool      
	IsResizing        bool
	CurrentResizePoint ResizePoint
	MoveStartX        int       
	MoveStartY        int       
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
	
	// Reset resize state
	h.IsResizing = false
	h.CurrentResizePoint = None

	if h.UI.State.CurrentAction == "scanline_fill" && ev.Button == desktop.MouseButtonPrimary {
		if h.UI.State.FillStage == "" || h.UI.State.FillStage == "area_selected" { // Start selecting area or restart selection
			// Clear previous selection rect if any visual remnants by redrawing
			if h.UI.State.SelectionRect != nil {
				h.UI.State.SelectionRect = nil
				h.UI.Canvas.Refresh() // Clear old one first
			}
			h.UI.State.SelectionRect = models.NewRectangle(adjustedPoint, adjustedPoint, color.RGBA{0, 0, 255, 100}, 1) // Blue, slightly transparent
			// It's better if fill color for selection is not set here, but handled by drawing logic for consistency
			// h.UI.State.SelectionRect.SetFillColor(color.RGBA{0, 0, 255, 50})
			h.UI.State.FillStage = "selecting_area"
			h.IsDrawing = true // Use IsDrawing to indicate selection rectangle drawing in progress
			h.UI.StatusLabel.SetText("Drag to select area for scanline fill. Release to confirm selection.")
			h.UI.Canvas.Refresh()
			return
		} else if h.UI.State.FillStage == "awaiting_fill_point" && h.UI.State.SelectionRect != nil {
			normalizedSelectionRect := *h.UI.State.SelectionRect // Make a copy for safety
			normalizedSelectionRect.Normalize()

			if normalizedSelectionRect.Contains(adjustedPoint) {
				// Area selected, and click is inside. Perform fill.
				canvasSize := h.UI.Canvas.Size()
				currentRenderedImage := h.UI.renderCanvas(int(canvasSize.Width), int(canvasSize.Height))

				bounds := currentRenderedImage.Bounds()
				width, height := bounds.Max.X, bounds.Max.Y
				
				originalCanvasData := make([][]color.Color, height)
				for y := 0; y < height; y++ {
					originalCanvasData[y] = make([]color.Color, width)
					for x := 0; x < width; x++ {
						originalCanvasData[y][x] = currentRenderedImage.At(x, y)
					}
				}

				// Create a separate buffer for the fill algorithm to operate on.
				fillBufferForAlgorithm := make([][]color.Color, height)
				for y := 0; y < height; y++ {
					fillBufferForAlgorithm[y] = make([]color.Color, width)
					copy(fillBufferForAlgorithm[y], originalCanvasData[y])
				}

				boundaryColor := color.RGBA{0, 0, 0, 255} // Assuming black boundary for now.
				fillColor := h.UI.State.FillColor
				if h.UI.State.UseImageFill {
					h.UI.StatusLabel.SetText("Scanline fill does not support image patterns. Using selected fill color.")
					if fillColor == nil {
						fillColor = color.RGBA{100, 100, 100, 255} // Default fill if none selected
					}
				}

				algPoint := algorithms.Point{X: adjustedPoint.X, Y: adjustedPoint.Y}
				algorithms.SmithScanlineFill(fillBufferForAlgorithm, algPoint, fillColor, boundaryColor)

				// Composite the filled area (from fillBufferForAlgorithm) onto the original image (originalCanvasData)
				// to create the final image, respecting the selection rectangle bounds.
				finalImage := image.NewRGBA(image.Rect(0, 0, width, height))
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						if x >= normalizedSelectionRect.TopLeft.X && x <= normalizedSelectionRect.BottomRight.X &&
						   y >= normalizedSelectionRect.TopLeft.Y && y <= normalizedSelectionRect.BottomRight.Y {
							finalImage.Set(x, y, fillBufferForAlgorithm[y][x])
						} else {
							finalImage.Set(x, y, originalCanvasData[y][x])
						}
					}
				}

				// Update the MainUI's base image for the canvas
				// Assumes MainUI has a field like 'BaseImage *image.RGBA' that its raster generator uses.
				h.UI.BaseImage = finalImage 

				h.UI.State.FillStage = "area_selected" // Indicate fill is done, ready for new selection or tool change
				h.UI.State.SelectionRect = nil      // Clear the selection rectangle from state
				h.UI.Canvas.Refresh()               // Refresh canvas to show the filled area and remove selection rect
				h.UI.StatusLabel.SetText("Area filled. Click to select new area or choose another tool.")
				return
			} else {
				// Clicked outside the selection rectangle
				h.UI.StatusLabel.SetText("Clicked outside selection. Click inside the blue rectangle to fill, or click and drag to reselect area.")
				// Optionally, reset to "selecting_area" or "" stage if a click outside means cancel current selection
				// h.UI.State.FillStage = ""
				// h.UI.State.SelectionRect = nil
				// h.UI.Canvas.Refresh()
				return
			}
		}
		// If FillStage is "selecting_area", MouseMove and MouseUp will handle it.
		// If FillStage is "awaiting_fill_point" but SelectionRect is nil (should not happen), ignore.
		return // Prevent other MouseDown actions when scanline_fill is active and in a specific stage
	}
	
	// ... (rest of MouseDown, e.g., for "select", "polygon", "pill", "line", etc.) ...
	// Ensure that if IsDrawing was set true for scanline_fill selection, it doesn't interfere here.
	// The return statements above should handle this for scanline_fill stages.

	if h.UI.State.CurrentAction == "select" && ev.Button == desktop.MouseButtonPrimary {
		// Check if we're clicking on a resize handle of the currently selected rectangle
		if h.UI.State.SelectedShape != nil {
			if rect, isRect := h.UI.State.SelectedShape.(*models.Rectangle); isRect {
				resizePoint := rect.GetResizePointAt(adjustedPoint)
				if resizePoint != models.None {
					h.IsResizing = true
					h.CurrentResizePoint = ResizePoint(resizePoint)
					h.UI.StatusLabel.SetText("Resizing rectangle...")
					return
				}
			}
		}
		
		// If not resizing, proceed with normal selection
		h.UI.State.SelectedShape = nil
		h.UI.PillLengthContainer.Hide()
		
		// Find if we clicked on a shape
		for i := len(h.UI.State.Shapes) - 1; i >= 0; i-- {
			shape := h.UI.State.Shapes[i]
			if shape.Contains(adjustedPoint) {
				h.UI.State.SelectedShape = shape
				h.IsMoving = true
				h.MoveStartX = adjustedPoint.X
				h.MoveStartY = adjustedPoint.Y
				
				// Handle special case for pill shapes
				if pill, isPill := shape.(*models.Pill); isPill {
					dx := pill.End.X - pill.Start.X
					dy := pill.End.Y - pill.Start.Y
					length := math.Sqrt(float64(dx*dx + dy*dy))
					h.UI.PillLengthSlider.SetValue(length)
					h.UI.PillLengthContainer.Show()
					h.UI.StatusLabel.SetText("Pill selected. Use slider to adjust length or drag to move.")
				} else if _, isRect := shape.(*models.Rectangle); isRect {
					h.UI.StatusLabel.SetText("Rectangle selected. Drag corners to resize or drag center to move. Press Delete to remove.")
				} else {
					h.UI.StatusLabel.SetText("Shape selected. Drag to move. Press Delete to remove.")
				}
				
				h.UI.Canvas.Refresh()
				return
			}
		}
		
		h.UI.StatusLabel.SetText("No shape selected.")
		return
	}

	if h.UI.State.CurrentAction == "scanline_fill" && ev.Button == desktop.MouseButtonPrimary {
		canvasImg := h.UI.renderCanvas(int(h.UI.Canvas.Size().Width), int(h.UI.Canvas.Size().Height))
		bounds := canvasImg.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y
		canvasData := make([][]color.Color, height)
		for y := 0; y < height; y++ {
			canvasData[y] = make([]color.Color, width)
			for x := 0; x < width; x++ {
				canvasData[y][x] = canvasImg.At(x, y)
			}
		}

		boundaryColor := color.RGBA{0, 0, 0, 255}
		fillColor := h.UI.State.FillColor
		if h.UI.State.UseImageFill { 
			h.UI.StatusLabel.SetText("Scanline fill does not support image patterns. Using selected fill color.")
			if fillColor == nil {
				fillColor = color.RGBA{100,100,100,255}
			}
		}

		 algPoint := algorithms.Point{X: adjustedPoint.X, Y: adjustedPoint.Y} 
		algorithms.SmithScanlineFill(canvasData, algPoint, fillColor, boundaryColor)


		newImg := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				newImg.Set(x, y, canvasData[y][x])
			}
		}
		

  	h.UI.Canvas.Refresh()

		h.UI.StatusLabel.SetText("Area filled using Scanline Fill.")
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
	adjustedPoint := h.adjustMousePosition(ev.PointEvent)

	if h.UI.State.CurrentAction == "scanline_fill" && h.UI.State.FillStage == "selecting_area" {
		if h.IsDrawing && h.UI.State.SelectionRect != nil {
			h.UI.State.SelectionRect.BottomRight = adjustedPoint
			h.UI.State.SelectionRect.Normalize() // Ensure TopLeft and BottomRight are correct
			h.IsDrawing = false
			h.UI.State.FillStage = "awaiting_fill_point"
			h.UI.StatusLabel.SetText("Area selected. Click inside the blue rectangle to pick a fill start point.")
			h.UI.Canvas.Refresh() // Refresh to show final selection rect and await click
			return
		}
	}
	
	if h.IsResizing && h.UI.State.SelectedShape != nil {
		h.IsResizing = false
		h.CurrentResizePoint = None
		h.UI.StatusLabel.SetText("Rectangle resized.")
		h.UI.Canvas.Refresh()
		return
	}
	
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
				
				if !selectedPoly.IsConvex() {
					h.UI.StatusLabel.SetText("Only convex polygons can be used as clippers.")
					return
				}
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

	if h.UI.State.CurrentAction == "scanline_fill" && h.UI.State.FillStage == "selecting_area" && h.IsDrawing {
		if h.UI.State.SelectionRect != nil {
			h.UI.State.SelectionRect.BottomRight = h.CurrentPoint
			h.UI.Canvas.Refresh() // Refresh to show selection rectangle being drawn
			return // Exclusive handling for selection drawing
		}
	}
	
	// Handle resizing of rectangle
	if h.IsResizing && h.UI.State.SelectedShape != nil {
		if rect, isRect := h.UI.State.SelectedShape.(*models.Rectangle); isRect {
			resizePoint := models.ResizePointType(h.CurrentResizePoint)
			rect.ResizeByCorner(resizePoint, h.CurrentPoint)
			h.UI.Canvas.Refresh()
		}
		return
	}
	
	// Handle moving shapes
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
	// Handle shape deletion with Delete or Backspace keys
	if (ev.Name == fyne.KeyDelete || ev.Name == fyne.KeyBackspace) && h.UI.State.CurrentAction == "select" && h.UI.State.SelectedShape != nil {
		// Find and remove the selected shape
		for i, shape := range h.UI.State.Shapes {
			if shape == h.UI.State.SelectedShape {
				// Remove shape from the slice
				h.UI.State.Shapes = append(h.UI.State.Shapes[:i], h.UI.State.Shapes[i+1:]...)
				h.UI.State.SelectedShape = nil
				h.UI.Canvas.Refresh()
				h.UI.StatusLabel.SetText("Shape deleted")
				break
			}
		}
		return
	}
	
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
	
	if h.IsResizing && h.UI.State.SelectedShape != nil {
		h.IsResizing = false
		h.CurrentResizePoint = None
		h.UI.StatusLabel.SetText("Rectangle resized.")
		h.UI.Canvas.Refresh()
		return
	}
	
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
	// Get the position of the canvas within the window
	canvasPos := h.UI.Canvas.Position()
	
	// Calculate the position relative to the canvas's position
	// by subtracting the canvas's position from the absolute mouse position
	x := int(ev.Position.X - canvasPos.X)
	y := int(ev.Position.Y - canvasPos.Y)
	
	// Get canvas size for bounds checking
	canvasSize := h.UI.Canvas.Size()
	maxX := int(canvasSize.Width) - 1
	maxY := int(canvasSize.Height) - 1
	
	// Constrain to canvas bounds
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