package models

import (
	"image/color"
	"math"
	"paint-drawer-pro/algorithms"
)


func NewPolygon(vertices []Point, color color.Color, thickness int) *Polygon {
	if thickness <= 0 {
		thickness = 1
	}
	return &Polygon{
		Vertices:  vertices,
		Color:     color,
		Thickness: thickness,
		IsFilled:  false,
	}
}


func (p *Polygon) Draw(canvas [][]color.Color, antiAliasing bool) {
	if len(p.Vertices) < 3 {
		return 
	}

	// Draw fill if enabled
	if p.IsFilled {
		p.drawFill(canvas)
	}
	
	// Draw the outline
	for i := 0; i < len(p.Vertices); i++ {
		start := p.Vertices[i]
		end := p.Vertices[(i+1)%len(p.Vertices)]

		if antiAliasing {
			drawXiaolinWuLine(canvas, start.X, start.Y, end.X, end.Y, p.Color)
		} else {
			if p.Thickness == 1 {
				drawMidpointLine(canvas, start.X, start.Y, end.X, end.Y, p.Color)
			} else {
				drawThickLine(canvas, start.X, start.Y, end.X, end.Y, p.Color, p.Thickness)
			}
		}
	}
}


func (p *Polygon) drawFill(canvas [][]color.Color) {
	// Convert vertices to algorithms.Point type
	algVertices := make([]algorithms.Point, len(p.Vertices))
	for i, v := range p.Vertices {
		algVertices[i] = algorithms.Point{X: v.X, Y: v.Y}
	}
	
	// First simplify the polygon to remove duplicate or near-duplicate vertices
	simplifiedVertices := algorithms.SimplifyPolygon(algVertices, 2.0)
	
	if p.UseImage && p.FillImage != nil {
		// Use the more robust parity-based fill algorithm for image filling
		algorithms.ParityFillPolygonWithImage(canvas, simplifiedVertices, p.FillImage)
	} else {
		// Use the more robust parity-based fill algorithm for solid color filling
		algorithms.ParityFillPolygon(canvas, simplifiedVertices, p.FillColor)
	}
}


func (p *Polygon) Contains(pt Point) bool {
	
	for _, vertex := range p.Vertices {
		dx := vertex.X - pt.X
		dy := vertex.Y - pt.Y
		if dx*dx+dy*dy <= 25 { 
			return true
		}
	}

	
	for i := 0; i < len(p.Vertices); i++ {
		start := p.Vertices[i]
		end := p.Vertices[(i+1)%len(p.Vertices)]


		lineLen := math.Sqrt(float64((end.X-start.X)*(end.X-start.X) + (end.Y-start.Y)*(end.Y-start.Y)))
		if lineLen == 0 {
			continue
		}

		t := float64((pt.X-start.X)*(end.X-start.X) + (pt.Y-start.Y)*(end.Y-start.Y)) / (lineLen * lineLen)
		if t < 0 || t > 1 {
			continue
		}

		nearestX := start.X + int(float64(end.X-start.X)*t)
		nearestY := start.Y + int(float64(end.Y-start.Y)*t)

		dist := math.Sqrt(float64((pt.X-nearestX)*(pt.X-nearestX) + (pt.Y-nearestY)*(pt.Y-nearestY)))
		if dist <= float64(p.Thickness+5) {
			return true
		}
	}

	return false
}


func (p *Polygon) GetControlPoints() []Point {
	return p.Vertices
}


func (p *Polygon) Move(deltaX, deltaY int) {
	for i := range p.Vertices {
		p.Vertices[i].X += deltaX
		p.Vertices[i].Y += deltaY
	}
}


func (p *Polygon) SetColor(color color.Color) {
	p.Color = color
}


func (p *Polygon) GetColor() color.Color {
	return p.Color
}


func (p *Polygon) Serialize() map[string]interface{} {
	serMap := map[string]interface{}{
		"type":      "polygon",
		"thickness": p.Thickness,
		"isFilled":  p.IsFilled,
		"useImage":  p.UseImage,
	}
	
	// Serialize vertices
	vertices := make([]map[string]interface{}, len(p.Vertices))
	for i, vertex := range p.Vertices {
		vertices[i] = map[string]interface{}{
			"X": vertex.X,
			"Y": vertex.Y,
		}
	}
	serMap["vertices"] = vertices
	
	// Serialize colors
	if p.Color != nil {
		r, g, b, a := p.Color.RGBA()
		serMap["color"] = map[string]interface{}{
			"R": uint8(r),
			"G": uint8(g),
			"B": uint8(b),
			"A": uint8(a),
		}
	}
	
	if p.IsFilled && !p.UseImage && p.FillColor != nil {
		r, g, b, a := p.FillColor.RGBA()
		serMap["fillColor"] = map[string]interface{}{
			"R": uint8(r),
			"G": uint8(g),
			"B": uint8(b),
			"A": uint8(a),
		}
	}
	
	// Image data would be too large to serialize directly
	// Consider saving it to a file instead
	
	return serMap
}


func (p *Polygon) Clone() Shape {
	vertices := make([]Point, len(p.Vertices))
	for i, vertex := range p.Vertices {
		vertices[i] = Point{X: vertex.X, Y: vertex.Y}
	}
	
	clone := NewPolygon(vertices, p.Color, p.Thickness)
	clone.FillColor = p.FillColor
	clone.IsFilled = p.IsFilled
	clone.UseImage = p.UseImage
	
	// Deep copy fill image if present
	if p.UseImage && p.FillImage != nil {
		height := len(p.FillImage)
		if height > 0 {
			width := len(p.FillImage[0])
			clone.FillImage = make([][]color.Color, height)
			for y := 0; y < height; y++ {
				clone.FillImage[y] = make([]color.Color, width)
				for x := 0; x < width; x++ {
					clone.FillImage[y][x] = p.FillImage[y][x]
				}
			}
		}
	}
	
	return clone
}


func (p *Polygon) SetFillColor(c color.Color) {
	p.FillColor = c
	p.IsFilled = true
	p.UseImage = false
}


func (p *Polygon) SetFillImage(img [][]color.Color) {
	p.FillImage = img
	p.IsFilled = true
	p.UseImage = true
}


func (p *Polygon) DisableFill() {
	p.IsFilled = false
}



func (p *Polygon) IsConvex() bool {
	if len(p.Vertices) <= 3 {
		// Triangles are always convex as long as they have distinct points
		return len(p.Vertices) == 3
	}
	
	// Add debugging information to show the number of vertices
	vertices := p.Vertices
	
	// Convert to algorithm points
	algVertices := make([]algorithms.Point, len(vertices))
	for i, v := range vertices {
		algVertices[i] = algorithms.Point{X: v.X, Y: v.Y}
	}
	
	// First simplify the polygon
	simplified := algorithms.SimplifyPolygon(algVertices, 2.0)
	
	// Then check if it's convex
	return algorithms.IsPolygonConvex(simplified)
}


func (p *Polygon) GetVertices() []Point {
	return p.Vertices
}