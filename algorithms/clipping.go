package algorithms

import (
	"math"
)

// SutherlandHodgman implements the Sutherland-Hodgman polygon clipping algorithm
// Subject polygon is the polygon to be clipped (can be any polygon, concave or convex)
// Clip polygon is the polygon to clip against (must be convex)
// Returns the clipped polygon vertices
func SutherlandHodgman(subject, clip []Point) []Point {
	if len(subject) < 3 || len(clip) < 3 {
		return nil // Invalid input
	}

	// Output will start as the subject polygon
	output := subject

	// For each edge of the clip polygon
	for i := 0; i < len(clip); i++ {
		// Edge vertices
		clipEdgeStart := clip[i]
		clipEdgeEnd := clip[(i+1)%len(clip)]

		// Current output becomes input for next edge
		input := output
		output = []Point{}

		// No vertices, nothing to do
		if len(input) == 0 {
			break
		}

		// Last point in input
		s := input[len(input)-1]

		// For each edge in input polygon
		for _, e := range input {
			// If current point is inside the clip edge
			if isInside(clipEdgeStart, clipEdgeEnd, e) {
				// If previous point was outside, add intersection
				if !isInside(clipEdgeStart, clipEdgeEnd, s) {
					intersection := computeIntersection(s, e, clipEdgeStart, clipEdgeEnd)
					output = append(output, intersection)
				}
				// Add the current point
				output = append(output, e)
			} else if isInside(clipEdgeStart, clipEdgeEnd, s) {
				// Current point is outside but previous was inside
				// Add the intersection point
				intersection := computeIntersection(s, e, clipEdgeStart, clipEdgeEnd)
				output = append(output, intersection)
			}
			// Update s for next iteration
			s = e
		}
	}

	return output
}

// isInside determines if a point is inside or outside a clip edge
// A point is inside if it's on the right side of the directed edge
func isInside(clipEdgeStart, clipEdgeEnd, point Point) bool {
	return (clipEdgeEnd.X - clipEdgeStart.X) * (point.Y - clipEdgeStart.Y) - 
	       (clipEdgeEnd.Y - clipEdgeStart.Y) * (point.X - clipEdgeStart.X) <= 0
}

// computeIntersection calculates the intersection point between two lines
func computeIntersection(s, e, clipEdgeStart, clipEdgeEnd Point) Point {
	// Convert to float for precise calculation
	x1, y1 := float64(s.X), float64(s.Y)
	x2, y2 := float64(e.X), float64(e.Y)
	x3, y3 := float64(clipEdgeStart.X), float64(clipEdgeStart.Y)
	x4, y4 := float64(clipEdgeEnd.X), float64(clipEdgeEnd.Y)

	// Line 1 as parametric: P = s + t(e-s), 0 <= t <= 1
	// Line 2 as parametric: P = clipEdgeStart + u(clipEdgeEnd-clipEdgeStart), 0 <= u <= 1
	
	// Calculate denominator for intersection formulas
	denom := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
	
	// Prevent division by zero
	if math.Abs(denom) < 0.0001 {
		// Lines are parallel, return midpoint as fallback
		return Point{
			X: int((x1 + x2) / 2),
			Y: int((y1 + y2) / 2),
		}
	}
	
	// Calculate intersection parameters
	ua := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / denom
	
	// Calculate intersection point
	intersectX := x1 + ua*(x2-x1)
	intersectY := y1 + ua*(y2-y1)
	
	return Point{
		X: int(math.Round(intersectX)),
		Y: int(math.Round(intersectY)),
	}
}

// IsPolygonConvex checks if a polygon is convex
// A polygon is convex if all interior angles are less than 180 degrees
func IsPolygonConvex(verticesInput interface{}) bool {
	// Convert to our internal Point type
	vertices := PointAdapter(verticesInput)
	
	// Get length
	length := len(vertices)
	
	if length < 3 {
		return false // Not a polygon
	}
	
	// Remove duplicate or near-duplicate points
	vertices = SimplifyPolygon(vertices, 2.0) // 2.0 pixel threshold
	
	// Count actual number of vertices after simplification
	length = len(vertices)
	
	// All triangles are convex (unless they have collinear points)
	if length == 3 {
		// Check if points are not collinear
		x1, y1 := vertices[0].X, vertices[0].Y
		x2, y2 := vertices[1].X, vertices[1].Y
		x3, y3 := vertices[2].X, vertices[2].Y
		
		// Calculate the area of the triangle using cross product
		area := (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2)) / 2
		return area != 0 // If area is not zero, the triangle is convex
	}
	
	// For a polygon to be convex, the cross products of consecutive edges must have the same sign
	sign := 0
	
	for i := 0; i < length; i++ {
		j := (i + 1) % length
		k := (i + 2) % length
		
		// Get coordinates for each point
		xi := vertices[i].X
		yi := vertices[i].Y
		xj := vertices[j].X
		yj := vertices[j].Y
		xk := vertices[k].X
		yk := vertices[k].Y
		
		// Vectors for consecutive edges
		dx1 := xj - xi
		dy1 := yj - yi
		dx2 := xk - xj
		dy2 := yk - yj
		
		// Cross product (z component in 2D is scalar)
		cross := dx1*dy2 - dy1*dx2
		
		if cross == 0 {
			continue // Collinear points, skip
		}
		
		// Check sign consistency
		currentSign := 0
		if cross > 0 {
			currentSign = 1
		} else {
			currentSign = -1
		}
		
		if sign == 0 {
			sign = currentSign
		} else if sign * currentSign < 0 {
			return false // Found both positive and negative cross products
		}
	}
	
	return true
}
