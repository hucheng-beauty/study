package main

import (
	"errors"
	"fmt"
	"sync"
)

var secureMap = &sMap{
	RWMutex: sync.RWMutex{},
	m:       make(map[string]int),
}

type sMap struct {
	sync.RWMutex
	m map[string]int
}

func (m *sMap) Get(key string) (int, error) {
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

func (m *sMap) Set(key string, value int) (bool, error) {
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
	for k, v := range map[string]int{
		"a": 97,
		"b": 98,
		"c": 99,
	} {
		_, _ = secureMap.Set(k, v)
	}
	fmt.Println("end")
}
