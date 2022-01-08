package main

import (
  "syscall"
  "fmt"
)

const port = 8888

func main() {
  fmt.Println(BANNER)
  fmt.Println("connecting to socket...")
  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
  if err != nil {
    panic(err)
  }
  defer func() {
    syscall.Close(fd)
  }()

  fmt.Printf("binding to port %d...\n", port)
  err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: port})
  if err != nil {
    panic(err)
  }

  fmt.Println("server is now listening")
  err = syscall.Listen(fd, 20)
  if err != nil {
    panic(err)
  }

  for {
    nfd, sa, err := syscall.Accept(fd)
    if err != nil {
      panic(err)
    }
    fmt.Println("received an incoming connection!")
    data := make([]byte, 4096)
    n, _, err := syscall.Recvfrom(nfd, data, 0)
    if err != nil {
      panic(err)
    }
    fmt.Printf("read %d bytes\n", n)

    syscall.Sendmsg(nfd, data, nil, sa, 0)
    syscall.Close(nfd)
  }
}


const BANNER = `  888888                   8888888b.
    "88b                   888   Y88b
     888                   888    888
     888  .d88b.   .d88b.  888   d88P 888d888 .d88b.  888  888 888  888
     888 d88""88b d8P  Y8b 8888888P"  888P"  d88""88b 'Y8bd8P' 888  888
     888 888  888 88888888 888        888    888  888   X88K   888  888
     88P Y88..88P Y8b.     888        888    Y88..88P .d8""8b. Y88b 888
     888  "Y88P"   "Y8888  888        888     "Y88P"  888  888  "Y88888
   .d88P                                                            888
 .d88P"                                                        Y8b d88P
888P"                                                           "Y88P"


`
