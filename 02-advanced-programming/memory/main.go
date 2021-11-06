package main

import (
  "fmt"
  "sync"
  "sync/atomic"
)


type AtomicMutex struct {
  v atomic.Value
}

func NewAtomicMutex() *AtomicMutex {
  m := &AtomicMutex{}
  m.v.Store(false)
  return m
}

func (m *AtomicMutex) Lock() {
  for !m.v.CompareAndSwap(false, true) {
  }
}

func (m *AtomicMutex) Unlock() {
  m.v.Store(false)
}


type ChannelMutex struct {
  ch chan int
}

func NewChannelMutex() *ChannelMutex {
  m := &ChannelMutex{ch: make(chan int, 1)}
  m.ch <- 0
  return m
}

func (m *ChannelMutex) Lock() {
  <-m.ch
}
func (m *ChannelMutex) Unlock() {
  m.ch<-0
}


func main() {
  // m := &sync.Mutex{}
  // m := NewChannelMutex()
  m := NewAtomicMutex()


  var a string
  wg := &sync.WaitGroup{}
  wg.Add(2)
  go func() {
    m.Lock()
    a = "a"
    m.Unlock()
    wg.Done()
  }()
  go func() {
    m.Lock()
    a = "b"
    m.Unlock()
    wg.Done()
  }()
  wg.Wait()
  m.Lock()
  fmt.Println(a)
  m.Unlock()
}
