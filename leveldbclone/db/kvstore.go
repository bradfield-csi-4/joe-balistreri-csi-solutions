package db

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
)

const MAX_MEMTABLE_SIZE_BYTES = 4096 // matches usual OS page size; randomly chosen
const SSTABLE_INDEX_INCR_BYTES = 400 // 10% of the SSTable size

func NewKVStore(name string) (DB, func()) {
	store := &KVStore{
		memtable:             NewSkipList(),
		name:                 name,
		maxMemtableSizeBytes: MAX_MEMTABLE_SIZE_BYTES,
	}
	var ssTable *SSTable
	if fileExists(store.ssTableFilename()) {
		f, err := os.Open(store.ssTableFilename())
		if err != nil {
			panic(err)
		}
		ssTable = LoadSSTable(f)
	}
	wal, wDone := NewWriteAheadLog(name + ".wal")
	store.wal = wal
	store.ssTable = ssTable
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
	ssTable              *SSTable
}

func (k *KVStore) Get(key []byte) ([]byte, error) {
	dbs := []ImmutableDB{k.memtable}
	if k.ssTable != nil {
		dbs = append(dbs, k.ssTable)
	}
	for _, db := range dbs {
		if db == nil {
			continue
		}
		v, err := db.Get(key)
		if err != nil {
			if _, ok := err.(*NotFoundError); ok {
				continue
			}
			return nil, err
		}
		return v, nil
	}
	return nil, &NotFoundError{}
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

func (k *KVStore) flushSSTable(f *os.File) error {
	var fileContents []byte

	i, err := k.memtable.RangeScan(nil, nil)
	if err != nil {
		return err
	}

	var index []IndexEntry
	var currBytes int
	var firstKey *[]byte
	for i.Next() {
		// append to file contents
		if firstKey == nil {
			k := i.Key()
			firstKey = &k
		}
		nextLine := toLogLine(i.Key(), i.Value())
		currBytes += len(nextLine)
		fileContents = append(fileContents, nextLine...)

		// append to index
		if currBytes > SSTABLE_INDEX_INCR_BYTES {
			index = append(index, IndexEntry{Key: *firstKey, Offset: 4 + len(fileContents) - currBytes, Length: currBytes})
			currBytes = 0
			firstKey = nil
		}
	}

	fmt.Printf("index is %+v\n", index)
	b := &bytes.Buffer{}
	gob.NewEncoder(b).Encode(index)

	indexOffset := len(fileContents) + 4
	fmt.Printf("index offset is %d\n", indexOffset)
	indexOffsetEncoded := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexOffsetEncoded, uint32(indexOffset))
	fileContents = append(indexOffsetEncoded, fileContents...)
	fileContents = append(fileContents, b.Bytes()...)

	// write to disc and flush
	n, err := f.Write(fileContents)
	if err != nil {
		return err
	}
	if n != len(fileContents) {
		return fmt.Errorf("wrote %d bytes but expected to write %d", n, len(fileContents))
	}
	err = f.Sync()
	if err == nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (k *KVStore) ssTableFilename() string {
	return k.name + ".sst"
}

func (k *KVStore) checkAndHandleFlush() error {
	// check the size of the memtable
	if k.memtable.SizeBytes() <= k.maxMemtableSizeBytes {
		return nil
	}
	fmt.Printf("memtable is %d bytes large, need to flush\n", k.memtable.SizeBytes())

	if fileExists(k.ssTableFilename()) {
		fmt.Println("SSTable already exists! Dropping the in memory values for now")
	} else {
		// create the SSTable file (and overwrite the old one - will fix that later)
		f, err := os.OpenFile(k.ssTableFilename(), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}

		// call flush on the memtable
		err = k.flushSSTable(f)
		if err != nil {
			return err
		}
		f.Seek(0, 0)
		k.ssTable = LoadSSTable(f)
	}

	// clear the writeahead log
	err := k.wal.Truncate()
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
