package db

func NewKVStore(name string) DB {
	store := &KVStore{
		memtable: NewSkipList(),
		name:     name,
		wal:      NewWriteAheadLog(name + ".wal"),
	}
	i, err := store.wal.Iterator()
	if err != nil {
		panic(err)
	}
	for i.Next() {
		k := i.Key()
		v := i.Value()
		if v == nil {
			store.memtable.Delete(k)
		} else {

			store.memtable.Put(k, v)
		}
	}
	return store
}

type KVStore struct {
	memtable DB
	wal      WriteAheadLog
	name     string
}

func (k *KVStore) Get(key []byte) (value []byte, err error) {
	return k.memtable.Get(key)
}

func (k *KVStore) Has(key []byte) (ret bool, err error) {
	return k.memtable.Has(key)
}

func (k *KVStore) Delete(key []byte) error {
	err := k.wal.Write(key, nil)
	if err != nil {
		return err
	}
	return k.memtable.Delete(key)
}

func (k *KVStore) Put(key, value []byte) error {
	err := k.wal.Write(key, value)
	if err != nil {
		return err
	}
	return k.memtable.Put(key, value)
}

func (k *KVStore) RangeScan(start, limit []byte) (Iterator, error) {
	return k.memtable.RangeScan(start, limit)
}
