package main

import (
  "fmt"
  "syscall"
  "strings"
  "encoding/binary"
  // "net"
)

const requestId = 14

func main() {
  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
  defer syscall.Close(fd)

  if err != nil {
    panic(err)
  }

  if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
    panic(err)
  }

  err = syscall.Bind(fd, &syscall.SockaddrInet4{
    Port: 7890,
  })
  if err != nil {
    panic(err)
  }

  payload := createQueryPayload()

  err = syscall.Sendto(fd, payload, 0, &syscall.SockaddrInet4{
    Port: 53,
    Addr: [4]byte{8,8,8,8},
  })
  if err != nil {
    panic(err)
  }

  response := make([]byte, 65536)
  for {
    n, _, err := syscall.Recvfrom(fd, response, 0)
    if err != nil {
      panic(err)
    }
    if n <= 0 {
      continue
    }
    response = response[:n]
    break
  }

  printResponse(response)
}

func printResponse(payload []byte) {
  fmt.Println(payload)

  // handle header values
  header := payload[:12]
  id := binary.BigEndian.Uint16(header[0:2])
  fmt.Printf("id: %d\n", id)
  if (header[3] >> 7) == 1 {
    fmt.Println("it's a response!")
  }
  opCode := (header[3] << 1) >> 4
  if opCode == 0 {
    fmt.Println("standard query!")
  } else {
    fmt.Println("other query type!")
  }

  questions := binary.BigEndian.Uint16(header[4:6])
  answerRRs := binary.BigEndian.Uint16(header[6:8])
  authorityRRs := binary.BigEndian.Uint16(header[8:10])
  additionalRRs := binary.BigEndian.Uint16(header[10:12])
  fmt.Printf("questions: %d\nanswerRRs: %d\nauthorityRRs: %d\nadditionalRRs: %d\n", questions, answerRRs, authorityRRs, additionalRRs)

  payload = payload[12:]

  if questions > 0 {
    fmt.Println(";; QUESTION SECTION:")
    for ; questions > 0; questions-- {
      payload = printQuestion(payload)
    }
    fmt.Println()
  }

  if answerRRs > 0 {
    fmt.Println(";; ANSWER SECTION:")
    for ; answerRRs > 0; answerRRs-- {
      payload = printAnswer(payload)
    }
    fmt.Println()
  }

  if authorityRRs > 0 {
    fmt.Println(";; AUTHORITY SECTION:")
    for ; authorityRRs > 0; authorityRRs-- {
      payload = printAnswer(payload)
    }
    fmt.Println()
  }

  // note: skipping additional section
}

var qtypes = map[int]string {
  1: "A",
  2: "NS",
  3: "MD",
  4: "MF",
  5: "CNAME",
  15: "MX",
  16: "TXT",
}

var qclasses = map[int]string {
  1: "IN",
  2: "CS",
  3: "CH",
  4: "HS",
}

func printQuestion(payload []byte) []byte {
  payload, name := readName(payload)
  fmt.Printf("%s", name)
  qtype := int(binary.BigEndian.Uint16(payload[0:2]))
  qclass := int(binary.BigEndian.Uint16(payload[2:4]))
  fmt.Printf("\t\t\t\t%s\t%s\n", qclasses[qclass], qtypes[qtype])
  return payload[4:]
}

func printAnswer(payload []byte) []byte {
  payload, name := readName(payload)
  fmt.Printf("%s", name)
  _type := payload[0:2]
  class := payload[2:4]
  ttl := binary.BigEndian.Uint32(payload[4:8]) // seconds
  rdLength := binary.BigEndian.Uint16(payload[8:10])
  rdData := payload[10:10+rdLength]
  fmt.Println(_type, class, ttl, rdLength, rdData)

  return payload[10+rdLength:]
}

func printAuthority(payload []byte) []byte {
  if payload[0] >> 6 == 3 {
    panic("message compression detected")
  }

  fmt.Println("authority")
  return payload
}

func readName(payload []byte) ([]byte, string) {
  if (payload[0] >> 6) == 3 {
    panic("message compression detected")
  }
  result := "" // TODO: use string buffer
  for l := payload[0]; l != 0; {
    label := payload[1:l+1]
    result += fmt.Sprintf("%s.", label)
    payload = payload[l+1:]
    l = payload[0]
  }
  payload = payload[1:] // remove trailing 0 byte from QNAME section
  return payload, result
}

func createQueryPayload() []byte {
  header := []byte{
    0, requestId, // ID
    1, // recursion desired
    0, // ra, z, rcode
    0, 1, // qd count (question)
    0, 0, // an count (answer rrs)
    0, 0, // ns count (name server rrs)
    0, 0, // ar count (additional rrs)
  }

  url := "www.google.com"

  encodedUrl := []byte{}
  for _, piece := range strings.Split(url, ".") {
    encodedUrl = append(encodedUrl, uint8(len(piece)))
    encodedUrl = append(encodedUrl, []byte(piece)...)
  }
  encodedUrl = append(encodedUrl, uint8(0))


  questionSection := append(encodedUrl,
    []byte{
      0, 1, // A record
      0, 1, // IN QClass
    }...,
  )


  payload := append(header, questionSection...)
  return payload
}
