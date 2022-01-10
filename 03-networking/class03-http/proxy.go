package main

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

func main() {
  fmt.Println(BANNER)
  proxyFd := tcpSocket()
  defer func() {
    syscall.Close(proxyFd)
  }()

  bind(proxyFd, port)

  fmt.Printf("server is now listening on port %d\n", port)
  err := syscall.Listen(proxyFd, 20)
  if err != nil {
    log.Fatalf("%v", err)
  }

  for {
    listenAndServe(proxyFd)
  }
}

func listenAndServe(proxyFd int) {
  clientFd, sa, err := syscall.Accept(proxyFd)
  if err != nil {
    log.Fatalf("%v", err)
  }
  fmt.Println("received an incoming connection!")
  data := make([]byte, 4096)
  n, _, err := syscall.Recvfrom(clientFd, data, 0)
  data = data[:n]
  if err != nil {
    log.Fatalf("%v", err)
  }

  req := parseRequest(data)

  var dataToReturn []byte
  cachedData, shouldCache, ok := cachedResponse(req)
  if ok {
    fmt.Println("using cached response!")
    dataToReturn = cachedData
  } else {
    dataToReturn = callDestination(dstPort, data)
    if shouldCache {
      pageCache[req.url] = dataToReturn
    }
  }

  syscall.Sendmsg(clientFd, dataToReturn, nil, sa, 0)
  syscall.Close(clientFd)
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

func parseRequest(data []byte) httpRequestMessage {
  messagePieces := strings.Split(string(data), "\r\n\r\n")
  headers := strings.Split(messagePieces[0], "\r\n")
  var body string
  if len(messagePieces) > 1 {
    body = messagePieces[1]
  }
  requestLinePieces := strings.Split(headers[0], " ")
  if len(requestLinePieces) != 3 {
    log.Fatalf("wrong number of request line pieces, %v", requestLinePieces)
  }
  requestLine := requestLine{
    method: requestLinePieces[0],
    url: requestLinePieces[1],
    version: requestLinePieces[2],
  }

  headers = headers[1:]
  headerMap := map[string]string{}
  for _, header := range headers {
    headerPieces := strings.Split(header, ": ")
    if len(headerPieces) != 2 {
      continue
    }
    headerMap[headerPieces[0]] = headerPieces[1]
  }

  return httpRequestMessage{
    body: []byte(body),
    headers: headerMap,
    requestLine: requestLine,
  }
}

type httpRequestMessage struct {
  requestLine
  headers map[string]string
  body []byte
}

type requestLine struct {
  method string
  url string
  version string
}

func callDestination(port int, data []byte) []byte {
  sa := &syscall.SockaddrInet4{Port: port}
  fd := tcpSocket()
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

func tcpSocket() int {
  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
  if err != nil {
    log.Fatalf("%v", err)
  }
  return fd
}

func bind(fd, port int) {
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
