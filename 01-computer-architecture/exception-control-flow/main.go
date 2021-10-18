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
  var allCommands []*exec.Cmd
  for {
    fmt.Print("🎃 ")
    command, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err == io.EOF {
        break
    } else if err != nil {
      panic(err)
    }
    allCommands = append(allCommands, run(command)...)
  }
  for _, cmd := range allCommands {
    if cmd != nil && cmd.Process != nil {
      cmd.Process.Kill()
    }
  }
  fmt.Println(ExitMessage)
}

type ArgumentGroup struct {
  args []string
  continueOnSuccess bool
  continueOnFailure bool
  backgroundTask bool
}

func run(command string) []*exec.Cmd {
  split := strings.Split(strings.TrimSuffix(command, "\n"), " ")
  if len(split) == 0 {
    fmt.Println()
    return nil
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
      return nil
    }
    argumentGroups = append(argumentGroups, ArgumentGroup{
      args: append(split[i:j], additionalArgs...),
      continueOnSuccess: token == "&&" || token == ";" || token == "&",
      continueOnFailure: token == "||" || token == ";" || token == "&",
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
  var commandsCreated []*exec.Cmd
  for _, group := range argumentGroups {
    cmd, exitStatus := runArgumentGroup(group)
    commandsCreated = append(commandsCreated, cmd)

    if exitStatus == 0 && !group.continueOnSuccess {
      break
    }
    if exitStatus != 0 && !group.continueOnFailure {
      break
    }
  }
  return commandsCreated
}

func runArgumentGroup(argGroup ArgumentGroup) (*exec.Cmd, int) {
  switch argGroup.args[0] {
  case "exit":
    fmt.Println(ExitMessage)
    os.Exit(0)
    return nil, 0
  case "cd":
    err := os.Chdir(strings.Join(argGroup.args[1:], "/"))
    if err != nil {
      fmt.Println(err)
      return nil, 1
    }
    return nil, 0
  }

  var cmd *exec.Cmd
  if len(argGroup.args) > 1 {
    cmd = exec.Command(argGroup.args[0], argGroup.args[1:]...)
  } else {
    cmd = exec.Command(argGroup.args[0])
  }

  return cmd, runKillableCommand(cmd, argGroup.backgroundTask)
}

func runKillableCommand(cmd *exec.Cmd, background bool) int {
  doneChannel := make(chan bool)

  // initiate command in a goroutine
  go func() {
    if background {
      cmd.SysProcAttr = &syscall.SysProcAttr{
        Setpgid: true,
      }
    }
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
      if e, ok := err.(*exec.Error); ok {
          fmt.Fprintf(os.Stdout, "💀 %s: command not found\n", e.Name)
      }
    }
    if !background {
      doneChannel<-true
    }
  }()

  if background {
    return 0
  }

  // complete command or kill if interrupt signal received
  select {
  case <-doneChannel:
    if cmd.Process == nil || cmd.ProcessState == nil {
      return 1
    }
    return cmd.ProcessState.ExitCode()
  case <-interruptChannel:
    if cmd != nil && cmd.Process != nil {
      cmd.Process.Kill()
    }
    fmt.Println()
    return 0
  }
}
