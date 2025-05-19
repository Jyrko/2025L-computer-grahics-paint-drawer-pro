package algorithms




func PointAdapter(modelPoints interface{}) []Point {
	var algorithmPoints []Point

	switch points := modelPoints.(type) {
	case []Point:
		return points 

	case []struct{ X, Y int }:
		algorithmPoints = make([]Point, len(points))
		for i, p := range points {
			algorithmPoints[i] = Point{X: p.X, Y: p.Y}
		}
		
	case []interface{}:
		algorithmPoints = make([]Point, len(points))
		for i, p := range points {
			if point, ok := p.(struct{ X, Y int }); ok {
				algorithmPoints[i] = Point{X: point.X, Y: point.Y}
			} else if pointMap, ok := p.(map[string]int); ok {
				algorithmPoints[i] = Point{X: pointMap["X"], Y: pointMap["Y"]}
			}
		}
	}

	return algorithmPoints
}
