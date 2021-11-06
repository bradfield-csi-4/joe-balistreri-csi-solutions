package main

import (
	"fmt"
	"unsafe"
)


// Return string length without use "len"
func strLenSafe(s string) int  {
	count := 0
	for range s {
		count++
	}
	return count
}

func strLenUnsafe(s string) int  {
	return *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + uintptr(8)))
}

// Return y coordinate of a point without using "p.y"
type Point struct { x , y int }

func getY(p Point) int {
	return *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + uintptr(8)))
}


// Return sum of []int without using range or []
func sumIntArray(a []int) int {
	sum := 0
	intSize := unsafe.Sizeof(int(0))
	arrLen := len(a)
	ap := (uintptr((*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)))))))
	for i := 0; i < arrLen; i++ {
		sum += *(*int)(unsafe.Pointer(ap + uintptr(i)*intSize))
	}
	return sum
}



// Given a map[int]int return the max value without using range or []

type hmap struct {
	count     	int
	flags     	uint8
	B         	uint8
	noverflow 	uint16
	hash0     	uint32
	buckets    	unsafe.Pointer
	oldbuckets 	unsafe.Pointer
	nevacuate  	uintptr
}

func getMapMaxVal(m map[int]int) int {

	max := 0

	hmap := *(**hmap)(unsafe.Pointer(uintptr(unsafe.Pointer(&m))))

	fmt.Println("count: ", hmap.count)
	fmt.Println("buckets: ", hmap.buckets)
	fmt.Println("hash0: ", hmap.hash0)
	fmt.Println("B: ", hmap.B)

	return max
}

func main() {

	// Try out string length
	s := "hello, world"
	l := strLenSafe(s)
	fmt.Println(l)
	l = strLenUnsafe(s)
	fmt.Println(l)

	// Try out getting y coordinate
	p := Point{5, 3}
	y := getY(p)
	fmt.Println(y)

	// Try out array sum
	a := []int{2, 3, 5, 7, 11}
	sum := sumIntArray(a)
	fmt.Println(sum)

  fmt.Println("joe")
	// Try out map max value
	m := map[int]int{
		1:	435,
		2:	824,
		3:	234,
		4: 	123,
	}
	max := getMapMaxVal(m)
	fmt.Println(max)

}
