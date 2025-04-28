package models

import (
	"image/color"
	"math"
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
	
	
	if length < float64(2*p.Radius) {
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

	
	if antiAliasing {
		drawXiaolinWuCircle(canvas, p.Start.X, p.Start.Y, p.Radius, p.Color)
	} else {
		drawMidpointCircle(canvas, p.Start.X, p.Start.Y, p.Radius, p.Color)
	}

	
	if antiAliasing {
		drawXiaolinWuCircle(canvas, p.End.X, p.End.Y, p.Radius, p.Color)
	} else {
		drawMidpointCircle(canvas, p.End.X, p.End.Y, p.Radius, p.Color)
	}

	
	
	x1 := int(float64(p.Start.X) + perpX*float64(p.Radius))
	y1 := int(float64(p.Start.Y) + perpY*float64(p.Radius))
	
	x2 := int(float64(p.Start.X) - perpX*float64(p.Radius))
	y2 := int(float64(p.Start.Y) - perpY*float64(p.Radius))
	
	x3 := int(float64(p.End.X) + perpX*float64(p.Radius))
	y3 := int(float64(p.End.Y) + perpY*float64(p.Radius))
	
	x4 := int(float64(p.End.X) - perpX*float64(p.Radius))
	y4 := int(float64(p.End.Y) - perpY*float64(p.Radius))
	
	
	if antiAliasing {
		drawXiaolinWuLine(canvas, x1, y1, x3, y3, p.Color)
		drawXiaolinWuLine(canvas, x2, y2, x4, y4, p.Color)
	} else {
		drawMidpointLine(canvas, x1, y1, x3, y3, p.Color)
		drawMidpointLine(canvas, x2, y2, x4, y4, p.Color)
	}
	
	
	
	
	steps := int(length) * 2
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps)
		pointX := int(float64(p.Start.X) + dirX*t*length)
		pointY := int(float64(p.Start.Y) + dirY*t*length)
		
		
		x1 := int(float64(pointX) + perpX*float64(p.Radius))
		y1 := int(float64(pointY) + perpY*float64(p.Radius))
		x2 := int(float64(pointX) - perpX*float64(p.Radius))
		y2 := int(float64(pointY) - perpY*float64(p.Radius))
		
		if antiAliasing {
			drawXiaolinWuLine(canvas, x1, y1, x2, y2, p.Color)
		} else {
			drawMidpointLine(canvas, x1, y1, x2, y2, p.Color)
		}
	}
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
		
		
		vx := point.X - p.Start.X
		vy := point.Y - p.Start.Y
		
		
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