package cache

import (
	"sync"
)

type InMemory[keyType comparable, valueType any] interface {
	Get(key keyType) (value valueType, ok bool)
	GetList(keys []keyType) (values []valueType)
	Save(key keyType, value valueType)
}

type inMemory[k comparable, v any] struct {
	mutex sync.RWMutex
	data  map[k]v
}

func New[keyType comparable, valueType any]() InMemory[keyType, valueType] {
	return &inMemory[keyType, valueType]{
		data: make(map[keyType]valueType),
	}
}

func (c *inMemory[k, v]) Get(key k) (value v, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, ok = c.data[key]

	return value, ok
}

func (c *inMemory[k, v]) GetList(keys []k) []v {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	values := make([]v, 0, len(keys))
	for i := range keys {
		value, ok := c.data[keys[i]]
		if ok {
			values = append(values, value)
		}
	}
	return values
}

func (c *inMemory[k, v]) Save(key k, value v) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}
