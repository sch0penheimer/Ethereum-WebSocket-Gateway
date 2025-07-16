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
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sch0penheimer/eth-ws-server/internal/gateway"
)

func main() {
	nodeCount := flag.Int("nodes", 0, "Total number of nodes")
	nodeAddresses := flag.String("addresses", "", "Comma-separated list of node IP addresses")
	nodePorts := flag.String("ports", "", "Comma-separated list of node ports")
	flag.Parse()

	if *nodeCount <= 0 {
		log.Fatal("Invalid or missing node count. Use the --nodes flag to specify the total number of nodes.")
	}
	if *nodeAddresses == "" || *nodePorts == "" {
		log.Fatal("Node addresses and ports must be provided. Use --addresses and --ports flags.")
	}

	addressList := strings.Split(*nodeAddresses, ",")
	portList := strings.Split(*nodePorts, ",")

	if len(addressList) != *nodeCount || len(portList) != *nodeCount {
		log.Fatalf("The number of addresses (%d) and ports (%d) must match the node count (%d).", len(addressList), len(portList), *nodeCount)
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
