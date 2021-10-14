package main

import (
  "bufio"
  "errors"
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

type resp struct {
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

  inputChan := make(chan resp)
  Loop:
    for {
      fmt.Print("🎃 ")
      go func() {
        text, err := bufio.NewReader(os.Stdin).ReadString('\n')
        inputChan <- resp{s: text, err: err}
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

  i, j := 0, 0
  for ; j < len(args); j++ {
    if args[j] == "&&" {
      if i == j {
        fmt.Println("Invalid syntax")
        return
      }
      status := handleSingleCommand(args[i:j])
      if status != 0 {
        return
      }
      i = j + 1
    }
  }
  if j > i {
    handleSingleCommand(args[i:j])
  }
}

func handleSingleCommand(args []string) int {
  switch args[0] {
  case "exit":
    os.Exit(0)
    return 0
  case "cd":
    err := os.Chdir(strings.Join(args[1:], "/"))
    if err != nil {
      fmt.Println(err)
      return 1
    }
    return 0
  }

  var cmd *exec.Cmd
  if len(args) > 1 {
    cmd = exec.Command(args[0], args[1:]...)
  } else {
    cmd = exec.Command(args[0])
  }
  resChan := make(chan resp)
  go func() {
    output, err := cmd.CombinedOutput()
    if err != nil {
      if e, ok := err.(*exec.Error); ok {
          resChan <- resp{s: fmt.Sprintf("💀 %s: command not found\n", e.Name), err: errors.New("fail")}
        } else {
          resChan <- resp{s: string(output)}
        }
    } else {
      resChan <- resp{s: string(output)}
    }
  }()
  strResp, ok := listenToChannels(resChan, c)
  fmt.Print(strResp)
  if strResp == "" {
    fmt.Println()
  }
  if ok {
    return 0
  }
  return 1
}

func listenToChannels(resChan chan resp, c chan os.Signal) (string, bool) {
  select {
    case r := <-resChan:
      if r.err != nil {
        return r.s, false
      }
      return r.s, true
    case <- c:
      return "", false
  }
}
