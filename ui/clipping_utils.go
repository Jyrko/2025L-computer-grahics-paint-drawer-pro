package ui

import (
	"paint-drawer-pro/algorithms"
	"paint-drawer-pro/models"
)



func ClipPolygon(subject, clip []models.Point) []models.Point {
	
	subjectPoints := make([]algorithms.Point, len(subject))
	for i, p := range subject {
		subjectPoints[i] = algorithms.Point{X: p.X, Y: p.Y}
	}
	
	clipPoints := make([]algorithms.Point, len(clip))
	for i, p := range clip {
		clipPoints[i] = algorithms.Point{X: p.X, Y: p.Y}
	}
	
	
	subjectPoints = algorithms.SimplifyPolygon(subjectPoints, 2.0) 
	clipPoints = algorithms.SimplifyPolygon(clipPoints, 2.0)
	
	
	algClippedPoints := algorithms.SutherlandHodgman(subjectPoints, clipPoints)
	
	
	clippedVertices := make([]models.Point, len(algClippedPoints))
	for i, p := range algClippedPoints {
		clippedVertices[i] = models.Point{X: p.X, Y: p.Y}
	}
	
	return clippedVertices
}
