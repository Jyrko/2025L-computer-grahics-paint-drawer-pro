package algorithms

// SimplifyPolygon removes duplicate or near-duplicate vertices from a polygon
// threshold determines how close points need to be to be considered duplicates
func SimplifyPolygon(points []Point, threshold float64) []Point {
	if len(points) < 3 {
		return points
	}

	// Initialize result with the first point
	result := []Point{points[0]}
	
	// Add points only if they are far enough from the previous point
	for i := 1; i < len(points); i++ {
		prev := result[len(result)-1]
		current := points[i]
		
		// Calculate squared distance
		dx := float64(current.X - prev.X)
		dy := float64(current.Y - prev.Y)
		distSquared := dx*dx + dy*dy
		
		if distSquared > threshold*threshold {
			result = append(result, current)
		}
	}
	
	// Check if the last point is too close to the first point
	if len(result) > 2 {
		last := result[len(result)-1]
		first := result[0]
		
		dx := float64(last.X - first.X)
		dy := float64(last.Y - first.Y)
		distSquared := dx*dx + dy*dy
		
		if distSquared < threshold*threshold {
			// Remove the last point if it's too close to the first
			result = result[:len(result)-1]
		}
	}
	
	// Ensure we have at least 3 points for a valid polygon
	if len(result) < 3 {
		return points // Return original if simplification fails
	}
	
	return result
}

// IsPolygonSimple checks if the polygon has no self-intersections
func IsPolygonSimple(vertices []Point) bool {
	n := len(vertices)
	if n < 3 {
		return false
	}
	
	// Check each pair of non-adjacent edges for intersections
	for i := 0; i < n; i++ {
		i1 := (i + 1) % n
		
		for j := i + 2; j < n + i - 1; j++ {
			j1 := (j + 1) % n
			
			// Skip if the edges share a vertex
			if i1 == j || i == j1 {
				continue
			}
			
			// Check if the edges intersect
			if doLinesIntersect(vertices[i], vertices[i1], vertices[j % n], vertices[j1 % n]) {
				return false
			}
		}
	}
	
	return true
}

// doLinesIntersect checks if two line segments intersect
func doLinesIntersect(p1, q1, p2, q2 Point) bool {
	// Calculate orientations
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)
	
	// General case
	if o1 != o2 && o3 != o4 {
		return true
	}
	
	// Special cases for collinear points
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}
	
	return false
}

// orientation calculates the orientation of triplet (p, q, r)
// 0 = collinear, 1 = clockwise, 2 = counterclockwise
func orientation(p, q, r Point) int {
	val := (q.Y - p.Y) * (r.X - q.X) - (q.X - p.X) * (r.Y - q.Y)
	if val == 0 {
		return 0 // collinear
	}
	if val > 0 {
		return 1 // clockwise
	}
	return 2 // counter-clockwise
}

// onSegment checks if point q lies on segment pr
func onSegment(p, q, r Point) bool {
	return q.X <= max(p.X, r.X) && q.X >= min(p.X, r.X) &&
		   q.Y <= max(p.Y, r.Y) && q.Y >= min(p.Y, r.Y)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
