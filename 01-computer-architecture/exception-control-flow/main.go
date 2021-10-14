package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
)

var cFlag = flag.String("c", "", "Run the shell as a single command instead of a repl")

func main() {
  flag.Parse()

  if cFlag != nil && *cFlag != "" {
    run(*cFlag)
    return
  }

  reader := bufio.NewReader(os.Stdin)
  for {
    fmt.Print("🎃 ")
    text, err := reader.ReadString('\n')
    if err == io.EOF {
      break
    } else if err != nil {
      panic(err)
    }
    run(text)
  }
  fmt.Println("\nHave a spooky good time! 🧙🏻🐈‍⬛")
}

func run(cmd string) {
  fmt.Print(cmd)
}
