package main

import (
	"flag"
    "fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sch0penheimer/eth-ws-server/blockchain"
	"github.com/sch0penheimer/eth-ws-server/websocket"
)

func main() {
	hostFlag := flag.String("host", "", "Host IP address of the Ethereum node cluster (e.g., 192.168.240.222)")
	flag.Parse()

	// Validate the input
    if *hostFlag == "" {
        log.Fatal("No host IP address provided. Use the --host flag to specify the host IP.")
    }

    // Dynamically format the node URLs using the host IP
    nodeURLs := []string{
        fmt.Sprintf("ws://%s:8545", *hostFlag),
        fmt.Sprintf("http://%s:8546", *hostFlag),
        fmt.Sprintf("http://%s:8547", *hostFlag),
        fmt.Sprintf("http://%s:8548", *hostFlag),
        fmt.Sprintf("http://%s:8549", *hostFlag),
    }

    log.Printf("Using the following WebSocket node URLs: %v", nodeURLs)

	miningController, err := blockchain.NewMiningController(nodeURLs)
	if err != nil {
		log.Fatalf("Failed to initialize mining controller: %v", err)
	}

	// Initialize blockchain client (for the first node)
	blockFetcher, err := blockchain.NewBlockFetcher(nodeURLs[0])
	if err != nil {
		log.Fatalf("Failed to initialize blockchain client: %v", err)
	}

	// Setup WebSocket handler
	wsHandler := websocket.NewWSHandler(blockFetcher, miningController)

	// Create router
	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", wsHandler.HandleConnections)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// CORS middleware
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
