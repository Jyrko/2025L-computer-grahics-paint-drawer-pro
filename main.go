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
	w.Resize(fyne.NewSize(800, 600))
	
	mainUI := ui.NewMainUI(w)
	
	
	mouseHandler := ui.NewMouseHandler(mainUI)
	
	
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		mouseHandler.KeyDown(ke)
	})
	
	
	drawingArea := container.NewStack(mouseHandler, mainUI.Canvas)
	
	
	mainUI.Container = container.NewBorder(
		nil, 
		container.NewHBox(mainUI.StatusLabel), 
		mainUI.ToolsContainer, 
		nil, 
		drawingArea, 
	)
	
	w.SetContent(mainUI.Container)
	w.ShowAndRun()
}