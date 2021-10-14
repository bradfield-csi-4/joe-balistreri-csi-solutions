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
  "syscall"
)

var cFlag = flag.String("c", "", "Run the shell as a single command instead of a repl")
var c = make(chan os.Signal, 1)

func main() {
  signal.Notify(c)
  flag.Parse()

  if cFlag != nil && *cFlag != "" {
    run(*cFlag)
    return
  }

  // TODO: have this read happen in a channel so we can listen for signals at the same time
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
  listenToChannels(resChan, c)
}

func listenToChannels(resChan chan string, c chan os.Signal) {
  select {
    case r := <-resChan:
      fmt.Print(r)
    case s := <- c:
      switch s {
      case syscall.SIGURG:
        fmt.Println("ignoring sigurg")
        listenToChannels(resChan, c)
      default:
        fmt.Println("Got signal:", s)
      }
  }
}
