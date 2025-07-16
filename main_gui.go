/*
==========================================================================================
  File:        main_gui.go
  Last Update: 2024-05-18
  Author:      Haitam Bidiouane (@sh0penheimer)
  Ownership:   © Haitam Bidiouane. All rights reserved.
------------------------------------------------------------------------------------------
  Scope:
    GUI entry point for the blockchain websocket gateway. Launches a Fyne-based desktop
    application for managing the gateway, with controls for configuration, start/stop,
    and status display. Designed for future integration with the modular gateway orchestration layer.
==========================================================================================
*/

package main

import (
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Blockchain Websocket Gateway GUI")

	// Input fields
	nodeCountEntry := widget.NewEntry()
	nodeCountEntry.SetPlaceHolder("Node count (e.g. 3)")

	nodeAddressesEntry := widget.NewEntry()
	nodeAddressesEntry.SetPlaceHolder("Comma-separated node addresses (e.g. 127.0.0.1,192.168.1.2)")

	nodePortsEntry := widget.NewEntry()
	nodePortsEntry.SetPlaceHolder("Comma-separated node ports (e.g. 8546,8545)")

	statusLabel := widget.NewLabel("Status: stopped")

	// Start/Stop buttons
	startBtn := widget.NewButton("Start Gateway", func() {
		statusLabel.SetText("Status: running (not yet wired up)")
	})
	stopBtn := widget.NewButton("Stop Gateway", func() {
		statusLabel.SetText("Status: stopped (not yet wired up)")
	})
	stopBtn.Disable() // Only enable after start in future logic

	// Layout
	form := container.NewVBox(
		widget.NewLabel("Blockchain Websocket Gateway"),
		widget.NewForm(
			widget.NewFormItem("Node Count", nodeCountEntry),
			widget.NewFormItem("Node Addresses", nodeAddressesEntry),
			widget.NewFormItem("Node Ports", nodePortsEntry),
		),
		container.NewHBox(startBtn, stopBtn),
		statusLabel,
	)

	w.SetContent(form)
	w.Resize(fyne.NewSize(500, 300))
	w.ShowAndRun()
} 