package store

import (
	"fmt"
	"math/rand"
	"sync"
)

//DataStore holds the value of the urls
type DataStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

//NewDataStore() creates a DataStore and returns its pointer
func NewDataStore() *DataStore {
	return &DataStore{urls: make(map[string]string)}
}

func (d *DataStore) set(key, val string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, present := d.urls[key]; present {
		return fmt.Errorf("val for key %s already exists")
	}
	d.urls[key] = val
	return nil
}

func (d *DataStore) Get(key string) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if val, ok := d.urls[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("val for key %s not found")
}

func (d *DataStore) Put(val string) string {
	for {
		randKey := generateRandomKey(32)
		if err := d.set(randKey, val); err == nil {
			return randKey
		}
	}
}
func (d *DataStore) GetKeys() []string {
	keys := make([]string, 0, len(d.urls))
	for key := range d.urls {
		keys = append(keys, key)
	}
	return keys
}

func generateRandomKey(len int) string {
	charset := "abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, len)
	for i := range b {
		b[i] = charset[rand.Intn(len)]
	}
	return string(b)
}
