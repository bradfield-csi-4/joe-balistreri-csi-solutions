package db

import (
	"encoding/binary"
	"io"
)

func KeyFromIterator(i int) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func toLogLine(key, value []byte) []byte {
	var logLine []byte

	keySize := make([]byte, 4)
	valueSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(keySize, uint32(len(key)))
	binary.LittleEndian.PutUint32(valueSize, uint32(len(value)))

	logLine = append(logLine, keySize...)
	logLine = append(logLine, valueSize...)
	logLine = append(logLine, key...)
	logLine = append(logLine, value...)

	return logLine
}

func readLogLine(f io.Reader) (readBytes int, key, value []byte, err error) {
	// read the sizes of next key and value
	sizes := make([]byte, 8)
	n, err := f.Read(sizes)
	if err != nil {
		return n, nil, nil, err
	}
	readBytes += n
	keySize := int(binary.LittleEndian.Uint32(sizes[:4]))
	valueSize := int(binary.LittleEndian.Uint32(sizes[4:]))

	// read the next key and value
	keyAndValue := make([]byte, keySize+valueSize)
	n, err = f.Read(keyAndValue)
	if err != nil {
		return readBytes + n, nil, nil, err
	}
	readBytes += n
	return readBytes, keyAndValue[:keySize], keyAndValue[keySize:], nil
}
