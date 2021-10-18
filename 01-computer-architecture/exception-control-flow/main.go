package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
  "os/exec"
  "os/signal"
  "strings"
)


var ExitMessage = "\nHave a spooky good time! 🧙🏻🐈‍⬛"
var cFlag = flag.String("c", "", "Run the shell as a single command instead of a repl")
var interruptChannel = make(chan os.Signal, 1)

func main() {
  signal.Notify(interruptChannel, os.Interrupt)
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
  fmt.Println(ExitMessage)
}

func run(command string) {
  split := strings.Split(strings.TrimSuffix(command, "\n"), " ")
  if len(split) == 0 {
    fmt.Println()
    return
  }

  switch split[0] {
  case "exit":
    fmt.Println(ExitMessage)
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
  runKillableCommand(cmd)
}

func runKillableCommand(cmd *exec.Cmd) {
  outputChannel := make(chan string)

  // initiate command in a goroutine
  go func() {
    output, err := cmd.CombinedOutput()
    if err != nil {
      if e, ok := err.(*exec.Error); ok {
          outputChannel<-fmt.Sprintf("💀 %s: command not found\n", e.Name)
      } else {
        outputChannel<-string(output)
      }
    } else {
      outputChannel<-string(output)
    }
  }()

  // complete command or kill if interrupt signal received
  select {
  case result := <-outputChannel:
    fmt.Print(result)
  case <-interruptChannel:
    if cmd != nil && cmd.Process != nil {
      cmd.Process.Kill()
    }
    fmt.Println()
  }
}
