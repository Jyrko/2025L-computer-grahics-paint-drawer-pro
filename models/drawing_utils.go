package models

import (
	"image/color"
	"paint-drawer-pro/algorithms"
)




func drawMidpointLine(canvas [][]color.Color, x0, y0, x1, y1 int, c color.Color) {
	algorithms.MidpointLine(canvas, x0, y0, x1, y1, c)
}


func drawThickLine(canvas [][]color.Color, x0, y0, x1, y1 int, c color.Color, thickness int) {
	algorithms.ThickLine(canvas, x0, y0, x1, y1, c, thickness)
}


func drawMidpointCircle(canvas [][]color.Color, centerX, centerY, radius int, c color.Color) {
	algorithms.MidpointCircle(canvas, centerX, centerY, radius, c)
}


func drawXiaolinWuLine(canvas [][]color.Color, x0, y0, x1, y1 int, c color.Color) {
	algorithms.XiaolinWuLine(canvas, x0, y0, x1, y1, c)
}


func drawXiaolinWuCircle(canvas [][]color.Color, centerX, centerY, radius int, c color.Color) {
	algorithms.XiaolinWuCircle(canvas, centerX, centerY, radius, c)
}