package types

import "sync"

type AtomicBool interface {
	Get() bool
	Set(value bool)
}

type atomicBool struct {
	value bool
	mutex *sync.RWMutex
}

func NewAtomicBool(initValue bool) AtomicBool {
	return &atomicBool{
		value: initValue,
		mutex: &sync.RWMutex{},
	}
}

func (at *atomicBool) Get() bool {
	at.mutex.RLock()
	defer at.mutex.RUnlock()
	return at.value
}

func (at *atomicBool) Set(value bool) {
	at.mutex.Lock()
	defer at.mutex.Unlock()
	at.value = value
}
