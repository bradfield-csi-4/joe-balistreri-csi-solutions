package main

import (
  "fmt"
  "syscall"
  "strconv"
  "log"
  "shared"
  "os"
)

const receiverPort = 9000

func main() {
  proxyPort, _ := strconv.Atoi(os.Args[1])
  fmt.Printf("starting receiver on port %d, proxying to port %d\n", receiverPort, proxyPort)

  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
  if err != nil {
    log.Fatalf("failed to create socket %v", err)
  }
  // allow the socket to be reused so we can immediately restart the server
  syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

  syscall.Bind(fd, &syscall.SockaddrInet4{Port: receiverPort})

  for {
    request := make([]byte, 65536)
    n, _, err := syscall.Recvfrom(fd, request, 0)
    if err != nil {
      panic(err)
    }
    request = request[:n]

    requestHeader := shared.HeaderFromBytes(request)

    fmt.Printf("got request: %+v\n", requestHeader)
    if shared.SumBytes(requestHeader) != shared.ValidSum {
      fmt.Println("got corrupted packet!")
    }

    payload := []byte(string(requestHeader.Data) + string(requestHeader.Data))
    payloadWithHeader := shared.HeaderToBytes(shared.NewHeader(payload))

    syscall.Sendto(fd, payloadWithHeader, 0, &syscall.SockaddrInet4{Port: proxyPort})
  }
}
