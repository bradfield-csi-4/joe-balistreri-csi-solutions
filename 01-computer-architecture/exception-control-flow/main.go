package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
  "os/exec"
  "strings"
)

var cFlag = flag.String("c", "", "Run the shell as a single command instead of a repl")

func main() {
  // allow the program to run a single command passed by a flag
  flag.Parse()
  if cFlag != nil && *cFlag != "" {
    run(*cFlag)
    return
  }

  // main loop
  for {
    fmt.Print("🎃 ")
    command, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err == io.EOF {
        break
    } else if err != nil {
      panic(err)
    }
    run(command)
  }
  fmt.Println("\nHave a spooky good time! 🧙🏻🐈‍⬛")
}

func run(command string) {
  split := strings.Split(strings.TrimSuffix(command, "\n"), " ")
  if len(split) == 0 {
    fmt.Println()
    return
  }

  switch split[0] {
  case "exit":
    os.Exit(0)
    return
  case "cd":
    err := os.Chdir(strings.Join(split[1:], "/"))
    if err != nil {
      fmt.Println(err)
      return
    }
    return
  }

  var cmd *exec.Cmd
  if len(split) > 1 {
    cmd = exec.Command(split[0], split[1:]...)
  } else {
    cmd = exec.Command(split[0])
  }
  output, err := cmd.CombinedOutput()
  if err != nil {
    if e, ok := err.(*exec.Error); ok {
        fmt.Printf("💀 %s: command not found\n", e.Name)
      } else {
        fmt.Print(string(output))
      }
  }
  fmt.Print(string(output))
}
