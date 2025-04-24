package models

import (
	"image/color"
	"math"
)


func NewCircle(center Point, radius int, color color.Color) *Circle {
	if radius <= 0 {
		radius = 1
	}
	return &Circle{
		Center: center,
		Radius: radius,
		Color:  color,
	}
}


func (c *Circle) Draw(canvas [][]color.Color, antiAliasing bool) {
	if antiAliasing {
		drawXiaolinWuCircle(canvas, c.Center.X, c.Center.Y, c.Radius, c.Color)
	} else {
		drawMidpointCircle(canvas, c.Center.X, c.Center.Y, c.Radius, c.Color)
	}
}


func (c *Circle) Contains(p Point) bool {
	distSq := (p.X-c.Center.X)*(p.X-c.Center.X) + (p.Y-c.Center.Y)*(p.Y-c.Center.Y)
	
	radiusSq := c.Radius * c.Radius
	return math.Abs(float64(distSq-radiusSq)) <= float64(5*5) 
}


func (c *Circle) GetControlPoints() []Point {
	
	return []Point{
		c.Center,
		{X: c.Center.X + c.Radius, Y: c.Center.Y},
	}
}


func (c *Circle) Move(deltaX, deltaY int) {
	c.Center.X += deltaX
	c.Center.Y += deltaY
}


func (c *Circle) SetColor(color color.Color) {
	c.Color = color
}


func (c *Circle) GetColor() color.Color {
	return c.Color
}


func (c *Circle) Serialize() map[string]interface{} {
	r, g, b, a := c.Color.RGBA()
	return map[string]interface{}{
		"type":    "circle",
		"centerX": c.Center.X,
		"centerY": c.Center.Y,
		"radius":  c.Radius,
		"color":   []uint32{r, g, b, a},
	}
}


func (c *Circle) Clone() Shape {
	return NewCircle(
		Point{X: c.Center.X, Y: c.Center.Y},
		c.Radius,
		c.Color,
	)
}