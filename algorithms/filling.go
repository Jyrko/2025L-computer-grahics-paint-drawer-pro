package algorithms

import (
	"image/color"
	"math"
	"sort"
)

type Edge struct {
	YMax   int     
	XOfYMin int     
	SlopeInv float64 
}

func EdgeTableFill(canvas [][]color.Color, vertices []Point, fillColor color.Color) {
	if len(vertices) < 3 {
		return 
	}
	
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

	edgeTable := make(map[int][]Edge)
	
	for i := 0; i < len(vertices); i++ {
		
		v1 := vertices[i]
		v2 := vertices[(i+1)%len(vertices)]
		
		if v1.Y == v2.Y {
			continue
		}
		
		
		if v1.Y > v2.Y {
			v1, v2 = v2, v1
		}
		
		
		slopeInv := float64(v2.X-v1.X) / float64(v2.Y-v1.Y)
		
		
		edge := Edge{
			YMax:    v2.Y,
			XOfYMin: v1.X,
			SlopeInv: slopeInv,
		}
		
		
		edgeTable[v1.Y] = append(edgeTable[v1.Y], edge)
	}
	
	
	var activeEdgeList []Edge
	
	
	for y := minY; y <= maxY; y++ {
		
		if edges, exists := edgeTable[y]; exists {
			activeEdgeList = append(activeEdgeList, edges...)
		}
		
		
		newAEL := make([]Edge, 0, len(activeEdgeList))
		for _, edge := range activeEdgeList {
			if edge.YMax > y {
				newAEL = append(newAEL, edge)
			}
		}
		activeEdgeList = newAEL
		
		
		sort.Slice(activeEdgeList, func(i, j int) bool {
			return activeEdgeList[i].XOfYMin < activeEdgeList[j].XOfYMin
		})
		
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
		
		
		for i := range activeEdgeList {
			activeEdgeList[i].XOfYMin += int(math.Round(activeEdgeList[i].SlopeInv))
		}
	}
}


func FillPolygonWithImage(canvas [][]color.Color, vertices []Point, fillImage [][]color.Color) {
	if len(vertices) < 3 || fillImage == nil || len(fillImage) == 0 || len(fillImage[0]) == 0 {
		return 
	}

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
	
	
	polygonWidth := maxX - minX
	polygonHeight := maxY - minY
	
	
	if polygonWidth <= 0 {
		polygonWidth = 1
	}
	if polygonHeight <= 0 {
		polygonHeight = 1
	}
	
	imgWidth := len(fillImage[0])
	imgHeight := len(fillImage)

	
	edgeTable := make(map[int][]Edge)
	
	
	for i := 0; i < len(vertices); i++ {
		
		v1 := vertices[i]
		v2 := vertices[(i+1)%len(vertices)]
		
		
		if v1.Y == v2.Y {
			continue
		}
		
		
		if v1.Y > v2.Y {
			v1, v2 = v2, v1
		}
		
		
		slopeInv := float64(v2.X-v1.X) / float64(v2.Y-v1.Y)
		
		
		edge := Edge{
			YMax:    v2.Y,
			XOfYMin: v1.X,
			SlopeInv: slopeInv,
		}
		
		
		edgeTable[v1.Y] = append(edgeTable[v1.Y], edge)
	}
	
	
	var activeEdgeList []Edge
	
	
	for y := minY; y <= maxY; y++ {
		
		if edges, exists := edgeTable[y]; exists {
			activeEdgeList = append(activeEdgeList, edges...)
		}
		
		
		newAEL := make([]Edge, 0, len(activeEdgeList))
		for _, edge := range activeEdgeList {
			if edge.YMax > y {
				newAEL = append(newAEL, edge)
			}
		}
		activeEdgeList = newAEL
		
		
		sort.Slice(activeEdgeList, func(i, j int) bool {
			return activeEdgeList[i].XOfYMin < activeEdgeList[j].XOfYMin
		})
		
		
		for i := 0; i < len(activeEdgeList)-1; i += 2 {
			if i+1 < len(activeEdgeList) {
				xStart := int(math.Floor(float64(activeEdgeList[i].XOfYMin)))
				xEnd := int(math.Ceil(float64(activeEdgeList[i+1].XOfYMin)))
				
				for x := xStart; x <= xEnd; x++ {
					if y >= 0 && y < len(canvas) && x >= 0 && x < len(canvas[0]) {
						
						
						tx := ((x - minX) * imgWidth) / polygonWidth % imgWidth
						ty := ((y - minY) * imgHeight) / polygonHeight % imgHeight
						
						
						if tx < 0 {
							tx += imgWidth
						}
						if ty < 0 {
							ty += imgHeight
						}
						
						
						if ty >= 0 && ty < imgHeight && tx >= 0 && tx < imgWidth {
							canvas[y][x] = fillImage[ty][tx]
						}
					}
				}
			}
		}
		
		
		for i := range activeEdgeList {
			activeEdgeList[i].XOfYMin += int(math.Round(activeEdgeList[i].SlopeInv))
		}
	}
}


type ScanlineSegment struct {
	Y, XLeft, XRight int
}


func SmithScanlineFill(canvas [][]color.Color, startPoint Point, fillColor color.Color, boundaryColor color.Color) {
	if startPoint.Y < 0 || startPoint.Y >= len(canvas) || startPoint.X < 0 || startPoint.X >= len(canvas[0]) {
		return 
	}

	if canvas[startPoint.Y][startPoint.X] == boundaryColor || canvas[startPoint.Y][startPoint.X] == fillColor {
		return 
	}

	stack := []ScanlineSegment{}
	canvasHeight := len(canvas)
	canvasWidth := len(canvas[0])

	
	xLeft, xRight := startPoint.X, startPoint.X

	
	for xLeft > 0 && canvas[startPoint.Y][xLeft-1] != boundaryColor {
		xLeft--
	}
	
	for xRight < canvasWidth-1 && canvas[startPoint.Y][xRight+1] != boundaryColor {
		xRight++
	}

	stack = append(stack, ScanlineSegment{Y: startPoint.Y, XLeft: xLeft, XRight: xRight})

	for len(stack) > 0 {
		
		segment := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		y, currentXLeft, currentXRight := segment.Y, segment.XLeft, segment.XRight

		
		for x := currentXLeft; x <= currentXRight; x++ {
			if y >= 0 && y < canvasHeight && x >= 0 && x < canvasWidth {
				canvas[y][x] = fillColor
			}
		}

		if y-1 >= 0 {
			processScanline(canvas, y-1, currentXLeft, currentXRight, fillColor, boundaryColor, &stack, canvasWidth)
		}
		
		if y+1 < canvasHeight {
			processScanline(canvas, y+1, currentXLeft, currentXRight, fillColor, boundaryColor, &stack, canvasWidth)
		}
	}
}

func processScanline(canvas [][]color.Color, y int, parentXLeft int, parentXRight int, fillColor color.Color, boundaryColor color.Color, stack *[]ScanlineSegment, canvasWidth int) {
	x := parentXLeft
	for x <= parentXRight {
		for x <= parentXRight && (canvas[y][x] == boundaryColor || canvas[y][x] == fillColor) {
			x++
		}
		if x > parentXRight {
			break
		}

		segmentXLeft := x

		for x <= parentXRight && canvas[y][x] != boundaryColor && canvas[y][x] != fillColor {
			x++
		}
		segmentXRight := x - 1

		*stack = append(*stack, ScanlineSegment{Y: y, XLeft: segmentXLeft, XRight: segmentXRight})

		if segmentXRight >= parentXRight {
			x = segmentXRight + 1
			for x < canvasWidth {
				for x < canvasWidth && (canvas[y][x] == boundaryColor || canvas[y][x] == fillColor) {
					x++
				}
				if x >= canvasWidth {
					break
				}
				newScanXLeft := x
				for x < canvasWidth && canvas[y][x] != boundaryColor && canvas[y][x] != fillColor {
					x++
				}
				newScanXRight := x - 1
				*stack = append(*stack, ScanlineSegment{Y: y, XLeft: newScanXLeft, XRight: newScanXRight})
			}

		}
	}
}
