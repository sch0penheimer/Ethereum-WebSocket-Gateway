/*
==========================================================================================
  File:        main.go
  Last Update: 2024-05-18
  Author:      Haitam Bidiouane (@sh0penheimer)
  Ownership:   © Haitam Bidiouane. All rights reserved.
------------------------------------------------------------------------------------------
  Scope:
    CLI entry point for the blockchain websocket gateway. Parses command-line flags,
    initializes the Gateway orchestration layer, and starts the HTTP server for websocket
    and health endpoints. Designed to be used as a CLI or as a reference for GUI startup.
==========================================================================================
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sch0penheimer/eth-ws-server/internal/gateway"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `\nBlockchain Websocket Gateway\n\n`)
	fmt.Fprintf(os.Stderr, `Usage: %s --nodes N --addresses IP1,IP2,... --ports PORT1,PORT2,...\n`, os.Args[0])
	fmt.Fprintf(os.Stderr, `\nFlags:\n`)
	flag.PrintDefaults()
}

func main() {
	nodeCount := flag.Int("nodes", 0, "Total number of nodes (required)")
	nodeAddresses := flag.String("addresses", "", "Comma-separated list of node IP addresses (required)")
	nodePorts := flag.String("ports", "", "Comma-separated list of node ports (required)")
	help := flag.Bool("help", false, "Show help message")
	flag.Usage = printUsage
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	if *nodeCount <= 0 {
		fmt.Fprintln(os.Stderr, "Error: --nodes must be greater than 0.")
		printUsage()
		os.Exit(1)
	}
	if *nodeAddresses == "" || *nodePorts == "" {
		fmt.Fprintln(os.Stderr, "Error: --addresses and --ports are required.")
		printUsage()
		os.Exit(1)
	}

	addressList := strings.Split(*nodeAddresses, ",")
	portList := strings.Split(*nodePorts, ",")

	if len(addressList) != *nodeCount || len(portList) != *nodeCount {
		fmt.Fprintf(os.Stderr, "Error: The number of addresses (%d) and ports (%d) must match the node count (%d).\n", len(addressList), len(portList), *nodeCount)
		printUsage()
		os.Exit(1)
	}

	cfg := gateway.GatewayConfig{
		NodeCount:     *nodeCount,
		NodeAddresses: addressList,
		NodePorts:     portList,
	}
	gw, err := gateway.NewGateway(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}
	gw.Start()
	log.Println(gw.Status())

	// Set up the HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/ws", gw.WSHandler().HandleConnections)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Add CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			next.ServeHTTP(w, r)
		})
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
