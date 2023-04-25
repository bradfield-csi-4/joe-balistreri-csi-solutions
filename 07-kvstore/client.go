package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const sockAddr = "/tmp/db.sock"

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the DB console!!")
	fmt.Print("> ")
	for scanner.Scan() {
		conn, err := net.Dial("unix", sockAddr)
		if err != nil {
			panic(err)
		}

		_, err = conn.Write(scanner.Bytes())
		if err != nil {
			panic(err)
		}

		b := make([]byte, 4096)
		if n, err := conn.Read(b); err != nil {
			if err == io.EOF {
				fmt.Println()
			} else {
				panic(err)
			}
		} else {
			fmt.Println(string(b[:n]))
		}

		fmt.Print("> ")
	}
}
