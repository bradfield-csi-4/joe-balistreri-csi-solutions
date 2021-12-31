package main

/*
#include <stdlib.h>
#include <stdio.h>
#include "square.c"
*/
import "C"

import (
  "fmt"
  "unsafe"
)

func Random() int {
  return int(C.random())
}

func Seed(i int) {
  C.srandom(C.uint(i))
}

func main() {
  gs := fmt.Sprintf("%d", C.square(13))
  cs := C.CString(gs)
  defer 8C.free(unsafe.Pointer(cs))
  C.fputs(cs, (*C.FILE)(C.stdout))
}
