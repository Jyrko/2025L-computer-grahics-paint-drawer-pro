package ui

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"paint-drawer-pro/models"
)


func SerializeColor(c color.Color) map[string]interface{} {
	if c == nil {
		return map[string]interface{}{
			"R": 0,
			"G": 0,
			"B": 0,
			"A": 0,
		}
	}
	
	r, g, b, a := c.RGBA()
	return map[string]interface{}{
		"R": uint8(r),
		"G": uint8(g),
		"B": uint8(b),
		"A": uint8(a),
	}
}


func DeserializeColor(colorMap map[string]interface{}) color.Color {
	r := uint8(colorMap["R"].(float64))
	g := uint8(colorMap["G"].(float64))
	b := uint8(colorMap["B"].(float64))
	a := uint8(colorMap["A"].(float64))
	
	return color.RGBA{r, g, b, a}
}


func SerializePoint(p models.Point) map[string]interface{} {
	return map[string]interface{}{
		"X": p.X,
		"Y": p.Y,
	}
}


func DeserializePoint(pointMap map[string]interface{}) models.Point {
	return models.Point{
		X: int(pointMap["X"].(float64)),
		Y: int(pointMap["Y"].(float64)),
	}
}


func (ui *MainUI) SaveShapesToFile(filePath string) error {
	
	shapesData := make([]map[string]interface{}, 0, len(ui.State.Shapes))
	
	for _, shape := range ui.State.Shapes {
		shapeData := shape.Serialize()
		shapesData = append(shapesData, shapeData)
	}
	
	
	data := map[string]interface{}{
		"shapes": shapesData,
	}
	
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing shapes: %v", err)
	}
	
	
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	
	return nil
}


func (ui *MainUI) LoadShapesFromFile(filePath string) error {
	
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	
	
	var data map[string]interface{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}
	
	
	shapesData, ok := data["shapes"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid shapes data format")
	}
	
	
	ui.State.Shapes = []models.Shape{}
	
	
	for _, shapeData := range shapesData {
		shapeMap, ok := shapeData.(map[string]interface{})
		if !ok {
			continue
		}
		
		
		shapeType, ok := shapeMap["type"].(string)
		if !ok {
			continue
		}
		
		var shape models.Shape
		
		switch shapeType {
		case "circle":
			shape = deserializeCircle(shapeMap)
		case "line":
			shape = deserializeLine(shapeMap)
		case "polygon":
			shape = deserializePolygon(shapeMap)
		case "rectangle":
			shape = deserializeRectangle(shapeMap)
		case "pill":
			shape = deserializePill(shapeMap)
		default:
			continue
		}
		
		if shape != nil {
			ui.State.Shapes = append(ui.State.Shapes, shape)
		}
	}
	
	ui.Canvas.Refresh()
	return nil
}


func deserializeCircle(data map[string]interface{}) *models.Circle {
	centerMap, ok := data["center"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	center := DeserializePoint(centerMap)
	radius := int(data["radius"].(float64))
	colorMap, ok := data["color"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	color := DeserializeColor(colorMap)
	
	return models.NewCircle(center, radius, color)
}

func deserializeLine(data map[string]interface{}) *models.Line {
	startMap, ok := data["start"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	endMap, ok := data["end"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	start := DeserializePoint(startMap)
	end := DeserializePoint(endMap)
	
	colorMap, ok := data["color"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	color := DeserializeColor(colorMap)
	thickness := int(data["thickness"].(float64))
	penType := data["penType"].(string)
	
	return models.NewLine(start, end, color, thickness, penType)
}

func deserializePolygon(data map[string]interface{}) *models.Polygon {
	verticesData, ok := data["vertices"].([]interface{})
	if !ok {
		return nil
	}
	
	vertices := make([]models.Point, len(verticesData))
	for i, vData := range verticesData {
		vMap, ok := vData.(map[string]interface{})
		if !ok {
			return nil
		}
		vertices[i] = DeserializePoint(vMap)
	}
	
	colorMap, ok := data["color"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	color := DeserializeColor(colorMap)
	thickness := int(data["thickness"].(float64))
	
	polygon := models.NewPolygon(vertices, color, thickness)
	
	
	isFilled, ok := data["isFilled"].(bool)
	if ok && isFilled {
		useImage, ok := data["useImage"].(bool)
		if ok && useImage {
			
			
			fillColorMap, ok := data["fillColor"].(map[string]interface{})
			if ok {
				fillColor := DeserializeColor(fillColorMap)
				polygon.SetFillColor(fillColor)
			}
		} else {
			fillColorMap, ok := data["fillColor"].(map[string]interface{})
			if ok {
				fillColor := DeserializeColor(fillColorMap)
				polygon.SetFillColor(fillColor)
			}
		}
	}
	
	return polygon
}

func deserializeRectangle(data map[string]interface{}) *models.Rectangle {
	topLeftMap, ok := data["topLeft"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	bottomRightMap, ok := data["bottomRight"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	topLeft := DeserializePoint(topLeftMap)
	bottomRight := DeserializePoint(bottomRightMap)
	
	colorMap, ok := data["color"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	color := DeserializeColor(colorMap)
	thickness := int(data["thickness"].(float64))
	
	rectangle := models.NewRectangle(topLeft, bottomRight, color, thickness)
	
	
	isFilled, ok := data["isFilled"].(bool)
	if ok && isFilled {
		useImage, ok := data["useImage"].(bool)
		if ok && useImage {
			
			
			fillColorMap, ok := data["fillColor"].(map[string]interface{})
			if ok {
				fillColor := DeserializeColor(fillColorMap)
				rectangle.SetFillColor(fillColor)
			}
		} else {
			fillColorMap, ok := data["fillColor"].(map[string]interface{})
			if ok {
				fillColor := DeserializeColor(fillColorMap)
				rectangle.SetFillColor(fillColor)
			}
		}
	}
	
	return rectangle
}

func deserializePill(data map[string]interface{}) *models.Pill {
	startMap, ok := data["start"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	endMap, ok := data["end"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	start := DeserializePoint(startMap)
	end := DeserializePoint(endMap)
	
	colorMap, ok := data["color"].(map[string]interface{})
	if !ok {
		return nil
	}
	
	color := DeserializeColor(colorMap)
	radius := int(data["radius"].(float64))
	
	pill := models.NewPill(start, radius, color)
	pill.End = end
	pill.Step = 3
	
	return pill
}
