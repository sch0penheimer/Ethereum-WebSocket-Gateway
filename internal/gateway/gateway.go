/*
==========================================================================================
  File:        gateway.go
  Last Update: 2024-05-18
  Author:      Haitam Bidiouane (@sh0penheimer)
  Ownership:   © Haitam Bidiouane. All rights reserved.
------------------------------------------------------------------------------------------
  Scope:
    Provides the Gateway struct and orchestration logic for the blockchain websocket gateway.
    Acts as a facade for configuration, startup, shutdown, and status, using the blockchain
    and websocket packages. Designed for use by both CLI and GUI entry points.
==========================================================================================
*/

package gateway

import (
    "fmt"
    "strings"
    "github.com/sch0penheimer/eth-ws-server/blockchain"
    "github.com/sch0penheimer/eth-ws-server/websocket"
)

type GatewayConfig struct {
    NodeCount     int
    NodeAddresses []string
    NodePorts     []string
    // Add more config fields as needed (e.g., listen port, log level, etc.)
}

type Gateway struct {
    config           GatewayConfig
    blockFetcher     *blockchain.BlockFetcher
    miningController *blockchain.MiningController
    wsHandler        *websocket.WSHandler
    running          bool
}

func NewGateway(cfg GatewayConfig) (*Gateway, error) {
    if cfg.NodeCount <= 0 {
        return nil, fmt.Errorf("invalid node count")
    }
    if len(cfg.NodeAddresses) != cfg.NodeCount || len(cfg.NodePorts) != cfg.NodeCount {
        return nil, fmt.Errorf("addresses and ports must match node count")
    }
    nodeURLs := make([]string, cfg.NodeCount)
    for i := 0; i < cfg.NodeCount; i++ {
        if i == 0 {
            nodeURLs[i] = fmt.Sprintf("ws://%s:%s", cfg.NodeAddresses[i], cfg.NodePorts[i])
        } else {
            nodeURLs[i] = fmt.Sprintf("http://%s:%s", cfg.NodeAddresses[i], cfg.NodePorts[i])
        }
    }
    miningController, err := blockchain.NewMiningController(nodeURLs)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize mining controller: %w", err)
    }
    blockFetcher, err := blockchain.NewBlockFetcher(nodeURLs[0])
    if err != nil {
        return nil, fmt.Errorf("failed to initialize block fetcher: %w", err)
    }
    wsHandler := websocket.NewWSHandler(blockFetcher, miningController)
    return &Gateway{
        config:           cfg,
        blockFetcher:     blockFetcher,
        miningController: miningController,
        wsHandler:        wsHandler,
        running:          false,
    }, nil
}

// Start launches the gateway (placeholder for future expansion)
func (g *Gateway) Start() error {
    g.running = true
    // In CLI: would start HTTP server, etc. In GUI: would trigger server start.
    return nil
}

// Stop shuts down the gateway (placeholder for future expansion)
func (g *Gateway) Stop() error {
    g.running = false
    // Add logic to gracefully stop HTTP server, close connections, etc.
    return nil
}

// Status returns a string summary of the gateway's current state
func (g *Gateway) Status() string {
    status := "stopped"
    if g.running {
        status = "running"
    }
    return fmt.Sprintf("Gateway status: %s | Nodes: %d | Addresses: %s | Ports: %s", status, g.config.NodeCount, strings.Join(g.config.NodeAddresses, ","), strings.Join(g.config.NodePorts, ","))
}

// Expose accessors for blockFetcher, miningController, wsHandler as needed for CLI/GUI
func (g *Gateway) BlockFetcher() *blockchain.BlockFetcher { return g.blockFetcher }
func (g *Gateway) MiningController() *blockchain.MiningController { return g.miningController }
func (g *Gateway) WSHandler() *websocket.WSHandler { return g.wsHandler }








