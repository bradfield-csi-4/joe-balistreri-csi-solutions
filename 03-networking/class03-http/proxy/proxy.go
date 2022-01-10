package proxy

import (
  "syscall"
  "fmt"
  "strings"
  "log"
)

const port = 8000
const dstPort = 9000


var cacheableRoutes = []string{
  "/website/",
}

var pageCache = map[string][]byte{}

func ListenAndServe(proxyFd int) {
  clientFd, sa, err := syscall.Accept(proxyFd)
  if err != nil {
    log.Fatalf("%v", err)
  }
  fmt.Println("received an incoming connection!")

  handleSingleRequest(clientFd, sa)
}

func handleSingleRequest(clientFd int, sa syscall.Sockaddr) {
  var httpBuf []byte // TODO: could use a bytes.Buffer, but it made checking for the double return harder

  // loop until we get two newlines, so we can receive the message in multiple pieces
  for {
    data := make([]byte, 4096)
    n, _, err := syscall.Recvfrom(clientFd, data, 0)
    if n == 0 {
      fmt.Println("client disconnected before sending a complete message")
      return
    }
    if err != nil {
      log.Fatalf("%v", err)
    }
    fmt.Println("received a message")
    data = data[:n]
    httpBuf = append(httpBuf, data...)
    if strings.Contains(string(httpBuf), DBL_BREAK) {
      break
    }
  }

  req := parseRequest(httpBuf)

  fmt.Printf("got request: %+v\n", req)

  var dataToReturn []byte
  cachedData, shouldCache, ok := cachedResponse(req)
  if ok {
    fmt.Println("using cached response!")
    dataToReturn = cachedData
  } else {
    dataToReturn = callDestination(dstPort, req.toHTTP())
    if shouldCache {
      pageCache[req.url] = dataToReturn
    }
  }

  syscall.Sendmsg(clientFd, dataToReturn, nil, sa, 0)

  if strings.ToLower(req.singleHopHeaders.connection) == "keep-alive" {
    // start a new goroutine to handle the next message on this connection
    go handleSingleRequest(clientFd, sa)
  } else {
    syscall.Close(clientFd)
  }
  return
}

func cachedResponse(req httpRequestMessage) (data []byte, shouldCache bool, ok bool) {
  if !cachingEnabled(req) {
    return nil, false, false
  }

  data, ok = pageCache[req.url]
  return data, true, ok
}

func cachingEnabled(req httpRequestMessage) (shouldCache bool) {
  for _, route := range cacheableRoutes {
    if strings.HasPrefix(req.url, route) {
      shouldCache = true
      break
    }
  }
  return shouldCache
}

func callDestination(port int, data []byte) []byte {
  sa := &syscall.SockaddrInet4{Port: port}
  fd := TcpSocket()
  err := syscall.Connect(fd, sa)
  if err != nil {
    log.Fatalf("failed to connect to destination. %v", err)
  }
  syscall.Sendmsg(fd, data, nil, sa, 0)
  fmt.Println("sent message to server")
  respData := make([]byte, 4096)
  _, _, err = syscall.Recvfrom(fd, respData, 0)
  if err != nil {
    log.Fatalf("failed to receive from destination. %v", err)
  }
  syscall.Close(fd)
  return respData
}

func TcpSocket() int {
  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
  if err != nil {
    log.Fatalf("%v", err)
  }
  return fd
}

func Bind(fd, port int) {
  err := syscall.Bind(fd, &syscall.SockaddrInet4{Port: port})
  if err != nil {
    log.Fatalf("%v", err)
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
