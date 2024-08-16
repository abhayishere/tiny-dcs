package main

import (
	"fmt"
	"net/http"
	"sync"
)

type CacheNode struct {
	data map[string]string
	port string
	mu   sync.RWMutex
}

func NewCacheNode() *CacheNode {
	return &CacheNode{data: make(map[string]string)}
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
			fmt.Println(w, "Value ", value, "was stored in ",node.port)
		}else{
			http.NotFound(w,r)
		}
	case http.MethodPost:
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		node.Set(key,value)
		fmt.Fprintf(w, "Stored key: %s, value: %s", key, value)
	default:
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
}

func (node *CacheNode)Start(){
	http.HandleFunc("/"+node.port,node.handler)
	fmt.Println("Cache Node is running on port :",node.port)
	http.ListenAndServe(node.port,nil)
}
