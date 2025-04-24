package models

import (
	"image/color"
	"math"
)


func NewPolygon(vertices []Point, color color.Color, thickness int) *Polygon {
	if thickness <= 0 {
		thickness = 1
	}
	return &Polygon{
		Vertices:  vertices,
		Color:     color,
		Thickness: thickness,
	}
}


func (p *Polygon) Draw(canvas [][]color.Color, antiAliasing bool) {
	if len(p.Vertices) < 3 {
		return 
	}

	
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
	r, g, b, a := p.Color.RGBA()
	
	vertices := make([]int, len(p.Vertices)*2)
	for i, vertex := range p.Vertices {
		vertices[i*2] = vertex.X
		vertices[i*2+1] = vertex.Y
	}
	
	return map[string]interface{}{
		"type":      "polygon",
		"vertices":  vertices,
		"color":     []uint32{r, g, b, a},
		"thickness": p.Thickness,
	}
}


func (p *Polygon) Clone() Shape {
	vertices := make([]Point, len(p.Vertices))
	for i, vertex := range p.Vertices {
		vertices[i] = Point{X: vertex.X, Y: vertex.Y}
	}
	
	return NewPolygon(vertices, p.Color, p.Thickness)
}