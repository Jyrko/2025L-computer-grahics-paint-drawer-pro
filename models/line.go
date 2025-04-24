package models

import (
	"image/color"
	"math"
)


type Line struct {
	Start     Point
	End       Point
	Color     color.Color
	Thickness int
	PenType   string 
}


func NewLine(start, end Point, color color.Color, thickness int, penType string) *Line {
	if thickness <= 0 {
		thickness = 1
	}
	return &Line{
		Start:     start,
		End:       end,
		Color:     color,
		Thickness: thickness,
		PenType:   penType,
	}
}


func (l *Line) Draw(canvas [][]color.Color, antiAliasing bool) {
	if l.PenType == "regular" {

		if antiAliasing {
			drawXiaolinWuLine(canvas, l.Start.X, l.Start.Y, l.End.X, l.End.Y, l.Color)
		} else {
			drawMidpointLine(canvas, l.Start.X, l.Start.Y, l.End.X, l.End.Y, l.Color)
		}
	} else { 
		if antiAliasing {
			drawXiaolinWuLine(canvas, l.Start.X, l.Start.Y, l.End.X, l.End.Y, l.Color)
		} else {
			if l.Thickness == 1 {
				drawMidpointLine(canvas, l.Start.X, l.Start.Y, l.End.X, l.End.Y, l.Color)
			} else {
				drawThickLine(canvas, l.Start.X, l.Start.Y, l.End.X, l.End.Y, l.Color, l.Thickness)
			}
		}
	}
}


func (l *Line) Contains(p Point) bool {
	
	lineLen := math.Sqrt(float64((l.End.X-l.Start.X)*(l.End.X-l.Start.X) + (l.End.Y-l.Start.Y)*(l.End.Y-l.Start.Y)))
	if lineLen == 0 {

		return math.Abs(float64(p.X-l.Start.X)) <= 5 && math.Abs(float64(p.Y-l.Start.Y)) <= 5
	}
	
	
	t := float64((p.X-l.Start.X)*(l.End.X-l.Start.X) + (p.Y-l.Start.Y)*(l.End.Y-l.Start.Y)) / (lineLen * lineLen)
	t = math.Max(0, math.Min(1, t))
	
	nearestX := l.Start.X + int(float64(l.End.X-l.Start.X)*t)
	nearestY := l.Start.Y + int(float64(l.End.Y-l.Start.Y)*t)
	
	dist := math.Sqrt(float64((p.X-nearestX)*(p.X-nearestX) + (p.Y-nearestY)*(p.Y-nearestY)))
	return dist <= float64(l.Thickness+5) 
}


func (l *Line) GetControlPoints() []Point {
	return []Point{l.Start, l.End}
}


func (l *Line) Move(deltaX, deltaY int) {
	l.Start.X += deltaX
	l.Start.Y += deltaY
	l.End.X += deltaX
	l.End.Y += deltaY
}


func (l *Line) SetColor(c color.Color) {
	l.Color = c
}


func (l *Line) GetColor() color.Color {
	return l.Color
}


func (l *Line) Serialize() map[string]interface{} {
	r, g, b, a := l.Color.RGBA()
	return map[string]interface{}{
		"type":      "line",
		"startX":    l.Start.X,
		"startY":    l.Start.Y,
		"endX":      l.End.X,
		"endY":      l.End.Y,
		"color":     []uint32{r, g, b, a},
		"thickness": l.Thickness,
		"penType":   l.PenType,
	}
}


func (l *Line) Clone() Shape {
	return NewLine(
		Point{X: l.Start.X, Y: l.Start.Y},
		Point{X: l.End.X, Y: l.End.Y},
		l.Color,
		l.Thickness,
		l.PenType,
	)
}