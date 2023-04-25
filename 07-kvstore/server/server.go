package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

const NEWLINE = '\n'
const FILENAME = "abc.db"
const sockAddr = "/tmp/db.sock"

var r = regexp.MustCompile("^(get|GET|set|SET) ([a-zA-Z0-9]+)(?:=(.+))?$")

// TODO: make this DB concurrency safe!!!
type DB map[string]string

type OP int

const (
	GET OP = iota + 1
	SET
)

type Command struct {
	op    OP
	key   string
	value string
}

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

			output := handleCommand(db, string(buf[:n]))

			_, err = conn.Write([]byte(output))
			if err != nil {
				panic(err)
			}
		}(conn)
	}
}

func handleCommand(db DB, command string) string {
	cmd, err := commandFromString(command)
	if err != nil {
		return err.Error()
	}

	switch cmd.op {
	case GET:
		return db[cmd.key]
	case SET:
		db[cmd.key] = cmd.value
		// TODO: may want to do this less often for improve performance - on exit?
		writeDb(FILENAME, db)

		return fmt.Sprintf("%s=%s", cmd.key, cmd.value)
	}

	return "unexpected error"
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

func commandFromString(s string) (*Command, error) {
	parts := r.FindStringSubmatch(s)
	if len(parts) != 4 {
		return nil, errors.New("invalid input")
	}

	cmd := Command{
		key: parts[2],
	}

	op := strings.ToLower(parts[1])
	switch op {
	case "get":
		cmd.op = GET
		if parts[3] != "" {
			return nil, errors.New("cannot set a value with GET")
		}
	case "set":
		cmd.op = SET
		cmd.value = parts[3]
	default:
		return nil, fmt.Errorf("invalid command %s", op)
	}

	return &cmd, nil
}
