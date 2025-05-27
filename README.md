# Ethereum WebSocket Gateway

A real-time WebSocket gateway service that provides unified access to multiple Ethereum blockchain nodes, enabling live block monitoring and distributed mining control through a single WebSocket interface.

## Overview

The Ethereum WebSocket Gateway implements a hub-and-spoke architecture where a central gateway service manages connections between WebSocket clients and multiple Ethereum nodes. The system uses a dual-protocol strategy: WebSocket connections for real-time data streaming and HTTP connections for mining control operations.

## Features

- **Real-time Block Streaming**: Subscribe to new block notifications with automatic broadcasting to connected clients
- **Multi-node Mining Control**: Start/stop mining operations across multiple Ethereum nodes simultaneously  
- **Network Metrics**: Live network statistics including block time, difficulty, hashrate, and latency
- **Historical Block Data**: Fetch latest blocks with transaction details and network metrics
- **WebSocket API**: JSON-based message protocol for all client interactions
- **CORS Support**: Cross-origin resource sharing enabled for web applications

## Architecture

### Core Components

The system consists of three primary components:

1. **Main Application** (`main.go`): 
   - HTTP server setup and routing
   - Command-line argument parsing for node configuration
   - Component initialization and dependency injection

2. **WebSocket Handler** (`websocket/websocket.go`):
   - Client connection management
   - Message routing and processing
   - Real-time block subscription handling

3. **Blockchain Services** (`blockchain/blockchain.go`): 
   - `BlockFetcher`: Ethereum node interaction for block data
   - `MiningController`: Multi-node mining operations

### Node Configuration

The gateway supports dynamic configuration of multiple Ethereum nodes: 

- First node (index 0): WebSocket connection for block subscriptions
- Additional nodes: HTTP connections for mining control only

## Installation & Usage

### Prerequisites

- Go 1.19 or higher
- Access to Ethereum nodes with WebSocket and HTTP RPC endpoints

### Running the Gateway

```bash
go run main.go --nodes=3 --addresses=192.168.1.10,192.168.1.11,192.168.1.12 --ports=8545,8545,8545
```

**Command-line Arguments:**
- `--nodes`: Total number of Ethereum nodes
- `--addresses`: Comma-separated list of node IP addresses  
- `--ports`: Comma-separated list of node RPC ports

The server starts on port 8080 with the following endpoints:
- `ws://localhost:8080/ws` - WebSocket connection
- `http://localhost:8080/health` - Health check endpoint

## WebSocket API

### Message Types

The gateway processes four distinct message types:

#### 1. Subscribe to New Blocks
```json
{"type": "subscribe"}
```
Response: `{"type": "subscribe", "status": true, "message": "Subscription status updated"}`

#### 2. Get Latest Blocks
```json
{"type": "latestblocks", "payload": {"count": 5}}
```
Response: 

#### 3. Get Mining Status
```json
{"type": "miningstatus"}
```
Response: `{"type": "miningStatus", "data": [true, false, true]}`

#### 4. Toggle Mining
```json
{"type": "togglemining", "payload": {"start": true}}
```
Response: `{"type": "toggleMining", "data": [true, true, true]}`

### Real-time Block Broadcasting

When subscribed, clients automatically receive new block notifications:

The system uses Ethereum's `SubscribeNewHead()` to monitor new blocks and broadcasts them to all subscribed clients with full block details and network metrics.

Looking at your request, you want markdown documentation for the **Block Structure** and **Network Metrics** data structures from the Ethereum WebSocket Gateway codebase.

## Data Structures

### Block Structure

The `Block` struct represents a complete Ethereum block with comprehensive metadata and transaction details:

**Key Fields:**
- **Block Identifiers**: `Number`, `Hash`, `ParentHash` for blockchain navigation
- **Merkle Roots**: `Sha3Uncles`, `TransactionRoot` for cryptographic verification  
- **Consensus Data**: `Validator` (fetched via `clique_getSigner` RPC), `Difficulty`
- **Gas Metrics**: `GasUsed`, `GasLimit` for network capacity tracking
- **Transaction Data**: `Transactions` array with detailed transaction information, `TransactionCount`, `TotalFees`
- **Block Metadata**: `Timestamp`, `Size` for block analysis

**Associated Transaction Structure:**

The block structure is populated by `GetBlockByNumber()` which fetches raw Ethereum block data and enriches it with calculated fields like `TotalFees` and validator information .

### Network Metrics

The network metrics provide real-time blockchain performance and system health data returned by `GetNetworkMetrics()`:

**Metrics Included:**
- **`averageBlockTime`**: Calculated from the last 25 blocks, returned in minutes
- **`difficulty`**: Current network difficulty from the latest block
- **`hashrate`**: Network hashrate via `eth_hashrate` RPC call
- **`latency`**: Network latency (currently mocked at 50ms)
- **`memoryUsage`**: Gateway application memory consumption in MB

These metrics are included in both `latestBlocks` responses and real-time `newBlock` broadcasts to WebSocket clients .

## Concurrency Model

The gateway implements a concurrent architecture with dedicated goroutines:

- `run()`: Client registration/unregistration management
- `watchNewBlocks()`: Real-time block subscription handling  
- `readPump()`/`writePump()`: Per-client message processing

## Development

### Project Structure
```
├── main.go                 # Application entry point and configuration
├── blockchain/
│   └── blockchain.go       # Ethereum client and mining controller
└── websocket/
    └── websocket.go        # WebSocket handler and client management
```

### Key Dependencies
- `github.com/ethereum/go-ethereum` - Ethereum client library
- `github.com/gorilla/websocket` - WebSocket implementation
- `github.com/gorilla/mux` - HTTP router

## Notes

The system is designed for Ethereum networks using the Clique consensus mechanism, as evidenced by the validator fetching using `clique_getSigner` RPC calls. The first node in the configuration must support WebSocket connections for real-time block subscriptions, while additional nodes only require HTTP RPC access for mining control operations.

Associated Wiki:
- [Wiki (sch0penheimer/Ethereum-WebSocket-Gateway)](https://deepwiki.com/sch0penheimer/Ethereum-WebSocket-Gateway)
