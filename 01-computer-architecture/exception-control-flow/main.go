package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
  "os/exec"
  // "reflect"
  "strings"
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

func run(text string) {
  text = strings.TrimSuffix(text, "\n")
  args := strings.Split(text, " ")

  switch args[0] {
  case "exit":
    os.Exit(0)
    return
  case "cd":
    os.Chdir(strings.Join(args[1:], "/"))
    return
  }

  var cmd *exec.Cmd
  if len(args) > 1 {
    cmd = exec.Command(args[0], args[1:]...)
  } else {
    cmd = exec.Command(args[0])
  }
  output, err := cmd.CombinedOutput()
  if err != nil {
    if e, ok := err.(*exec.Error); ok {
      fmt.Printf("💀 %s: command not found\n", e.Name)
    } else {
      fmt.Print(string(output))
    }
  } else {
    fmt.Print(string(output))
  }
}
