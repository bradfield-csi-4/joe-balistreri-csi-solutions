package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jdbalistreri/bradfield-csi-solutions/07-kvstore/encoding"
)

const NEWLINE = '\n'
const FILENAME = "abc.db"
const sockAddr = "/tmp/db.sock"

// TODO: make this DB concurrency safe!!!
type DB map[string][]byte

func main() {

	socket, err := net.Listen("unix", sockAddr)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(sockAddr)
		os.Exit(1)
	}()

	// TODO: individual replicas will need separate db files
	db := loadDb(FILENAME)

	for {
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 4096)

			n, err := conn.Read(buf)
			if err != nil {
				panic(err)
			}

			output := handleCommand(db, buf[:n])
			cmd := encoding.Command{
				Op:    encoding.MSG,
				Value: output,
			}

			_, err = conn.Write(cmd.ToBinaryV1())
			if err != nil {
				panic(err)
			}
		}(conn)
	}
}

func handleCommand(db DB, command []byte) []byte {
	cmd, err := encoding.CommandFromBinary(command)
	if err != nil {
		return []byte(err.Error())
	}

	switch cmd.Op {
	case encoding.GET:
		return db[string(cmd.Key)]
	case encoding.SET:
		db[string(cmd.Key)] = cmd.Value
		// TODO: may want to do this less often for improve performance - on exit?
		writeDb(FILENAME, db)

		return []byte(fmt.Sprintf("%s=%s", cmd.Key, cmd.Value))
	}

	return []byte("unexpected error")
}

func writeDb(filename string, db DB) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(f)
	return encoder.Encode(db)
}

func loadDb(filename string) DB {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to load from disk: %s\n", err.Error())
		fmt.Println("starting a fresh db")
		return DB{}
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("failed to read from file: %s\n", err.Error())
		fmt.Println("starting a fresh db")
		return DB{}
	}

	decoder := gob.NewDecoder(bytes.NewReader(b))
	db := DB{}
	err = decoder.Decode(&db)
	if err != nil {
		fmt.Printf("failed to decode db from file data: %s\n", err.Error())
		fmt.Println("starting a fresh db")
		return DB{}
	}

	return db
}
