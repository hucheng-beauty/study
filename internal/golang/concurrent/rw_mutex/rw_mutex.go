package main

import (
	"errors"
	"fmt"
	"sync"
)

var secureMapObject = &secureMap{
	RWMutex: sync.RWMutex{},
	m:       make(map[string]int),
}

type secureMap struct {
	sync.RWMutex
	m map[string]int
}

func (m *secureMap) Get(key string) (int, error) {
	if len(key) <= 0 {
		return 0, errors.New("key invalid")
	}

	m.RLock()
	defer m.RUnlock()
	if value, ok := m.m[key]; ok {
		return value, nil
	}
	return 0, nil
}

func (m *secureMap) Set(key string, value int) (bool, error) {
	if len(key) <= 0 {
		return false, errors.New("key invalid")
	}

	_, isExist := m.m[key]
	if isExist {
		return false, errors.New("key is exist")
	}

	m.Lock()
	m.m[key] = value
	m.Unlock()
	return true, nil
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
