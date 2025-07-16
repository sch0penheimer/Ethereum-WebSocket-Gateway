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
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sch0penheimer/eth-ws-server/internal/gateway"
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

	var gw *gateway.Gateway

	// Start/Stop buttons
	startBtn := widget.NewButton("Start Gateway", func() {
		// Validate input
		nodeCount, err := strconv.Atoi(nodeCountEntry.Text)
		if err != nil || nodeCount <= 0 {
			dialog.ShowError(
				fyne.NewError("Invalid Input", "Node count must be a positive integer."), w)
			return
		}
		addresses := strings.Split(nodeAddressesEntry.Text, ",")
		ports := strings.Split(nodePortsEntry.Text, ",")
		if len(addresses) != nodeCount || len(ports) != nodeCount {
			dialog.ShowError(
				fyne.NewError("Invalid Input", "Number of addresses and ports must match node count."), w)
			return
		}
		cfg := gateway.GatewayConfig{
			NodeCount:     nodeCount,
			NodeAddresses: addresses,
			NodePorts:     ports,
		}
		var gwErr error
		gw, gwErr = gateway.NewGateway(cfg)
		if gwErr != nil {
			statusLabel.SetText("Status: error")
			dialog.ShowError(
				fyne.NewError("Gateway Error", gwErr.Error()), w)
			return
		}
		gw.Start()
		statusLabel.SetText(gw.Status())
		startBtn.Disable()
		stopBtn.Enable()
	})
	stopBtn := widget.NewButton("Stop Gateway", func() {
		if gw != nil {
			gw.Stop()
			statusLabel.SetText(gw.Status())
		}
		stopBtn.Disable()
		startBtn.Enable()
	})
	stopBtn.Disable() // Only enable after start

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