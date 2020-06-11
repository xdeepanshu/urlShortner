package store

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
)

//DataStore holds the value of the urls
type DataStore struct {
	urls map[string]string
	mu   sync.RWMutex
	file *os.File
	l    *log.Logger
}

type record struct {
	Key, Val string
}

//NewDataStore() creates a DataStore and returns its pointer
func NewDataStore(filename string, logger *log.Logger) *DataStore {
	ds := &DataStore{urls: make(map[string]string)}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		ds.l.Fatal("DataStore : ", err)
	}
	ds.file, ds.l = f, logger
	if err := ds.load(); err != nil {
		ds.l.Printf("DataStore error while loading [%s]\n", err)
	}
	return ds

}

func (d *DataStore) save(key, val string) error {
	e := gob.NewEncoder(d.file)
	return e.Encode(record{key, val})
}

func (d *DataStore) load() error {
	if _, err := d.file.Seek(0, 0); err != nil {
		return err
	}
	decoder := gob.NewDecoder(d.file)
	var err error
	for err == nil {
		var r record
		if err = decoder.Decode(&r); err == nil {
			d.set(r.Key, r.Val)
		}
		if err == io.EOF {
			return nil
		}
	}
	return err
}

func (d *DataStore) set(key, val string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, present := d.urls[key]; present {
		d.l.Printf("val for key %s already exists\n", key)
		return fmt.Errorf("val for key %s already exists\n", key)
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
	d.l.Printf("val for key %s not found\n", key)
	return "", fmt.Errorf("val for key %s not found\n", key)
}

func (d *DataStore) Put(val string) string {
	for {
		randKey := generateRandomKey(32)
		if err := d.set(randKey, val); err == nil {
			if err := d.save(randKey, val); err != nil {
				d.l.Println("Error saving to the data store : ", err)
			}

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
