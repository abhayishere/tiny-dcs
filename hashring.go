package main

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashRing struct {
	nodes       []int
	nodeMap     map[int]string
	numReplicas int
}

func NewHashRing(numReplicas int) *HashRing {
	return &HashRing{
		nodeMap:     make(map[int]string),
		numReplicas: numReplicas,
	}
}

func (hr *HashRing) AddNde(node string) {
	for i := 0; i < hr.numReplicas; i++ {
		hash := int(crc32.ChecksumIEEE([]byte(node + strconv.Itoa(i))))
		hr.nodes=append(hr.nodes, hash)
		hr.nodeMap[hash]=node
	}
	sort.Ints(hr.nodes)
}

func (hr *HashRing)GetNode(key string)string{
	if len(hr.nodes)==0{
		return ""
	}
	hash:=int(crc32.ChecksumIEEE([]byte(key)))
	idx:=sort.Search(len(hr.nodes),func(i int)bool{return hr.nodes[i]>=hash})
	if idx==len(hr.nodes){
		idx=0;
	}
	return hr.nodeMap[hr.nodes[idx]]
}
