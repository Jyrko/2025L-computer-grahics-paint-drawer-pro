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
	FillColor color.Color
	IsFilled  bool
	FillImage [][]color.Color
	UseImage  bool
}


type DrawingState struct {
	Shapes         []Shape
	SelectedShape  Shape
	CurrentShape   Shape  
	CurrentAction  string
	AntiAliasing   bool
	PenType        string 
	BrushThickness int    
	CurrentColor   color.RGBA
	FillEnabled    bool
	FillColor      color.Color
	FillImage      [][]color.Color
	UseImageFill   bool 
	SelectionRect  *Rectangle // Added for scanline fill area selection
	FillStage      string     // Added for scanline fill stages (e.g., selecting_area, area_selected)
}

// Normalize ensures that TopLeft is actually the top-left
// and BottomRight is the bottom-right.
func (r *Rectangle) Normalize() {
	if r.TopLeft.X > r.BottomRight.X {
		r.TopLeft.X, r.BottomRight.X = r.BottomRight.X, r.TopLeft.X
	}
	if r.TopLeft.Y > r.BottomRight.Y {
		r.TopLeft.Y, r.BottomRight.Y = r.BottomRight.Y, r.TopLeft.Y
	}
}

// Contains checks if a point is inside the rectangle.
// ...existing code...