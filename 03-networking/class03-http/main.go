package main

import (
  "syscall"
  "fmt"
  "os"
  "os/signal"
  "log"
  "github.com/jbalistreri/bradfield-csi-solutions/03-networking/class03-http/proxy"
)

const port = 8000
const dstPort = 9000

func main() {
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  fmt.Println(BANNER)
  proxyFd := proxy.TcpSocket()
  // allow the socket to be reused so we can immediately restart the server
  syscall.SetsockoptInt(proxyFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
  go func() {
    <-sigs
    fmt.Println("gracefully shutting down")
    syscall.Close(proxyFd)
    os.Exit(0)
  }()

  proxy.Bind(proxyFd, port)

  fmt.Printf("server is now listening on port %d\n", port)
  err := syscall.Listen(proxyFd, 20)
  if err != nil {
    log.Fatalf("%v", err)
  }

  for {
    proxy.ListenAndServe(proxyFd)
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
