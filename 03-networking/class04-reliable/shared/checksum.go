package shared

import (
    "log"
)

const zero uint16 = 0
const ValidSum uint16 = ^zero

func SumBytes(header Header) (sum uint16) {
  b := HeaderToBytes(header)
  if len(b) % 2 != 0 {
    b = append(b, 0)
  }
  for i := 0; i < len(b); i += 2 {
    sum = sumWithCarry(bytesToUint16(b[i:i+2]), sum)
  }
  return sum
}

func sumWithCarry(a, b uint16) uint16 {
  sum := a + b
  if sum < a || sum < b {
    sum += 1
  }
  return sum
}

func bytesToUint16(bytes []byte) (result uint16) {
	b1, b2 := bytes[0], bytes[1]
	return result + uint16(b1)<<8 + uint16(b2)
}

func uInt16ToBytes(u uint16) (result []byte) {
	result = append(result, byte(u>>8))
	result = append(result, byte(u))
	return result
}

type Header struct {
  Length uint16 // in bytes
  Checksum uint16 // 1s complement of sum of all 16 bit words; padding at end of data if needed
  Data []byte
}


func NewHeader(data []byte) Header {
  result := Header{
    Length: uint16(2 + len(data)),
    Data: data,
  }
  result.Checksum = ^SumBytes(result)

  return result
}

func HeaderFromBytes(encoded []byte) Header {
  if len(encoded) < 4 {
    log.Fatalf("invalid encoded header %v", encoded)
  }
  return Header{
    Length: bytesToUint16(encoded[0:2]),
    Checksum: bytesToUint16(encoded[2:4]),
    Data: encoded[4:],
  }
}

func HeaderToBytes(header Header) (result []byte) {
  result = append(result, uInt16ToBytes(header.Length)...)
  result = append(result, uInt16ToBytes(header.Checksum)...)
  result = append(result, header.Data...)
  return result
}
