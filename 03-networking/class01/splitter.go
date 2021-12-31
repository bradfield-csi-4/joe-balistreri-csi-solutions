package main

import (
  "os"
  "bytes"
  "strings"
  "strconv"
  "sort"
  "fmt"
  "encoding/binary"
  "encoding/hex"
)

const filename = "net.cap"
const capPacketHeaderLen = 16

type serverResp struct {
  sequenceNo int
  payload []byte
}

func main() {

  data, err := os.ReadFile(filename)
  if err != nil {
    panic(err)
  }

  // remove cap file header
  data = data[24:]

  var clientIp []byte = nil

  var responses []serverResp
  seen := map[int]bool{}

  stdoutDumper := hex.Dumper(os.Stdout)
	defer stdoutDumper.Close()
  for len(data) > 0 {
    fmt.Println("\n\n** new packet **")

    // parse cap file header
    packetLengthBytes := data[8:12]
    untruncatedPacketLengthBytes := data[12:16]
    pl := binary.LittleEndian.Uint32(packetLengthBytes)
    upl := binary.LittleEndian.Uint32(untruncatedPacketLengthBytes)
    if pl != upl {
      panic("help!")
    }


    // parse ethernet frame header
    ethernetFrame := data[capPacketHeaderLen : capPacketHeaderLen + pl]


    // parse ipv4 datagram header
    ipV4Datagram := ethernetFrame[14:]
    sourceIp := ipV4Datagram[12:16]
    if clientIp == nil {
      clientIp = sourceIp
    }
    fromServer := false
    if string(sourceIp) == string(clientIp) {
      fmt.Println("client request")
    } else {
      fromServer = true
      fmt.Println("server response")
    }
    ipDatagramLength := binary.BigEndian.Uint16(ipV4Datagram[2:4])
    ipV4Datagram = ipV4Datagram[:ipDatagramLength]

    ipV4HeaderLengthBytes := int(ipV4Datagram[0] << 2 >> 2) * 4 // header length is in 32-bit words


    // parse tcp segment header and collect server responses
    tcpSegment := ipV4Datagram[ipV4HeaderLengthBytes:]
    seqNo := int(binary.BigEndian.Uint32(tcpSegment[4:8]))
    ackNo := int(binary.BigEndian.Uint32(tcpSegment[8:12]))
    fmt.Printf("seqNo: %d\nackNo: %d\n", seqNo, ackNo)

    tcpHeaderLengthBytes := (tcpSegment[12] >> 4) * 4 // header length in 32 bit words
    tcpPayload := tcpSegment[tcpHeaderLengthBytes:]
    stdoutDumper.Write(tcpPayload)

    if fromServer {
      if !seen[seqNo] {
        seen[seqNo] = true
        responses = append(responses, serverResp{sequenceNo: seqNo, payload: tcpPayload})
      }
    }

    data = data[capPacketHeaderLen + pl:]
  }

  // sort and combine server responses into single payload
  sort.SliceStable(responses, func(i, j int) bool {
      return responses[i].sequenceNo < responses[j].sequenceNo
  })
  result := &bytes.Buffer{}
  for _, resp := range responses {
    result.Write(resp.payload)
  }

  // remove HTTP headers and write payload body to output file
  r := result.String()
  var contentLength int
  for _, line := range strings.Split(r, "\n") {
    if strings.HasPrefix(line, "Content-Length:") {
      v := strings.Split(line, ": ")[1]
      v = v[:len(v)-1]
      contentLength, err = strconv.Atoi(v)
      break
    }
  }
  payload := r[len(r)-contentLength:len(r)]
  os.WriteFile("output", []byte(payload), 0644)
}
