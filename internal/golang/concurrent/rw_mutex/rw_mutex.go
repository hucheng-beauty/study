package main

import (
	"fmt"
	"log"
	"sync"
)

var secureMapObject *secureMap = NewSecureMap()

type secureMap struct {
	mutex sync.RWMutex
	m     map[string]int
}

func (m *secureMap) Get(key string) (int, bool) {
	if len(key) <= 0 {
		return -1, false
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *secureMap) Set(key string, value int) bool {
	if len(key) <= 0 {
		log.Println("key invalid")
		return false
	}

	_, isExist := m.m[key]
	if isExist {
		log.Printf("key is exist,key: %s\n", key)
	}

	m.mutex.Lock()
	m.m[key] = value
	m.mutex.Unlock()
	return true
}

func NewSecureMap() *secureMap {
	return &secureMap{
		mutex: sync.RWMutex{},
		m:     make(map[string]int),
	}
}

func main() {
	t := map[string]int{
		"a": 97,
		"b": 98,
		"c": 99,
	}

	for k, v := range t {
		secureMapObject.Set(k, v)
	}
	fmt.Println("end")
}
