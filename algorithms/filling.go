package algorithms

import (
	"image/color"
	"math"
	"sort"
)

// Edge represents an edge in the Edge Table algorithm
type Edge struct {
	YMax   int     // Maximum y-coordinate of the edge
	XOfYMin int     // X-coordinate at the minimum y-coordinate
	SlopeInv float64 // 1/slope (dx/dy)
}

// EdgeTableFill implements the Edge Table polygon filling algorithm
func EdgeTableFill(canvas [][]color.Color, vertices []Point, fillColor color.Color) {
	if len(vertices) < 3 {
		return // Need at least 3 vertices for a polygon
	}

	// Find min and max Y coordinates to set scanning range
	minY := vertices[0].Y
	maxY := vertices[0].Y
	for _, v := range vertices {
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	// Create edge table
	edgeTable := make(map[int][]Edge)
	
	// Process each edge
	for i := 0; i < len(vertices); i++ {
		// Get the current edge points
		v1 := vertices[i]
		v2 := vertices[(i+1)%len(vertices)]
		
		// Skip horizontal edges
		if v1.Y == v2.Y {
			continue
		}
		
		// Ensure v1.Y < v2.Y for consistent processing
		if v1.Y > v2.Y {
			v1, v2 = v2, v1
		}
		
		// Calculate inverse slope (dx/dy)
		slopeInv := float64(v2.X-v1.X) / float64(v2.Y-v1.Y)
		
		// Create edge and add to edge table at its yMin
		edge := Edge{
			YMax:    v2.Y,
			XOfYMin: v1.X,
			SlopeInv: slopeInv,
		}
		
		// Add edge to the edge table, keyed by yMin
		edgeTable[v1.Y] = append(edgeTable[v1.Y], edge)
	}
	
	// Active edge list (initially empty)
	var activeEdgeList []Edge
	
	// Process each scanline from bottom to top
	for y := minY; y <= maxY; y++ {
		// Add edges starting at this scanline to AEL
		if edges, exists := edgeTable[y]; exists {
			activeEdgeList = append(activeEdgeList, edges...)
		}
		
		// Remove edges that end at this scanline
		newAEL := make([]Edge, 0, len(activeEdgeList))
		for _, edge := range activeEdgeList {
			if edge.YMax > y {
				newAEL = append(newAEL, edge)
			}
		}
		activeEdgeList = newAEL
		
		// Sort edges by x-coordinate
		sort.Slice(activeEdgeList, func(i, j int) bool {
			return activeEdgeList[i].XOfYMin < activeEdgeList[j].XOfYMin
		})
		
		// Fill between pairs of edges
		for i := 0; i < len(activeEdgeList)-1; i += 2 {
			if i+1 < len(activeEdgeList) {
				xStart := int(math.Floor(float64(activeEdgeList[i].XOfYMin)))
				xEnd := int(math.Ceil(float64(activeEdgeList[i+1].XOfYMin)))
				
				for x := xStart; x <= xEnd; x++ {
					if y >= 0 && y < len(canvas) && x >= 0 && x < len(canvas[0]) {
						canvas[y][x] = fillColor
					}
				}
			}
		}
		
		// Update x-coordinates for next scanline
		for i := range activeEdgeList {
			activeEdgeList[i].XOfYMin += int(math.Round(activeEdgeList[i].SlopeInv))
		}
	}
}

// FillPolygonWithImage fills a polygon with an image texture using the Edge Table algorithm
func FillPolygonWithImage(canvas [][]color.Color, vertices []Point, fillImage [][]color.Color) {
	if len(vertices) < 3 || fillImage == nil || len(fillImage) == 0 || len(fillImage[0]) == 0 {
		return // Need valid inputs
	}

	// Find min and max coordinates to set scanning range and for texture mapping
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

	// Create edge table
	edgeTable := make(map[int][]Edge)
	
	// Process each edge
	for i := 0; i < len(vertices); i++ {
		// Get the current edge points
		v1 := vertices[i]
		v2 := vertices[(i+1)%len(vertices)]
		
		// Skip horizontal edges
		if v1.Y == v2.Y {
			continue
		}
		
		// Ensure v1.Y < v2.Y for consistent processing
		if v1.Y > v2.Y {
			v1, v2 = v2, v1
		}
		
		// Calculate inverse slope (dx/dy)
		slopeInv := float64(v2.X-v1.X) / float64(v2.Y-v1.Y)
		
		// Create edge and add to edge table at its yMin
		edge := Edge{
			YMax:    v2.Y,
			XOfYMin: v1.X,
			SlopeInv: slopeInv,
		}
		
		// Add edge to the edge table, keyed by yMin
		edgeTable[v1.Y] = append(edgeTable[v1.Y], edge)
	}
	
	// Active edge list (initially empty)
	var activeEdgeList []Edge
	
	// Process each scanline from bottom to top
	for y := minY; y <= maxY; y++ {
		// Add edges starting at this scanline to AEL
		if edges, exists := edgeTable[y]; exists {
			activeEdgeList = append(activeEdgeList, edges...)
		}
		
		// Remove edges that end at this scanline
		newAEL := make([]Edge, 0, len(activeEdgeList))
		for _, edge := range activeEdgeList {
			if edge.YMax > y {
				newAEL = append(newAEL, edge)
			}
		}
		activeEdgeList = newAEL
		
		// Sort edges by x-coordinate
		sort.Slice(activeEdgeList, func(i, j int) bool {
			return activeEdgeList[i].XOfYMin < activeEdgeList[j].XOfYMin
		})
		
		// Fill between pairs of edges
		for i := 0; i < len(activeEdgeList)-1; i += 2 {
			if i+1 < len(activeEdgeList) {
				xStart := int(math.Floor(float64(activeEdgeList[i].XOfYMin)))
				xEnd := int(math.Ceil(float64(activeEdgeList[i+1].XOfYMin)))
				
				for x := xStart; x <= xEnd; x++ {
					if y >= 0 && y < len(canvas) && x >= 0 && x < len(canvas[0]) {
						// Calculate texture coordinates
						// Map from polygon space to texture space
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
		
		// Update x-coordinates for next scanline
		for i := range activeEdgeList {
			activeEdgeList[i].XOfYMin += int(math.Round(activeEdgeList[i].SlopeInv))
		}
	}
}
