#Distributed Cache System

## Overview

This project implements a distributed cache system using a consistent hashing mechanism and a heartbeat protocol to ensure node availability. The system is designed to manage a distributed cache across multiple nodes and handle node failures gracefully.

## Components

- **CacheNode:** Represents a single node in the cache cluster, responsible for storing and retrieving cache data. It also sends periodic heartbeat messages to the central coordinator to signal that it is alive.
- **HashRing:** Implements a consistent hashing algorithm to distribute cache keys across nodes. This ensures an even distribution of data and minimal reorganization when nodes are added or removed.
- **DistributedCache:** Acts as the central coordinator that manages the nodes, handles cache operations, and monitors node health through heartbeats.

## Features

- **Consistent Hashing:** Distributes cache data across nodes using consistent hashing to ensure even distribution and minimal reorganization.
- **Heartbeat Mechanism:** Nodes periodically send heartbeat messages to the central coordinator to indicate they are alive. The coordinator monitors these messages to detect node failures.
- **Dynamic Node Management:** Nodes can be added or removed dynamically. The system automatically redistributes cache data as needed.
- **Fault Tolerance:** The system detects node failures and updates its state accordingly, ensuring continued operation with the remaining nodes.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/distributed-cache-system.git
   cd distributed-cache-system
   ```
2. **Build the Project:**
   ```bash
   go run .
   ```
3. **Testing Node Failures:**

   The system automatically stops one of the nodes (port 3100) after 10 seconds to simulate a node failure. This demonstrates the heartbeat monitoring feature.

## Usage
### Adding a Node
To add a new cache node to the distributed system:
```go
dc.AddNode(":3103")
```
### Storing Data
To store a key-value pair in the distributed cache:
```http
POST /cache?key=mykey&value=myvalue
```
### Retrieving Data
To retrieve a value by its key:
```http
GET /cache?key=mykey
```
### Node Heartbeat
Each node sends a heartbeat signal every 5 seconds to the central coordinator to indicate that it is still operational. If a node fails to send a heartbeat within 10 seconds, it is considered down and removed from the active node list.
#### Example Output
```terminal
Node :3100 is sending heartbeat!
Cache Node is running on port: :3101
Cache Node is running on port: :3100
Node :3100 is up!
Central Coordinator running on :8080
Cache Node is running on port: :3102
Node :3101 is sending heartbeat!
Node :3101 is up!
Node :3102 is sending heartbeat!
Node :3102 is up!
Stopping node on port: :3100
Node :3100 is down!
```
## Contributing
Contributions are welcome! Please fork this repository and submit a pull request for any enhancements or bug fixes.
