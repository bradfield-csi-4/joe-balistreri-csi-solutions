package main

import (
  "fmt"
  "log"
  "os"
  "strconv"
  "syscall"
  "time"
  "shared"
)

const senderPort = 8000

func main() {
  proxyPort, _ := strconv.Atoi(os.Args[1])
  fmt.Printf("starting sender on port %d, proxying to port %d\n", senderPort, proxyPort)

  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
  if err != nil {
    log.Fatalf("failed to create socket %v", err)
  }
  // allow the socket to be reused so we can immediately restart the server
  syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

  syscall.Bind(fd, &syscall.SockaddrInet4{Port: senderPort})

  i := 1
  for {
    payload := []byte(fmt.Sprintf("Hello, it's me! %d", i))

    payloadWithHeader := shared.HeaderToBytes(shared.NewHeader(payload))

    syscall.Sendto(fd, payloadWithHeader, 0, &syscall.SockaddrInet4{Port: proxyPort})

    response := make([]byte, 65536)
    n, _, err := syscall.Recvfrom(fd, response, 0)
    if err != nil {
      panic(err)
    }
    response = response[:n]

    responseHeader := shared.HeaderFromBytes(response)
    if shared.SumBytes(responseHeader) != shared.ValidSum {
      fmt.Println("got corrupted packet!")
    }

    fmt.Printf("got response: %+v\n", responseHeader)
    i++

    time.Sleep(time.Second)
  }
}
