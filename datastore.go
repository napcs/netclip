package netclip

import (
	"sort"
	"sync"
)

// DataStore holds our text data
type DataStore struct {
	data map[string]string
	mu   sync.Mutex
}

// NewDataStore initializes a new data store
func NewDataStore() DataStore {
	return DataStore{
		data: make(map[string]string),
	}
}

// Store saves a field to the datastore
func (ds *DataStore) Store(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = value
}

// Range lets us loop over all the records.
func (ds *DataStore) Range() map[string]string {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	keys := make([]string, 0, len(ds.data))
	for key := range ds.data {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	sortedData := make(map[string]string, len(ds.data))
	for _, key := range keys {
		sortedData[key] = ds.data[key]
	}

	return sortedData
}

// Delete removes a record
func (ds *DataStore) Delete(key string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	delete(ds.data, key)
}

// GetValue gets the value at the key
func (ds *DataStore) GetValue(key string) (string, bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	value, ok := ds.data[key]
	return value, ok
}
