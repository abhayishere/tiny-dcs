package main

import (
	"fmt"
	"net/http"
	"time"
)

type DistributedCache struct {
	hashring          *HashRing
	nodes             map[string]*CacheNode
	activeNodes       map[string]bool
	heartbeatListener chan string
}

func NewDistributedCache() *DistributedCache {
	return &DistributedCache{
		hashring:          NewHashRing(3),
		nodes:             make(map[string]*CacheNode),
		activeNodes:       make(map[string]bool),
		heartbeatListener: make(chan string),
	}
}

func (dc *DistributedCache) ListenToHeartbeat() {
	for {
		select {
		case nodePort := <-dc.heartbeatListener:
			dc.activeNodes[nodePort] = true
			fmt.Println("Node", nodePort, "is up!")
		case <-time.After(10 * time.Second):
			for nodePort := range dc.activeNodes {
				if !dc.activeNodes[nodePort] {
					fmt.Println("Node", nodePort, "is down!")
					delete(dc.activeNodes, nodePort)
				} else {
					dc.activeNodes[nodePort] = false
				}
			}
		}
	}
}

func (dc *DistributedCache) AddNode(address string) {
	node := NewCacheNode(dc)
	node.port = address
	dc.hashring.AddNode(address)
	dc.nodes[address] = node
	go node.Start()
	go node.heartBeat()
}

func (dc *DistributedCache) Set(key, value string) {
	nodePort := dc.hashring.GetNode(key)
	dc.nodes[nodePort].Set(key, value)
}

func (dc *DistributedCache) Get(key string) (string, bool) {
	nodePort := dc.hashring.GetNode(key)
	return dc.nodes[nodePort].Get(key)
}

func main() {
	dc := NewDistributedCache()
	dc.AddNode(":3100")
	dc.AddNode(":3101")
	dc.AddNode(":3102")
	go dc.ListenToHeartbeat()
	go func() {
		time.Sleep(10 * time.Second)
		dc.nodes[":3100"].Stop()
	}()
	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key := r.URL.Query().Get("key")
			if value, exists := dc.Get(key); exists {
				fmt.Fprintf(w, "Value: %s", value)
			} else {
				http.NotFound(w, r)
			}
		case http.MethodPost:
			key := r.FormValue("key")
			value := r.FormValue("value")
			dc.Set(key, value)
			fmt.Fprintf(w, "Stored key: %s, value: %s", key, value)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Central Coordinator running on :8080")
	http.ListenAndServe(":8080", nil)
}
