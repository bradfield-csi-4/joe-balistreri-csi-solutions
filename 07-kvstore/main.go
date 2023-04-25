package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const NEWLINE = '\n'
const FILENAME = "abc.db"

var r = regexp.MustCompile("^(get|GET|set|SET) ([a-zA-Z0-9]+)(?:=(.+))?$")

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

	scanner := bufio.NewScanner(os.Stdin)

	// TODO: individual replicas will need separate db files
	db := loadDb(FILENAME)

	fmt.Println("Welcome to the DB console!!")
	fmt.Print("> ")
	for scanner.Scan() {
		cmd, err := commandFromString(scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
			fmt.Print("> ")
			continue
		}

		switch cmd.op {
		case GET:
			fmt.Println(db[cmd.key])
			fmt.Print("> ")
			continue
		case SET:
			db[cmd.key] = cmd.value
			fmt.Printf("%s=%s\n", cmd.key, cmd.value)
			fmt.Print("> ")

			// TODO: may want to do this less often for improve performance - on exit?
			writeDb(FILENAME, db)
			continue
		}

		fmt.Print("> ")
	}
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
