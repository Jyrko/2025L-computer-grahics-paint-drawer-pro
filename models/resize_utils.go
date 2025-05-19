package models


type ResizePointType int

const (
	None ResizePointType = iota
	TopLeft
	TopRight
	BottomRight
	BottomLeft
)
