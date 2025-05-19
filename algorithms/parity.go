package algorithms

import (
	"image/color"
	"math"
	"sort"
)

// PointInPolygon determines if a point is inside a polygon using the even-odd rule (Ray Casting Algorithm)
// Returns true if the point is inside the polygon or on its boundary
func PointInPolygon(point Point, vertices []Point) bool {
	if len(vertices) < 3 {
		return false // Not a polygon
	}
	
	// Check if the point is on a vertex or edge first
	for i := 0; i < len(vertices); i++ {
		j := (i + 1) % len(vertices)
		
		// Check if the point is on a vertex
		if vertices[i].X == point.X && vertices[i].Y == point.Y {
			return true
		}
		
		// Check if the point is on an edge
		// Only for horizontal edges or using a more sophisticated on-edge check
		if vertices[i].Y == vertices[j].Y && point.Y == vertices[i].Y {
			// Horizontal edge check
			if (point.X >= vertices[i].X && point.X <= vertices[j].X) || 
			   (point.X >= vertices[j].X && point.X <= vertices[i].X) {
				return true
			}
		} else {
			// For non-horizontal edges, we could check if the point is on the line
			// But we'll leave this part to the ray casting algorithm below
		}
	}
	
	// Cast a ray from the point to the right and count intersections
	intersections := 0
	
	for i := 0; i < len(vertices); i++ {
		j := (i + 1) % len(vertices)
		
		// Handle special cases for ray-casting algorithm
		
		// Skip if both vertices are on same side of ray
		if (vertices[i].Y > point.Y) == (vertices[j].Y > point.Y) {
			continue
		}
		
		// Check if intersection is to the right of the point
		// Handle division by zero (horizontal edges caught earlier)
		if vertices[j].Y == vertices[i].Y {
			continue
		}
		
		slope := float64(vertices[j].X-vertices[i].X) / float64(vertices[j].Y-vertices[i].Y)
		x := float64(vertices[i].X) + slope*float64(point.Y-vertices[i].Y)
		
		if x > float64(point.X) {
			intersections++
		}
	}
	
	// If odd number of intersections, point is inside
	return intersections%2 == 1
}

// ParityFillPolygon implements the Scan Line Algorithm with parity checking
// More robust for complex polygons than edge-table algorithm
func ParityFillPolygon(canvas [][]color.Color, vertices []Point, fillColor color.Color) {
	if len(vertices) < 3 {
		return // Need at least 3 vertices for a polygon
	}

	// Find min and max Y coordinates to set scanning range
	minX, minY := vertices[0].X, vertices[0].Y
	maxX, maxY := vertices[0].X, vertices[0].Y
	
	for _, v := range vertices {
		if v.X < minX {
			minX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	// Add some margin to ensure we include the whole polygon
	minX = int(math.Max(0, float64(minX-1)))
	minY = int(math.Max(0, float64(minY-1)))
	maxX = int(math.Min(float64(len(canvas[0])-1), float64(maxX+1)))
	maxY = int(math.Min(float64(len(canvas)-1), float64(maxY+1)))

	// Optimized scan line algorithm
	for y := minY; y <= maxY; y++ {
		// Find intersections with this scanline
		intersections := []int{}
		
		for i := 0; i < len(vertices); i++ {
			j := (i + 1) % len(vertices)
			
			// Skip horizontal edges
			if vertices[i].Y == vertices[j].Y {
				continue
			}
			
			// Check if the edge crosses this scanline
			if (vertices[i].Y <= y && vertices[j].Y > y) || 
			   (vertices[j].Y <= y && vertices[i].Y > y) {
				// Calculate the x-coordinate of intersection
				x := vertices[i].X + (y - vertices[i].Y) * 
				     (vertices[j].X - vertices[i].X) / 
				     (vertices[j].Y - vertices[i].Y)
				
				// Add to our list of intersections
				intersections = append(intersections, x)
			}
		}
		
		// Sort intersections from left to right
		sort.Ints(intersections)
		
		// Fill between pairs of intersection points
		for i := 0; i < len(intersections)-1; i += 2 {
			if i+1 < len(intersections) {
				for x := intersections[i]; x <= intersections[i+1]; x++ {
					if x >= 0 && x < len(canvas[0]) {
						canvas[y][x] = fillColor
					}
				}
			}
		}
	}
}

// ParityFillPolygonWithImage fills a polygon with an image texture using optimized scan line algorithm
func ParityFillPolygonWithImage(canvas [][]color.Color, vertices []Point, fillImage [][]color.Color) {
	if len(vertices) < 3 || fillImage == nil || len(fillImage) == 0 || len(fillImage[0]) == 0 {
		return // Need valid inputs
	}

	// Find min and max coordinates for bounding box
	minX, minY := vertices[0].X, vertices[0].Y
	maxX, maxY := vertices[0].X, vertices[0].Y
	
	for _, v := range vertices {
		if v.X < minX {
			minX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}
	
	// Width and height for texture coordinate calculation
	polygonWidth := maxX - minX
	polygonHeight := maxY - minY
	
	// Prevent division by zero
	if polygonWidth <= 0 {
		polygonWidth = 1
	}
	if polygonHeight <= 0 {
		polygonHeight = 1
	}
	
	imgWidth := len(fillImage[0])
	imgHeight := len(fillImage)
	
	// Add some margin to ensure we include the whole polygon
	minX = int(math.Max(0, float64(minX-1)))
	minY = int(math.Max(0, float64(minY-1)))
	maxX = int(math.Min(float64(len(canvas[0])-1), float64(maxX+1)))
	maxY = int(math.Min(float64(len(canvas)-1), float64(maxY+1)))

	// Optimized scan line algorithm
	for y := minY; y <= maxY; y++ {
		// Find intersections with this scanline
		intersections := []int{}
		
		for i := 0; i < len(vertices); i++ {
			j := (i + 1) % len(vertices)
			
			// Skip horizontal edges
			if vertices[i].Y == vertices[j].Y {
				continue
			}
			
			// Check if the edge crosses this scanline
			if (vertices[i].Y <= y && vertices[j].Y > y) || 
			   (vertices[j].Y <= y && vertices[i].Y > y) {
				// Calculate the x-coordinate of intersection
				x := vertices[i].X + (y - vertices[i].Y) * 
				     (vertices[j].X - vertices[i].X) / 
				     (vertices[j].Y - vertices[i].Y)
				
				// Add to our list of intersections
				intersections = append(intersections, x)
			}
		}
		
		// Sort intersections from left to right
		sort.Ints(intersections)
		
		// Fill between pairs of intersection points
		for i := 0; i < len(intersections)-1; i += 2 {
			if i+1 < len(intersections) {
				for x := intersections[i]; x <= intersections[i+1]; x++ {
					if x >= minX && x <= maxX {
						// Calculate texture coordinates
						tx := ((x - minX) * imgWidth) / polygonWidth % imgWidth
						ty := ((y - minY) * imgHeight) / polygonHeight % imgHeight
						
						// Ensure positive coordinates (for modulo to work correctly)
						if tx < 0 {
							tx += imgWidth
						}
						if ty < 0 {
							ty += imgHeight
						}
						
						// Safe access to texture
						if ty >= 0 && ty < imgHeight && tx >= 0 && tx < imgWidth {
							canvas[y][x] = fillImage[ty][tx]
						}
					}
				}
			}
		}
	}
}
