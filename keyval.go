package main

import "sync"

type KeyVal struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKeyVal() *KeyVal {
	return &KeyVal{
		data: make(map[string][]byte),
	}
}

func (kv *KeyVal) Set(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[key] = []byte(value)
}

func (kv *KeyVal) Get(key string) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[key]
	return val, ok
}
