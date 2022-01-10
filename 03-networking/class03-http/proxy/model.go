package proxy

import (
  "bytes"
  "strings"
  "log"
  "fmt"
)

const LINE_BREAK = "\r\n"
const DBL_BREAK = LINE_BREAK + LINE_BREAK
const SPACE = " "
const HEADER_SEP = ": "

func parseRequest(data []byte) httpRequestMessage {
  messagePieces := strings.Split(string(data), DBL_BREAK)
  headers := strings.Split(messagePieces[0], LINE_BREAK)
  var body string
  if len(messagePieces) > 1 {
    body = messagePieces[1]
  }
  requestLinePieces := strings.Split(headers[0], SPACE)
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
    headerPieces := strings.Split(header, HEADER_SEP)
    if len(headerPieces) != 2 {
      continue
    }
    headerMap[headerPieces[0]] = headerPieces[1]
  }

  singleHop := removeSingleHopHeaders(headerMap)

  return httpRequestMessage{
    body: []byte(body),
    headers: headerMap,
    singleHopHeaders: singleHop,
    requestLine: requestLine,
  }
}

func removeSingleHopHeaders(headerMap map[string]string) singleHopHeaders {
  result := singleHopHeaders{}

  for k, v := range headerMap {
    switch strings.ToLower(k) {
    case "connection":
      result.connection = v
    case "keep-alive", "transfer-encoding", "te", "trailer", "upgrade", "proxy-authorization", "proxy-authenticate":
      // fall through so the value is deleted
      // could add them to the singleHopHeaders struct as a future improvement
    default:
      continue
    }
    delete(headerMap, k)
  }

  return result
}

func (h *httpRequestMessage) toHTTP() []byte {
  b := &bytes.Buffer{}

  // write request line
  b.WriteString(fmt.Sprintf("%s %s %s", h.requestLine.method, h.requestLine.url, h.requestLine.version))
  b.WriteString(LINE_BREAK)

  // write headers
  for k, v := range h.headers {
    b.WriteString(k + HEADER_SEP + v + LINE_BREAK)
  }

  b.WriteString(LINE_BREAK)
  b.Write(h.body)
  return b.Bytes()
}

type httpRequestMessage struct {
  requestLine
  headers map[string]string
  body []byte
  singleHopHeaders
}

type singleHopHeaders struct {
  connection string
}

type requestLine struct {
  method string
  url string
  version string
}
