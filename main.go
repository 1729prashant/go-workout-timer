package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&appTheme{})

	win := a.NewWindow("Workout Pause Timer")
	win.Resize(fyne.NewSize(1000, 640))

	state := NewAppState()
	ui := NewUI(win, state)
	ui.StartTicker()

	win.ShowAndRun()
}
