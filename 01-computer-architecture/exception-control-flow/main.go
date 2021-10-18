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

type ArgumentGroup struct {
  args []string
  continueOnSuccess bool
  continueOnFailure bool
  backgroundTask bool
}

func run(command string) {
  split := strings.Split(strings.TrimSuffix(command, "\n"), " ")
  if len(split) == 0 {
    fmt.Println()
    return
  }

  var argumentGroups []ArgumentGroup
  i, j := 0, 0
  for ; j < len(split); j++ {
    token := split[j]
    var additionalArgs []string
    if len(token) > 0 && token[len(token)-1] == ';' {
      additionalArgs = []string{token[:len(token)-1]}
      token = ";"
    }
    if len(token) == 0 || !(token == "&&" || token == "||" || token == ";" || token == "&")  {
      continue
    }
    if i == j && len(additionalArgs) == 0 {
      fmt.Println("Invalid syntax")
      return
    }
    argumentGroups = append(argumentGroups, ArgumentGroup{
      args: append(split[i:j], additionalArgs...),
      continueOnSuccess: token == "&&" || token == ";",
      continueOnFailure: token == "||" || token == ";",
      backgroundTask: token == "&",
    })
    i = j + 1
  }
  if j > i {
    argumentGroups = append(argumentGroups, ArgumentGroup{
      args: split[i:j],
    })
  }

  // run each argument group
  for _, group := range argumentGroups {
    exitStatus := runArgumentGroup(group)
    if exitStatus == 0 && !group.continueOnSuccess {
      break
    }
    if exitStatus != 0 && !group.continueOnFailure {
      break
    }
  }
}

func runArgumentGroup(argGroup ArgumentGroup) int {
  switch argGroup.args[0] {
  case "exit":
    fmt.Println(ExitMessage)
    os.Exit(0)
    return 0
  case "cd":
    err := os.Chdir(strings.Join(argGroup.args[1:], "/"))
    if err != nil {
      fmt.Println(err)
      return 1
    }
    return 0
  }

  var cmd *exec.Cmd
  if len(argGroup.args) > 1 {
    cmd = exec.Command(argGroup.args[0], argGroup.args[1:]...)
  } else {
    cmd = exec.Command(argGroup.args[0])
  }

  return runKillableCommand(cmd)
}

func runKillableCommand(cmd *exec.Cmd) int {
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
    if cmd.Process == nil {
      return 1
    }
    state, err := cmd.Process.Wait()
    if err != nil || state == nil {
      return 1
    }
    return state.ExitCode()
  case <-interruptChannel:
    if cmd != nil && cmd.Process != nil {
      cmd.Process.Kill()
    }
    fmt.Println()
    return 0
  }
}
