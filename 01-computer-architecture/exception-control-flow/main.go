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

var cFlag = flag.String("c", "", "Run the shell as a single command instead of a repl")
var c = make(chan os.Signal, 1)

type inputResp struct {
  s string
  err error
}

func main() {
  signal.Notify(c, os.Interrupt)
  flag.Parse()

  if cFlag != nil && *cFlag != "" {
    run(*cFlag)
    return
  }

  inputChan := make(chan inputResp)
  Loop:
    for {
      fmt.Print("🎃 ")
      go func() {
        text, err := bufio.NewReader(os.Stdin).ReadString('\n')
        inputChan <- inputResp{s: text, err: err}
      }()

      select {
      case input := <-inputChan:
        if input.err == io.EOF {
          break Loop
        } else if input.err != nil {
          panic(input.err)
        }
        run(input.s)
      case <-c:
        break Loop
      }
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
  resChan := make(chan string)
  go func() {
    output, err := cmd.CombinedOutput()
    if err != nil {
      if e, ok := err.(*exec.Error); ok {
        resChan <- fmt.Sprintf("💀 %s: command not found\n", e.Name)
        } else {
          resChan <- string(output)
        }
      } else {
        resChan <- string(output)
      }
  }()
  if resp, ok := listenToChannels(resChan, c); ok {
    fmt.Print(resp)
  } else {
    fmt.Println()
    // TODO: terminate child goroutine
  }
}

func listenToChannels(resChan chan string, c chan os.Signal) (string, bool) {
  select {
    case r := <-resChan:
      return r, true
    case <- c:
      return "", false
  }
}
