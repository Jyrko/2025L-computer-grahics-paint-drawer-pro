package main

import (
	"paint-drawer-pro/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.NewWithID("com.university.paintdrawerpro")
	w := a.NewWindow("Paint Drawer Pro")
	w.Resize(fyne.NewSize(1024, 768))
	
	mainUI := ui.NewMainUI(w)
	
	
	mouseHandler := ui.NewMouseHandler(mainUI)
	
	
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		mouseHandler.KeyDown(ke)
	})
	
	
	// Create a container that allows the canvas to resize and position correctly
	drawingArea := container.NewStack(mouseHandler, mainUI.Canvas)
	
	// Use a padding container to ensure canvas is properly centered
	paddedDrawingArea := container.NewPadded(drawingArea)
	
	// Setup main container with proper layout
	mainUI.Container = container.NewBorder(
		nil, 
		container.NewHBox(mainUI.StatusLabel), 
		mainUI.ToolsContainer, 
		nil, 
		paddedDrawingArea)
	
	w.SetContent(mainUI.Container)
	w.ShowAndRun()
}