package models

// ResizePointType represents which corner of a rectangle is being resized
type ResizePointType int

const (
	None ResizePointType = iota
	TopLeft
	TopRight
	BottomRight
	BottomLeft
)
