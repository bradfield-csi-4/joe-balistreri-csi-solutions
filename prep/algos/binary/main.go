package main

import (
  "bytes"
  "fmt"
)

const DIGITS = "0123456789abcdef"

func convert_to_base(decimal_num int, base int) string {
  var stack []byte
  for decimal_num != 0 {
    stack = append(stack, DIGITS[decimal_num % base])
    decimal_num /= base
  }

  b := &bytes.Buffer{}
  for i := len(stack) -1; i >= 0; i-- {
    b.WriteByte(stack[i])
  }
  return b.String()
}

func main() {
  fmt.Println(convert_to_base(42, 2))
  fmt.Println(convert_to_base(233, 2))
  fmt.Println(convert_to_base(25, 2))
  fmt.Println(convert_to_base(26, 16))
}
