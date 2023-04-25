package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs)
	// done := make(chan bool, 1)

	for {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
	}

	// 	done <- true
	// }()

	// fmt.Println("awaiting signal")
	// <-done
	// fmt.Println("exiting")
}
