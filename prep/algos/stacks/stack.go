package main

import (
  "fmt"
  "bytes"
)

type Stack interface {
  Push(interface{})
  Pop() interface{}
  Peek() interface{}
  IsEmpty() bool
  Size() int
}


type StackImpl struct {
  items []interface{}
}

func NewStack() *StackImpl {
  return &StackImpl{}
}

func (s *StackImpl) Push(i interface{}) {
  s.items = append(s.items, i)
}

func (s *StackImpl) Pop() interface{} {
  return nil
}

func (s *StackImpl) Peek() interface{} {
  return nil
}

func (s *StackImpl) IsEmpty() bool {
  return false
}

func (s *StackImpl) Size() int {
  return len(s.items)
}

func (s *StackImpl) String() string {
  b := &bytes.Buffer{}
  b.WriteString("{ ")
  for _, v := range s.items {
    fmt.Fprintf(b, "%s, ", v)
  }
  b.WriteString(" }")
  return b.String()
}

type A struct {
  B interface{}
}

func (a *A) String() string {
  return fmt.Sprintf("%v", a.B)
}

func main() {
  s := NewStack()
  fmt.Println(s)
  s.Push("b")
  fmt.Println(s)
}
