package ui

import (
	"paint-drawer-pro/algorithms"
	"paint-drawer-pro/models"
)

// ClipPolygon performs polygon clipping using the Sutherland-Hodgman algorithm
// It handles the conversion between model.Point and algorithm.Point
func ClipPolygon(subject, clip []models.Point) []models.Point {
	// Convert model points to algorithm points
	subjectPoints := make([]algorithms.Point, len(subject))
	for i, p := range subject {
		subjectPoints[i] = algorithms.Point{X: p.X, Y: p.Y}
	}
	
	clipPoints := make([]algorithms.Point, len(clip))
	for i, p := range clip {
		clipPoints[i] = algorithms.Point{X: p.X, Y: p.Y}
	}
	
	// Perform clipping
	algClippedPoints := algorithms.SutherlandHodgman(subjectPoints, clipPoints)
	
	// Convert back to model points
	clippedVertices := make([]models.Point, len(algClippedPoints))
	for i, p := range algClippedPoints {
		clippedVertices[i] = models.Point{X: p.X, Y: p.Y}
	}
	
	return clippedVertices
}
