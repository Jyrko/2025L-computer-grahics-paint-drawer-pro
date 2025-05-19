package algorithms

import (
	"math"
)


func SutherlandHodgman(subject, clip []Point) []Point {
	if len(subject) < 3 || len(clip) < 3 {
		return nil 
	}
	
	output := subject

	for i := 0; i < len(clip); i++ {
		clipEdgeStart := clip[i]
		clipEdgeEnd := clip[(i+1)%len(clip)]
		
		input := output
		output = []Point{}
		
		if len(input) == 0 {
			break
		}

		
		s := input[len(input)-1]

		
		for _, e := range input {
			
			if isInside(clipEdgeStart, clipEdgeEnd, e) {
				
				if !isInside(clipEdgeStart, clipEdgeEnd, s) {
					intersection := computeIntersection(s, e, clipEdgeStart, clipEdgeEnd)
					output = append(output, intersection)
				}
				
				output = append(output, e)
			} else if isInside(clipEdgeStart, clipEdgeEnd, s) {
				
				
				intersection := computeIntersection(s, e, clipEdgeStart, clipEdgeEnd)
				output = append(output, intersection)
			}
			
			s = e
		}
	}

	return output
}



func isInside(clipEdgeStart, clipEdgeEnd, point Point) bool {
	return (clipEdgeEnd.X - clipEdgeStart.X) * (point.Y - clipEdgeStart.Y) - 
	       (clipEdgeEnd.Y - clipEdgeStart.Y) * (point.X - clipEdgeStart.X) <= 0
}


func computeIntersection(s, e, clipEdgeStart, clipEdgeEnd Point) Point {
	
	x1, y1 := float64(s.X), float64(s.Y)
	x2, y2 := float64(e.X), float64(e.Y)
	x3, y3 := float64(clipEdgeStart.X), float64(clipEdgeStart.Y)
	x4, y4 := float64(clipEdgeEnd.X), float64(clipEdgeEnd.Y)

	
	
	
	
	denom := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)
	
	
	if math.Abs(denom) < 0.0001 {
		
		return Point{
			X: int((x1 + x2) / 2),
			Y: int((y1 + y2) / 2),
		}
	}
	
	
	ua := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / denom
	
	
	intersectX := x1 + ua*(x2-x1)
	intersectY := y1 + ua*(y2-y1)
	
	return Point{
		X: int(math.Round(intersectX)),
		Y: int(math.Round(intersectY)),
	}
}



func IsPolygonConvex(verticesInput interface{}) bool {
	
	vertices := PointAdapter(verticesInput)
	
	
	length := len(vertices)
	
	if length < 3 {
		return false 
	}
	
	
	vertices = SimplifyPolygon(vertices, 2.0) 
	
	
	length = len(vertices)
	
	
	if length == 3 {
		
		x1, y1 := vertices[0].X, vertices[0].Y
		x2, y2 := vertices[1].X, vertices[1].Y
		x3, y3 := vertices[2].X, vertices[2].Y
		
		
		area := (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2)) / 2
		return area != 0 
	}
	
	
	sign := 0
	
	for i := 0; i < length; i++ {
		j := (i + 1) % length
		k := (i + 2) % length
		
		
		xi := vertices[i].X
		yi := vertices[i].Y
		xj := vertices[j].X
		yj := vertices[j].Y
		xk := vertices[k].X
		yk := vertices[k].Y
		
		
		dx1 := xj - xi
		dy1 := yj - yi
		dx2 := xk - xj
		dy2 := yk - yj
		
		
		cross := dx1*dy2 - dy1*dx2
		
		if cross == 0 {
			continue 
		}
		
		
		currentSign := 0
		if cross > 0 {
			currentSign = 1
		} else {
			currentSign = -1
		}
		
		if sign == 0 {
			sign = currentSign
		} else if sign * currentSign < 0 {
			return false 
		}
	}
	
	return true
}
