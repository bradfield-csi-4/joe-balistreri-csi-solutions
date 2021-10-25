package main

import (
  "fmt"
  "unsafe"
)

var WordSize uintptr = 0

func init() {
  var i int
  WordSize = unsafe.Sizeof(i)
}

type Point struct {
  x int
  y int
}

func main() {
  s := "hello bob"
  fmt.Println(myLen(s))
  p := Point{10, 123}
  fmt.Println(yValue(p))
  i := []int{1,2,3,4,5,6}
  fmt.Println(mySum(i))
}

func myLen(s string) int {
  return *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + WordSize))
}

func yValue(p Point) int {
  return *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + WordSize))
}

func mySum(a []int) int {
  var total int
  l := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + WordSize))
  for i := 0; i < l; i++ {
    total += *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&a)) + WordSize * uintptr(i)))))
  }
  return total
}
