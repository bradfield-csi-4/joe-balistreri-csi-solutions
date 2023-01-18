package db

import (
	"encoding/binary"
	"fmt"
	"os"
)

func NewWriteAheadLog(filename string) WriteAheadLog {
	return WriteAheadLog{
		filename: filename,
	}
}

type WriteAheadLog struct {
	filename string
}

func (w *WriteAheadLog) Write(key, value []byte) error {
	f, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var logLine []byte

	keySize := make([]byte, 4)
	valueSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(keySize, uint32(len(key)))
	binary.LittleEndian.PutUint32(valueSize, uint32(len(value)))

	logLine = append(logLine, keySize...)
	logLine = append(logLine, valueSize...)
	logLine = append(logLine, key...)
	logLine = append(logLine, value...)

	n, err := f.Write(logLine)
	if err != nil {
		return err
	}
	if n != len(logLine) {
		return fmt.Errorf("wrote %d bytes but expected to write %d", n, len(logLine))
	}
	err = f.Sync()
	if err == nil {
		return err
	}
	return nil
}

func (k *WriteAheadLog) Iterator() (Iterator, error) {
	// this file remains open as long as the iterator is in the process of being read
	f, err := os.Open(k.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &WALIterator{}, nil
		}
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &WALIterator{
		maxPosition: int(s.Size()),
		file:        f,
	}, nil
}

func (i *WALIterator) Next() bool {
	// make sure we haven't exceeded the max size of the file at the time the iterator began
	if i.position >= i.maxPosition {
		if i.file != nil {
			i.file.Close()
			i.file = nil
		}
		i.currErr = nil
		i.currKey = nil
		i.currValue = nil
		return false
	}

	// read the sizes of next key and value
	sizes := make([]byte, 8)
	n, err := i.file.Read(sizes)
	if err != nil {
		i.currErr = err
		return false
	}
	i.position += n
	keySize := int(binary.LittleEndian.Uint32(sizes[:4]))
	valueSize := int(binary.LittleEndian.Uint32(sizes[4:]))

	// read the next key and value
	keyAndValue := make([]byte, keySize+valueSize)
	n, err = i.file.Read(keyAndValue)
	if err != nil {
		i.currErr = err
		return false
	}
	i.position += n
	i.currKey = keyAndValue[:keySize]
	i.currValue = keyAndValue[keySize:]
	return true
}

func (i *WALIterator) Error() error {
	return i.currErr
}

func (i *WALIterator) Key() []byte {
	return i.currKey
}

func (i *WALIterator) Value() []byte {
	return i.currValue
}

type WALIterator struct {
	maxPosition int
	position    int
	file        *os.File
	currKey     []byte
	currValue   []byte
	currErr     error
}
