package main

import "fmt"

const CHARS = "0123456789abcdef"


func convert(n, base int) string {
  if n < base {
    return string(CHARS[n])
  }
  return convert(n / base, base) + string(CHARS[n % base])
}


func main() {
  fmt.Println(convert(1453, 16))
}
