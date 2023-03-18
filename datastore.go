package netclip

import "sync"

// DataStore holds our text data
type DataStore struct {
	data map[string]string
	mu   sync.Mutex
}

// Store saves a field to the datastore
func (ds *DataStore) Store(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = value
}

// Range lets us loop over all the records.
func (ds *DataStore) Range(f func(key, value string) bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for key, value := range ds.data {
		if !f(key, value) {
			break
		}
	}
}

// Delete removes a record
func (ds *DataStore) Delete(key string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	delete(ds.data, key)
}
