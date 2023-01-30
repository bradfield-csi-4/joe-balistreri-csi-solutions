package db

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
)

const MAX_MEMTABLE_SIZE_BYTES = 4096 // matches usual OS page size; randomly chosen
const SSTABLE_INDEX_INCR_BYTES = 400 // 10% of the SSTable size

func NewKVStore(name string) (DB, func()) {
	// load stored ssTable metadata

	metadataFile := metadataFilename(name)
	metadata := KVStoreMetadata{}
	if fileExists(metadataFile) {
		f, err := os.Open(metadataFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		err = gob.NewDecoder(bytes.NewReader(b)).Decode(&metadata)
		if err != nil {
			panic(err)
		}
	}

	// make the store
	store := &KVStore{
		memtable:        NewSkipList(),
		name:            name,
		KVStoreMetadata: metadata,
	}

	// load all SSTables
	var ssTable *SSTable
	for i := store.KVStoreMetadata.SSTableIncr; i >= 0; i-- {
		if fileExists(store.ssTableFilename(i)) {
			f, err := os.Open(store.ssTableFilename(i))
			if err != nil {
				panic(err)
			}
			ssTable = LoadSSTable(f)
			store.SSTables = append(store.SSTables, ssTable)
		}
	}

	// load any values from the write-ahead log into the memtable
	wal, wDone := NewWriteAheadLog(name + ".wal")
	store.wal = wal
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

	// return
	return store, func() {
		wDone()
		store.Close()
	}
}

func metadataFilename(storeName string) string {
	return storeName + ".metadata"
}

type KVStore struct {
	memtable *SkipList
	wal      WriteAheadLog
	name     string
	KVStoreMetadata
	SSTables []*SSTable
}

type KVStoreMetadata struct {
	SSTableIncr int
}

func (k *KVStore) Get(key []byte) ([]byte, error) {
	for _, db := range k.dbs() {
		v, err := db.Get(key)
		if err != nil {
			if err == ErrNotFound {
				continue
			}
			return nil, err // Deletes handled by ErrKeyDeleted
		}
		return v, nil
	}
	return nil, ErrNotFound
}

func (k *KVStore) Has(key []byte) (ret bool, err error) {
	for _, db := range k.dbs() {
		v, err := db.Has(key)
		if err != nil {
			if err == ErrNotFound {
				continue
			}
			return false, err
		}
		return v, nil
	}
	return false, ErrNotFound
}

func (k *KVStore) dbs() []ImmutableDB {
	dbs := []ImmutableDB{k.memtable}
	for _, s := range k.SSTables {
		dbs = append(dbs, s)
	}
	return dbs
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

func (k *KVStore) flushSSTable(f *os.File, it Iterator) error {
	var fileContents []byte

	var index []IndexEntry
	var currBytes int
	var firstKey *[]byte
	for it.Next() {
		// append to file contents
		if firstKey == nil {
			k := it.Key()
			firstKey = &k
		}
		nextLine := toLogLine(it.Key(), it.Value())
		currBytes += len(nextLine)
		fileContents = append(fileContents, nextLine...)

		// append to index
		if currBytes > SSTABLE_INDEX_INCR_BYTES {
			index = append(index, IndexEntry{Key: *firstKey, Offset: 4 + len(fileContents) - currBytes, Length: currBytes})
			currBytes = 0
			firstKey = nil
		}
	}
	if currBytes > 0 {
		index = append(index, IndexEntry{Key: *firstKey, Offset: 4 + len(fileContents) - currBytes, Length: currBytes})
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

func (k *KVStore) ssTableFilename(i int) string {
	return k.name + strconv.Itoa(i) + ".sst"
}

func (k *KVStore) writeMetadata() error {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(k.KVStoreMetadata)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(metadataFilename(k.name), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b.Bytes())
	return err
}

func (k *KVStore) checkAndHandleFlush() error {
	// check the size of the memtable
	if k.memtable.SizeBytes() <= MAX_MEMTABLE_SIZE_BYTES {
		return nil
	}
	fmt.Printf("memtable is %d bytes large, need to flush\n", k.memtable.SizeBytes())

	// create the SSTable file (and overwrite the old one - will fix that later)
	f, err := os.OpenFile(k.ssTableFilename(k.KVStoreMetadata.SSTableIncr+1), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	// call flush on the memtable
	it, err := k.memtable.RangeScan(nil, nil)
	if err != nil {
		return err
	}

	// add the new SSTable to our list
	err = k.flushSSTable(f, it)
	if err != nil {
		return err
	}
	f.Seek(0, 0)
	k.SSTables = append([]*SSTable{LoadSSTable(f)}, k.SSTables...)

	// update the kvstore metadata
	k.KVStoreMetadata.SSTableIncr += 1
	err = k.writeMetadata()
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
	for _, s := range k.SSTables {
		s.Close()
	}
}
