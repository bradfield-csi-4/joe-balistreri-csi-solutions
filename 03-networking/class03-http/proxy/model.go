package proxy

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
