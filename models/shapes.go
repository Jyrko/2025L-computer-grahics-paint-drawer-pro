package models

import (
	"image/color"
)


type Point struct {
	X, Y int
}


type Shape interface {
	Draw(canvas [][]color.Color, antiAliasing bool)
	Contains(p Point) bool
	Move(deltaX, deltaY int)
	GetControlPoints() []Point
	SetColor(c color.Color)
	GetColor() color.Color
	Serialize() map[string]interface{}
	Clone() Shape
}


type Circle struct {
	Center Point
	Radius int
	Color  color.Color
}


type Polygon struct {
	Vertices  []Point
	Color     color.Color
	Thickness int
}


type DrawingState struct {
	Shapes         []Shape
	SelectedShape  Shape
	CurrentShape   Shape  
	CurrentAction  string
	AntiAliasing   bool
	PenType        string 
	BrushThickness int    
}