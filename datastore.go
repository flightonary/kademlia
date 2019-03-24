package kademlia

type dataStore struct {
	store map[string][]byte
}

func newDataStore() *dataStore {
	return &dataStore{make(map[string][]byte)}
}

func (ds *dataStore) Get(key string) []byte {
	return ds.store[key]
}

func (ds *dataStore) Put(key string, value []byte) {
	ds.store[key] = value
}

func (ds *dataStore) Delete(key string) {
	delete(ds.store, key)
}

func (ds *dataStore) Exist(key string) bool {
	_, ok := ds.store[key]
	return ok
}