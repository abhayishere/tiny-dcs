package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CacheNode struct {
	data   map[string]string
	port   string
	mu     sync.RWMutex
	server *http.Server
	dc     *DistributedCache
	done   chan bool
}

func NewCacheNode(dc *DistributedCache) *CacheNode {
	return &CacheNode{
		data: make(map[string]string),
		dc:   dc,
		done: make(chan bool),
	}
}

func (node *CacheNode) heartBeat() {
	for {
		select {
		case <-node.done:
			fmt.Println("Node", node.port, "stopped sending heartbeat!")
			return
		default:
			node.dc.heartbeatListener <- node.port
			time.Sleep(5 * time.Second)
		}
	}
}

func (node *CacheNode) Set(key, value string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.data[key] = value
}

func (node *CacheNode) Get(key string) (string, bool) {
	node.mu.RLock()
	defer node.mu.RUnlock()
	value, ok := node.data[key]
	return value, ok
}

func (node *CacheNode) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := r.URL.Query().Get("key")
		if value, exists := node.Get(key); exists {
			fmt.Println(w, "Value ", value, "was stored in ", node.port)
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		node.Set(key, value)
		fmt.Fprintf(w, "Stored key: %s, value: %s", key, value)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (node *CacheNode) Start() {
	server := &http.Server{
		Addr:    node.port,
		Handler: http.HandlerFunc(node.handler),
	}
	node.server = server

	fmt.Println("Cache Node is running on port:", node.port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("ListenAndServe(): %s\n", err)
	}
}

func (node *CacheNode) Stop() {
	if node.server != nil {
		node.done<-true
		fmt.Println("Stopping node on port:", node.port)
		node.server.Close()
	}
}
