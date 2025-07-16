/*
==========================================================================================
  File:        main_gui.go
  Last Update: 2024-05-18
  Author:      Haitam Bidiouane (@sh0penheimer)
  Ownership:   © Haitam Bidiouane. All rights reserved.
------------------------------------------------------------------------------------------
  Scope:
    GUI entry point for the blockchain websocket gateway. Launches a Fyne-based desktop
    application for managing the gateway, with placeholders for controls and status display.
    Designed for future integration with the modular gateway orchestration layer.
==========================================================================================
*/

package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Blockchain Websocket Gateway GUI")

	w.SetContent(
		container.NewVBox(
			widget.NewLabel("Gateway GUI Placeholder"),
			widget.NewLabel("(Controls and status will appear here.)"),
		),
	)

	w.Resize(fyne.NewSize(400, 200))
	w.ShowAndRun()
} 