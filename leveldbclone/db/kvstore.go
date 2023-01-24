package db

import (
	"fmt"
	"os"
)

const MAX_MEMTABLE_SIZE_BYTES = 4096 // matches usual OS page size; randomly chosen

func NewKVStore(name string) (DB, func()) {
	wal, wDone := NewWriteAheadLog(name + ".wal")
	store := &KVStore{
		memtable:             NewSkipList(),
		name:                 name,
		wal:                  wal,
		maxMemtableSizeBytes: MAX_MEMTABLE_SIZE_BYTES,
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
	return store, func() {
		wDone()
		store.Close()
	}
}

type KVStore struct {
	memtable             *SkipList
	wal                  WriteAheadLog
	name                 string
	maxMemtableSizeBytes int
	ssTable              *os.File
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
	err = k.memtable.Delete(key)
	if err != nil {
		return err
	}
	return k.checkAndHandleFlush()
}

func (k *KVStore) Put(key, value []byte) error {
	err := k.wal.Write(key, value)
	if err != nil {
		return err
	}
	err = k.memtable.Put(key, value)
	if err != nil {
		return err
	}
	return k.checkAndHandleFlush()
}

func (k *KVStore) checkAndHandleFlush() error {
	// check the size of the memtable
	if k.memtable.SizeBytes() <= k.maxMemtableSizeBytes {
		return nil
	}
	fmt.Printf("memtable is %d bytes large, need to flush\n", k.memtable.SizeBytes())

	// create the SSTable file (and overwrite the old one - will fix that later)
	f, err := os.OpenFile(k.name+".sst", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	k.ssTable = f

	// call flush on the memtable
	err = k.memtable.flushSSTable(f)
	if err != nil {
		return err
	}

	// clear the writeahead log
	err = k.wal.Truncate()
	if err != nil {
		return err
	}

	// new memtable
	k.memtable = NewSkipList()
	return nil
}

func (k *KVStore) RangeScan(start, limit []byte) (Iterator, error) {
	return k.memtable.RangeScan(start, limit)
}

func (k *KVStore) Close() {
	if k.ssTable != nil {
		k.ssTable.Close()
	}
}
