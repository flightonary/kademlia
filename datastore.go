package kademlia

type dataStore struct {
	store map[string][]byte
}

func newDataStore() *dataStore {
	return &dataStore{make(map[string][]byte)}
}

func (ds *dataStore) Get(kid *KadID) []byte {
	return ds.store[String(kid)]
}

func (ds *dataStore) GetAsString(kid *KadID) string {
	return string(ds.Get(kid))
}

func (ds *dataStore) Put(kid *KadID, value []byte) {
	ds.store[String(kid)] = value
}

func (ds *dataStore) Delete(kid *KadID) {
	delete(ds.store, String(kid))
}
