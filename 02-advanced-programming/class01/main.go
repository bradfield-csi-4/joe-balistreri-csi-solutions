package main

import (
  "fmt"
  "unsafe"
)

var WordSize uintptr = 0
var bucketsOffset uintptr = 0
var BOffset uintptr = 0

const bucketCnt = 8
const keySize = 8
const elemSize = 8
const topHashSize = bucketCnt * 1
const bmapSize = topHashSize + (keySize + elemSize) * bucketCnt + 8 //word size
const elemOffset = topHashSize + keySize * bucketCnt

func init() {
  var i int
  WordSize = unsafe.Sizeof(i)
  fmt.Println("wordSize: ", WordSize)
  BOffset = WordSize + 1
  bucketsOffset = WordSize + 8
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

  m := map[int]int {
    1: 10,
    2: 23124,
    3: 5432,
  }
  for i := 4; i < 10000; i++ {
    m[i] = i
  }
  fmt.Println(myMax(m))
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

type hmap struct {
  count int
  flags uint8
  B uint8
  noverflow uint16
  hash0 uint32
  buckets *intBucket
  oldbuckets *intBucket
  nevacuate uintptr
  extra *mapextra
}

type mapextra struct {
  overflow *[]*intBucket
  oldoverflow *[]*intBucket
  nextOverflow *intBucket
}

type intBucket struct {
  tophash [8]uint8
  keys [8]int
  elems [8]int
  overflow *intBucket
}

func myMax(m map[int]int) int {
  hmap := *(*hmap)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&m))))
  totalBuckets := 1 << hmap.B
  max := 0
  var exampleBucket intBucket

  // loop through each bucket
  for i := 0; i < totalBuckets; i++ {
    bucket := *(*intBucket)(unsafe.Pointer(uintptr(unsafe.Pointer(hmap.buckets)) + uintptr(i) * unsafe.Sizeof(exampleBucket)))

    // for each bucket, loop through the elements and then check the overflow bucket
  overflowloop:
    for {
      for j := 0; j < 8; j++ {
        nextVal := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&bucket.elems)) + uintptr(j) * elemSize))
        if nextVal > max {
          max = nextVal
        }
      }
      if bucket.overflow == nil {
        break overflowloop
      }
      bucket = *bucket.overflow
    }
  }

  return max
}
