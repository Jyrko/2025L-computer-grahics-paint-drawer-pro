package algorithms

import (
	"image/color"
	"math"
)


type Point struct {
	X, Y int
}


func SetPixel(canvas [][]color.Color, x, y int, c color.Color) {
	if x >= 0 && y >= 0 && y < len(canvas) && x < len(canvas[y]) {
		canvas[y][x] = c
	}
}


func SetPixelWithAlpha(canvas [][]color.Color, x, y int, c color.Color, alpha float64) {
	if x < 0 || y < 0 || y >= len(canvas) || x >= len(canvas[y]) {
		return
	}
	
	r1, g1, b1, a1 := c.RGBA()
	r1, g1, b1 = r1>>8, g1>>8, b1>>8
	a1 = uint32(float64(a1>>8) * alpha)
	
	if a1 == 0 {
		return 
	}
	
	if a1 == 255 {
		canvas[y][x] = c
		return
	}
	
	
	r2, g2, b2, a2 := canvas[y][x].RGBA()
	r2, g2, b2, a2 = r2>>8, g2>>8, b2>>8, a2>>8
	
	
	a2f := float64(a2) / 255.0
	af := float64(a1) / 255.0
	outAlpha := af + a2f*(1-af)
	
	var r, g, b uint8
	if outAlpha > 0 {
		r = uint8((float64(r1)*af + float64(r2)*a2f*(1-af)) / outAlpha)
		g = uint8((float64(g1)*af + float64(g2)*a2f*(1-af)) / outAlpha)
		b = uint8((float64(b1)*af + float64(b2)*a2f*(1-af)) / outAlpha)
	}
	
	canvas[y][x] = color.RGBA{R: r, G: g, B: b, A: uint8(outAlpha * 255)}
}



func MidpointLine(canvas [][]color.Color, x0, y0, x1, y1 int, c color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	
	
	if dx == 0 {
		startY, endY := y0, y1
		if y0 > y1 {
			startY, endY = y1, y0
		}
		for y := startY; y <= endY; y++ {
			SetPixel(canvas, x0, y, c)
		}
		return
	}
	
	
	if dy == 0 {
		startX, endX := x0, x1
		if x0 > x1 {
			startX, endX = x1, x0
		}
		for x := startX; x <= endX; x++ {
			SetPixel(canvas, x, y0, c)
		}
		return
	}
	
	
	steep := math.Abs(float64(dy)) > math.Abs(float64(dx))
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	
	
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	
	dx = x1 - x0
	dy = y1 - y0
	
	
	yStep := 1
	if dy < 0 {
		yStep = -1
		dy = -dy
	}
	
	
	d := 2*dy - dx
	y := y0
	
	
	for x := x0; x <= x1; x++ {
		if steep {
			SetPixel(canvas, y, x, c)
		} else {
			SetPixel(canvas, x, y, c)
		}

		if d > 0 {
			y += yStep
			d += 2 * (dy - dx)
		} else {
			d += 2 * dy
		}
	}
}


func MakeCircularBrush(thickness int) [][]bool {
	
	if thickness%2 == 0 {
		thickness++
	}
	
	brush := make([][]bool, thickness)
	radius := thickness / 2
	radiusSq := radius * radius
	
	for y := 0; y < thickness; y++ {
		brush[y] = make([]bool, thickness)
		for x := 0; x < thickness; x++ {
	
			dx := x - radius
			dy := y - radius
			distSq := dx*dx + dy*dy
	
	
			brush[y][x] = distSq <= radiusSq
		}
	}
	
	return brush
}


func ThickLine(canvas [][]color.Color, x0, y0, x1, y1 int, c color.Color, thickness int) {
	
	brush := MakeCircularBrush(thickness)
	radius := thickness / 2
	
	
	linePixels := make([]Point, 0)
	
	
	dx := x1 - x0
	dy := y1 - y0
	
	
	if dx == 0 {
		startY, endY := y0, y1
		if y0 > y1 {
			startY, endY = y1, y0
		}
		for y := startY; y <= endY; y++ {
			linePixels = append(linePixels, Point{X: x0, Y: y})
		}
	} else if dy == 0 { 
		startX, endX := x0, x1
		if x0 > x1 {
			startX, endX = x1, x0
		}
		for x := startX; x <= endX; x++ {
			linePixels = append(linePixels, Point{X: x, Y: y0})
		}
	} else {
		steep := math.Abs(float64(dy)) > math.Abs(float64(dx))
		if steep {
			x0, y0 = y0, x0
			x1, y1 = y1, x1
		}

		if x0 > x1 {
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}

		dx = x1 - x0
		dy = y1 - y0

		yStep := 1
		if dy < 0 {
			yStep = -1
			dy = -dy
		}

		d := 2*dy - dx
		y := y0

		for x := x0; x <= x1; x++ {
			if steep {
				linePixels = append(linePixels, Point{X: y, Y: x})
			} else {
				linePixels = append(linePixels, Point{X: x, Y: y})
			}
	
			if d > 0 {
				y += yStep
				d += 2 * (dy - dx)
			} else {
				d += 2 * dy
			}
		}
	}
	
	
	for _, p := range linePixels {
		for by := 0; by < thickness; by++ {
			for bx := 0; bx < thickness; bx++ {
				if brush[by][bx] {
			
					nx := p.X + bx - radius
					ny := p.Y + by - radius
					SetPixel(canvas, nx, ny, c)
				}
			}
		}
	}
}



func MidpointCircle(canvas [][]color.Color, centerX, centerY, radius int, c color.Color) {
	x := radius
	y := 0
	err := 0
	
	for x >= y {
		SetPixel(canvas, centerX+x, centerY+y, c)
		SetPixel(canvas, centerX+y, centerY+x, c)
		SetPixel(canvas, centerX-y, centerY+x, c)
		SetPixel(canvas, centerX-x, centerY+y, c)
		SetPixel(canvas, centerX-x, centerY-y, c)
		SetPixel(canvas, centerX-y, centerY-x, c)
		SetPixel(canvas, centerX+y, centerY-x, c)
		SetPixel(canvas, centerX+x, centerY-y, c)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}



func XiaolinWuLine(canvas [][]color.Color, x0, y0, x1, y1 int, c color.Color) {
	
	steep := math.Abs(float64(y1-y0)) > math.Abs(float64(x1-x0))
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	gradient := float64(1.0)
	if dx != 0 {
		gradient = dy / dx
	}
	
	
	xend := float64(x0)
	yend := float64(y0) + gradient*(xend-float64(x0))
	xgap := 1.0 - math.Mod(float64(x0)+0.5, 1.0)
	xpxl1 := x0
	ypxl1 := int(yend)
	
	if steep {
		SetPixelWithAlpha(canvas, ypxl1, xpxl1, c, (1-math.Mod(yend, 1.0))*xgap)
		SetPixelWithAlpha(canvas, ypxl1+1, xpxl1, c, math.Mod(yend, 1.0)*xgap)
	} else {
		SetPixelWithAlpha(canvas, xpxl1, ypxl1, c, (1-math.Mod(yend, 1.0))*xgap)
		SetPixelWithAlpha(canvas, xpxl1, ypxl1+1, c, math.Mod(yend, 1.0)*xgap)
	}
	
	
	intery := yend + gradient
	
	
	xend = float64(x1)
	yend = float64(y1) + gradient*(xend-float64(x1))
	xgap = math.Mod(float64(x1)+0.5, 1.0)
	xpxl2 := x1
	ypxl2 := int(yend)
	
	if steep {
		SetPixelWithAlpha(canvas, ypxl2, xpxl2, c, (1-math.Mod(yend, 1.0))*xgap)
		SetPixelWithAlpha(canvas, ypxl2+1, xpxl2, c, math.Mod(yend, 1.0)*xgap)
	} else {
		SetPixelWithAlpha(canvas, xpxl2, ypxl2, c, (1-math.Mod(yend, 1.0))*xgap)
		SetPixelWithAlpha(canvas, xpxl2, ypxl2+1, c, math.Mod(yend, 1.0)*xgap)
	}
	
	
	if steep {
		for x := xpxl1 + 1; x < xpxl2; x++ {
			SetPixelWithAlpha(canvas, int(intery), x, c, 1-math.Mod(intery, 1.0))
			SetPixelWithAlpha(canvas, int(intery)+1, x, c, math.Mod(intery, 1.0))
			intery += gradient
		}
	} else {
		for x := xpxl1 + 1; x < xpxl2; x++ {
			SetPixelWithAlpha(canvas, x, int(intery), c, 1-math.Mod(intery, 1.0))
			SetPixelWithAlpha(canvas, x, int(intery)+1, c, math.Mod(intery, 1.0))
			intery += gradient
		}
	}
}


func XiaolinWuCircle(canvas [][]color.Color, centerX, centerY, radius int, c color.Color) {
	
	x := radius
	y := 0
	err := 0
	
	
	drawPixel := func(xi, yi int, alpha float64) {
		SetPixelWithAlpha(canvas, centerX+xi, centerY+yi, c, alpha)
		SetPixelWithAlpha(canvas, centerX+yi, centerY+xi, c, alpha)
		SetPixelWithAlpha(canvas, centerX-yi, centerY+xi, c, alpha)
		SetPixelWithAlpha(canvas, centerX-xi, centerY+yi, c, alpha)
		SetPixelWithAlpha(canvas, centerX-xi, centerY-yi, c, alpha)
		SetPixelWithAlpha(canvas, centerX-yi, centerY-xi, c, alpha)
		SetPixelWithAlpha(canvas, centerX+yi, centerY-xi, c, alpha)
		SetPixelWithAlpha(canvas, centerX+xi, centerY-yi, c, alpha)
	}
	
	
	drawPixel(x, y, 1.0)
	
	
	for x > y {
		y++

		if err <= 0 {
			err += 2*y + 1
		}

		if err > 0 {
			x--
			err -= 2*x + 1
		}


		if x >= y {
			idealRadius := math.Sqrt(float64(x*x + y*y))
			diff := idealRadius - float64(radius)
			alpha := 1.0 - math.Abs(diff)
			if alpha < 0 {
				alpha = 0
			}
	
			drawPixel(x, y, alpha)
		}
	}
}