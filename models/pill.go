package models

import (
	"image/color"
	"math"
	"paint-drawer-pro/algorithms"
)


type Pill struct {
	Start    Point     
	End      Point     
	Radius   int       
	Color    color.Color
	Step     int       
}


func NewPill(start Point, radius int, color color.Color) *Pill {
	return &Pill{
		Start:   start,
		End:     start, 
		Radius:  radius,
		Color:   color,
		Step:    1, 
	}
}


func (p *Pill) Draw(canvas [][]color.Color, antiAliasing bool) {
	
	if p.Step == 1 {	
		drawMidpointCircle(canvas, p.Start.X, p.Start.Y, 5, p.Color)
		return
	}
	
	
	if p.Step == 2 {
		if antiAliasing {
			drawXiaolinWuCircle(canvas, p.Start.X, p.Start.Y, p.Radius, p.Color)
		} else {
			drawMidpointCircle(canvas, p.Start.X, p.Start.Y, p.Radius, p.Color)
		}
		return
	}

	
	dx := p.End.X - p.Start.X
	dy := p.End.Y - p.Start.Y
	length := math.Sqrt(float64(dx*dx + dy*dy))
	
	
	if length < float64(p.Radius) {
		if antiAliasing {
			drawXiaolinWuCircle(canvas, p.Start.X, p.Start.Y, p.Radius, p.Color)
		} else {
			drawMidpointCircle(canvas, p.Start.X, p.Start.Y, p.Radius, p.Color)
		}
		return
	}
	
	
	dirX := float64(dx) / length
	dirY := float64(dy) / length
	
	
	perpX := -dirY
	perpY := dirX

	
	
	rectStartX := p.Start.X + int(dirX*float64(p.Radius))
	rectStartY := p.Start.Y + int(dirY*float64(p.Radius))
	
	
	rectEndX := p.End.X - int(dirX*float64(p.Radius))
	rectEndY := p.End.Y - int(dirY*float64(p.Radius))
	
	
	topLeft := Point{
		X: rectStartX + int(perpX*float64(p.Radius)),
		Y: rectStartY + int(perpY*float64(p.Radius)),
	}
	
	bottomLeft := Point{
		X: rectStartX - int(perpX*float64(p.Radius)),
		Y: rectStartY - int(perpY*float64(p.Radius)),
	}
	
	topRight := Point{
		X: rectEndX + int(perpX*float64(p.Radius)),
		Y: rectEndY + int(perpY*float64(p.Radius)),
	}
	
	bottomRight := Point{
		X: rectEndX - int(perpX*float64(p.Radius)),
		Y: rectEndY - int(perpY*float64(p.Radius)),
	}
	
	
	minX := min(min(topLeft.X, topRight.X), min(bottomLeft.X, bottomRight.X))
	maxX := max(max(topLeft.X, topRight.X), max(bottomLeft.X, bottomRight.X))
	minY := min(min(topLeft.Y, topRight.Y), min(bottomLeft.Y, bottomRight.Y))
	maxY := max(max(topLeft.Y, topRight.Y), max(bottomLeft.Y, bottomRight.Y))

	
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			
			if isPointInRect(Point{X: x, Y: y}, topLeft, topRight, bottomRight, bottomLeft) {
				algorithms.SetPixel(canvas, x, y, p.Color) 
			}
		}
	}
	
	
	drawFilledSemicircle(canvas, p.Start.X, p.Start.Y, p.Radius, dirX, dirY, p.Color, antiAliasing)
	
	
	drawFilledSemicircle(canvas, p.End.X, p.End.Y, p.Radius, -dirX, -dirY, p.Color, antiAliasing)
}


func isPointInRect(p, a, b, c, d Point) bool {
	
	ab := crossProduct(b.X-a.X, b.Y-a.Y, p.X-a.X, p.Y-a.Y)
	bc := crossProduct(c.X-b.X, c.Y-b.Y, p.X-b.X, p.Y-b.Y)
	cd := crossProduct(d.X-c.X, d.Y-c.Y, p.X-c.X, p.Y-c.Y)
	da := crossProduct(a.X-d.X, a.Y-d.Y, p.X-d.X, p.Y-d.Y)
	
	
	return (ab >= 0 && bc >= 0 && cd >= 0 && da >= 0) || (ab <= 0 && bc <= 0 && cd <= 0 && da <= 0)
}


func crossProduct(x1, y1, x2, y2 int) int {
	return x1*y2 - y1*x2
}


func drawFilledSemicircle(canvas [][]color.Color, centerX, centerY, radius int, dirX, dirY float64, c color.Color, antiAliasing bool) {
	
	for y := centerY - radius; y <= centerY + radius; y++ {
		for x := centerX - radius; x <= centerX + radius; x++ {
			
			if y < 0 || x < 0 || y >= len(canvas) || x >= len(canvas[0]) {
				continue
			}
			
			
			vx := float64(x - centerX)
			vy := float64(y - centerY)
			
			
			distSq := vx*vx + vy*vy
			if distSq > float64(radius*radius) {
				continue
			}
			
			
			dotProduct := vx*dirX + vy*dirY
			
			
			if dotProduct >= 0 {
				algorithms.SetPixel(canvas, x, y, c) 
			}
		}
	}
	
	
	if antiAliasing {
		drawXiaolinWuCircle(canvas, centerX, centerY, radius, c)
	} else {
		drawMidpointCircle(canvas, centerX, centerY, radius, c)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}


func (p *Pill) Contains(point Point) bool {
	
	distSq1 := (point.X-p.Start.X)*(point.X-p.Start.X) + (point.Y-p.Start.Y)*(point.Y-p.Start.Y)
	if math.Abs(float64(distSq1)-float64(p.Radius*p.Radius)) <= float64(5*5) {
		return true
	}
	
	
	distSq2 := (point.X-p.End.X)*(point.X-p.End.X) + (point.Y-p.End.Y)*(point.Y-p.End.Y)
	if math.Abs(float64(distSq2)-float64(p.Radius*p.Radius)) <= float64(5*5) {
		return true
	}
	
	
	dx := p.End.X - p.Start.X
	dy := p.End.Y - p.Start.Y
	length := math.Sqrt(float64(dx*dx + dy*dy))
	
	if length > 0 {
		
		dirX := float64(dx) / length
		dirY := float64(dy) / length
		
		
		perpX := -dirY
		perpY := dirX
		
		
		vx := float64(point.X - p.Start.X)
		vy := float64(point.Y - p.Start.Y)
		
		projDir := vx*dirX + vy*dirY
		
		if projDir >= 0 && projDir <= length {
			
			projPerp := math.Abs(vx*perpX + vy*perpY)
			
			return projPerp <= float64(p.Radius+5)
		}
	}
	
	return false
}


func (p *Pill) GetControlPoints() []Point {
	return []Point{
		p.Start,
		p.End,
		{X: p.Start.X + p.Radius, Y: p.Start.Y}, 
	}
}


func (p *Pill) Move(deltaX, deltaY int) {
	p.Start.X += deltaX
	p.Start.Y += deltaY
	p.End.X += deltaX
	p.End.Y += deltaY
}


func (p *Pill) SetColor(c color.Color) {
	p.Color = c
}


func (p *Pill) GetColor() color.Color {
	return p.Color
}


func (p *Pill) Serialize() map[string]interface{} {
	r, g, b, a := p.Color.RGBA()
	return map[string]interface{}{
		"type":    "pill",
		"startX":  p.Start.X,
		"startY":  p.Start.Y,
		"endX":    p.End.X, 
		"endY":    p.End.Y,
		"radius":  p.Radius,
		"color":   []uint32{r, g, b, a},
	}
}


func (p *Pill) Clone() Shape {
	return &Pill{
		Start:  Point{X: p.Start.X, Y: p.Start.Y},
		End:    Point{X: p.End.X, Y: p.End.Y},
		Radius: p.Radius,
		Color:  p.Color,
		Step:   p.Step,
	}
}