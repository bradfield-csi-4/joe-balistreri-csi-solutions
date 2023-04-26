package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/jdbalistreri/bradfield-csi-solutions/07-kvstore/encoding"
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

		cmdString := scanner.Bytes()
		cmd, err := encoding.CommandFromString(string(cmdString))
		if err != nil {
			fmt.Println(err)
			fmt.Print("> ")
			continue
		}

		_, err = conn.Write(cmd.ToBinaryV1())
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
			respCmd, err := encoding.CommandFromBinary(b[:n])
			if err != nil {
				panic(err)
			}
			if respCmd.Op != encoding.MSG {
				panic("wrong op returned!")
			}
			fmt.Println(string(respCmd.Value))
		}

		fmt.Print("> ")
	}
}
