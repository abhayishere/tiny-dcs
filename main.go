package main

import (
	"fmt"
	"net/http"
)

type DistributedCache struct {
	hashring *HashRing
	nodes    map[string]*CacheNode
}

func NewDistributedCache() *DistributedCache {
	return &DistributedCache{
		hashring: NewHashRing(3),
		nodes:    make(map[string]*CacheNode),
	}
}

func (dc *DistributedCache) AddNode(address string) {
	node := NewCacheNode()
	node.port = address
	dc.hashring.AddNde(address)
	dc.nodes[address] = node
	go node.Start()
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
