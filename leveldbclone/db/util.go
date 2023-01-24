package db

func KeyFromIterator(i int) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}
