package models

import (
	"image/color"
	"math"
)

// Rectangle represents a rectangle shape
type Rectangle struct {
	TopLeft     Point
	BottomRight Point
	Color       color.Color
	Thickness   int
	FillColor   color.Color
	IsFilled    bool
	FillImage   [][]color.Color
	UseImage    bool
}

// NewRectangle creates a new rectangle with the given corner points and color
func NewRectangle(topLeft, bottomRight Point, color color.Color, thickness int) *Rectangle {
	if thickness <= 0 {
		thickness = 1
	}
	
	// Make sure topLeft is actually top-left and bottomRight is bottom-right
	x1, y1 := topLeft.X, topLeft.Y
	x2, y2 := bottomRight.X, bottomRight.Y
	
	// Swap if necessary
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	
	return &Rectangle{
		TopLeft:     Point{X: x1, Y: y1},
		BottomRight: Point{X: x2, Y: y2},
		Color:       color,
		Thickness:   thickness,
		IsFilled:    false,
	}
}

// Draw draws the rectangle on the canvas
func (r *Rectangle) Draw(canvas [][]color.Color, antiAliasing bool) {
	// Draw filled rectangle first if enabled
	if r.IsFilled {
		r.drawFill(canvas)
	}
	
	// Draw the rectangle outline
	topRight := Point{X: r.BottomRight.X, Y: r.TopLeft.Y}
	bottomLeft := Point{X: r.TopLeft.X, Y: r.BottomRight.Y}
	
	// Draw four edges of the rectangle
	if antiAliasing {
		drawXiaolinWuLine(canvas, r.TopLeft.X, r.TopLeft.Y, topRight.X, topRight.Y, r.Color)
		drawXiaolinWuLine(canvas, topRight.X, topRight.Y, r.BottomRight.X, r.BottomRight.Y, r.Color)
		drawXiaolinWuLine(canvas, r.BottomRight.X, r.BottomRight.Y, bottomLeft.X, bottomLeft.Y, r.Color)
		drawXiaolinWuLine(canvas, bottomLeft.X, bottomLeft.Y, r.TopLeft.X, r.TopLeft.Y, r.Color)
	} else {
		if r.Thickness == 1 {
			drawMidpointLine(canvas, r.TopLeft.X, r.TopLeft.Y, topRight.X, topRight.Y, r.Color)
			drawMidpointLine(canvas, topRight.X, topRight.Y, r.BottomRight.X, r.BottomRight.Y, r.Color)
			drawMidpointLine(canvas, r.BottomRight.X, r.BottomRight.Y, bottomLeft.X, bottomLeft.Y, r.Color)
			drawMidpointLine(canvas, bottomLeft.X, bottomLeft.Y, r.TopLeft.X, r.TopLeft.Y, r.Color)
		} else {
			drawThickLine(canvas, r.TopLeft.X, r.TopLeft.Y, topRight.X, topRight.Y, r.Color, r.Thickness)
			drawThickLine(canvas, topRight.X, topRight.Y, r.BottomRight.X, r.BottomRight.Y, r.Color, r.Thickness)
			drawThickLine(canvas, r.BottomRight.X, r.BottomRight.Y, bottomLeft.X, bottomLeft.Y, r.Color, r.Thickness)
			drawThickLine(canvas, bottomLeft.X, bottomLeft.Y, r.TopLeft.X, r.TopLeft.Y, r.Color, r.Thickness)
		}
	}
}

// drawFill fills the rectangle with a solid color or an image
func (r *Rectangle) drawFill(canvas [][]color.Color) {
	startX := r.TopLeft.X + 1
	endX := r.BottomRight.X - 1
	startY := r.TopLeft.Y + 1
	endY := r.BottomRight.Y - 1
	
	// Don't try to fill if rectangle is too small
	if startX >= endX || startY >= endY {
		return
	}
	
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			if x >= 0 && y >= 0 && y < len(canvas) && x < len(canvas[0]) {
				if r.UseImage && r.FillImage != nil {
					// Map the pixel to the image coordinates
					imgY := (y - startY) % len(r.FillImage)
					imgX := (x - startX) % len(r.FillImage[0])
					if imgY >= 0 && imgX >= 0 && imgY < len(r.FillImage) && imgX < len(r.FillImage[0]) {
						canvas[y][x] = r.FillImage[imgY][imgX]
					}
				} else {
					canvas[y][x] = r.FillColor
				}
			}
		}
	}
}

// Contains checks if a point is contained within or on the rectangle
func (r *Rectangle) Contains(p Point) bool {
	// Check if point is near any of the corners (for vertex selection)
	corners := []Point{
		r.TopLeft, 
		{X: r.BottomRight.X, Y: r.TopLeft.Y},     // Top Right
		r.BottomRight, 
		{X: r.TopLeft.X, Y: r.BottomRight.Y},     // Bottom Left
	}
	
	for _, corner := range corners {
		dx := corner.X - p.X
		dy := corner.Y - p.Y
		if dx*dx+dy*dy <= 25 { // 5*5 = 25 pixel radius for selection
			return true
		}
	}
	
	// Check if point is near any of the edges
	// Top edge
	if p.X >= r.TopLeft.X && p.X <= r.BottomRight.X && 
	   math.Abs(float64(p.Y-r.TopLeft.Y)) <= float64(r.Thickness+5) {
		return true
	}
	
	// Right edge
	if p.Y >= r.TopLeft.Y && p.Y <= r.BottomRight.Y && 
	   math.Abs(float64(p.X-r.BottomRight.X)) <= float64(r.Thickness+5) {
		return true
	}
	
	// Bottom edge
	if p.X >= r.TopLeft.X && p.X <= r.BottomRight.X && 
	   math.Abs(float64(p.Y-r.BottomRight.Y)) <= float64(r.Thickness+5) {
		return true
	}
	
	// Left edge
	if p.Y >= r.TopLeft.Y && p.Y <= r.BottomRight.Y && 
	   math.Abs(float64(p.X-r.TopLeft.X)) <= float64(r.Thickness+5) {
		return true
	}
	
	// Check if point is inside the rectangle (when filled)
	if r.IsFilled && p.X > r.TopLeft.X && p.X < r.BottomRight.X && 
	   p.Y > r.TopLeft.Y && p.Y < r.BottomRight.Y {
		return true
	}
	
	return false
}

// GetControlPoints returns the control points of the rectangle (four corners)
func (r *Rectangle) GetControlPoints() []Point {
	return []Point{
		r.TopLeft,
		{X: r.BottomRight.X, Y: r.TopLeft.Y},     // Top Right
		r.BottomRight,
		{X: r.TopLeft.X, Y: r.BottomRight.Y},     // Bottom Left
	}
}

// Move moves the rectangle by the specified delta
func (r *Rectangle) Move(deltaX, deltaY int) {
	r.TopLeft.X += deltaX
	r.TopLeft.Y += deltaY
	r.BottomRight.X += deltaX
	r.BottomRight.Y += deltaY
}

// SetColor sets the color of the rectangle outline
func (r *Rectangle) SetColor(c color.Color) {
	r.Color = c
}

// GetColor returns the color of the rectangle outline
func (r *Rectangle) GetColor() color.Color {
	return r.Color
}

// SetFillColor sets the fill color of the rectangle
func (r *Rectangle) SetFillColor(c color.Color) {
	r.FillColor = c
	r.IsFilled = true
	r.UseImage = false
}

// SetFillImage sets an image to fill the rectangle
func (r *Rectangle) SetFillImage(img [][]color.Color) {
	r.FillImage = img
	r.IsFilled = true
	r.UseImage = true
}

// DisableFill disables filling the rectangle
func (r *Rectangle) DisableFill() {
	r.IsFilled = false
}

// IsConvex returns true since rectangles are always convex
func (r *Rectangle) IsConvex() bool {
	return true
}

// GetVertices returns the vertices of the rectangle as a slice of Points
func (r *Rectangle) GetVertices() []Point {
	return []Point{
		r.TopLeft,
		{X: r.BottomRight.X, Y: r.TopLeft.Y},     // Top Right
		r.BottomRight,
		{X: r.TopLeft.X, Y: r.BottomRight.Y},     // Bottom Left
	}
}

// Serialize converts the rectangle to a map for serialization
func (r *Rectangle) Serialize() map[string]interface{} {
	serMap := map[string]interface{}{
		"type":      "rectangle",
		"thickness": r.Thickness,
		"isFilled":  r.IsFilled,
		"useImage":  r.UseImage,
	}
	
	// Serialize points
	serMap["topLeft"] = map[string]interface{}{
		"X": r.TopLeft.X,
		"Y": r.TopLeft.Y,
	}
	
	serMap["bottomRight"] = map[string]interface{}{
		"X": r.BottomRight.X,
		"Y": r.BottomRight.Y,
	}
	
	// Serialize colors
	if r.Color != nil {
		r, g, b, a := r.Color.RGBA()
		serMap["color"] = map[string]interface{}{
			"R": uint8(r),
			"G": uint8(g),
			"B": uint8(b),
			"A": uint8(a),
		}
	}
	
	if r.IsFilled && !r.UseImage && r.FillColor != nil {
		fr, fg, fb, fa := r.FillColor.RGBA()
		serMap["fillColor"] = map[string]interface{}{
			"R": uint8(fr),
			"G": uint8(fg),
			"B": uint8(fb),
			"A": uint8(fa),
		}
	}
	
	// Image data would be too large to serialize directly
	// Consider saving it to a file instead
	
	return serMap
}

// Clone creates a copy of the rectangle
func (r *Rectangle) Clone() Shape {
	newRect := &Rectangle{
		TopLeft:     Point{X: r.TopLeft.X, Y: r.TopLeft.Y},
		BottomRight: Point{X: r.BottomRight.X, Y: r.BottomRight.Y},
		Color:       r.Color,
		Thickness:   r.Thickness,
		FillColor:   r.FillColor,
		IsFilled:    r.IsFilled,
		UseImage:    r.UseImage,
	}
	
	if r.UseImage && r.FillImage != nil {
		// Deep copy the fill image
		height := len(r.FillImage)
		if height > 0 {
			width := len(r.FillImage[0])
			newRect.FillImage = make([][]color.Color, height)
			for y := 0; y < height; y++ {
				newRect.FillImage[y] = make([]color.Color, width)
				for x := 0; x < width; x++ {
					newRect.FillImage[y][x] = r.FillImage[y][x]
				}
			}
		}
	}
	
	return newRect
}
