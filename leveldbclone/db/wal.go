package db

import (
	"fmt"
	"os"
)

func NewWriteAheadLog(filename string) (WriteAheadLog, func()) {
	w := WriteAheadLog{
		filename: filename,
	}
	return w, w.Close
}

type WriteAheadLog struct {
	filename string
}

func (w *WriteAheadLog) Close() {}

func (w *WriteAheadLog) Truncate() error {
	// TODO: make sure we don't delete the file while there's a currently open iterator
	return os.Remove(w.filename)
}

func (w *WriteAheadLog) Write(key, value []byte) error {
	f, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	logLine := toLogLine(key, value)
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

	n, key, val, err := readLogLine(i.file)
	i.position += n
	i.currKey = key
	i.currValue = val
	i.currErr = err
	return err == nil
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
