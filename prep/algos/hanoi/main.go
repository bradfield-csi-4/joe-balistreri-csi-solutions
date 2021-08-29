package main

import (
  "bytes"
  "fmt"
)

type Hanoi struct {
  towers [][]int
  n int
}

func NewHanoi(n int) *Hanoi {
  var slot1 []int
  for i := n; i > 0; i-- {
    slot1 = append(slot1, i)
  }
  return &Hanoi{
    n: n,
    towers: [][]int{slot1, []int{}, []int{}},
  }
}

func (h *Hanoi) can_move(from_i, to_i int) bool {
  from, to := h.towers[from_i], h.towers[to_i]
  if len(from) < 1 {
    panic(fmt.Sprintf("illegal from is empty %d", from_i))
  }
  from_value := from[len(from) - 1]
  return len(to) == 0 || to[len(to)-1] > from_value
}

func (h *Hanoi) move_disc(from_i, to_i int) {
  if !h.can_move(from_i, to_i) {
    panic("cannot move onto a smaller disk")
  }
  from := h.towers[from_i]
  h.towers[to_i] = append(h.towers[to_i], from[len(from)-1])
  h.towers[from_i] = h.towers[from_i][:len(from)-1]
  fmt.Println(h)
}

func (h *Hanoi) other_i(from_i, to_i int) int {
  switch from_i + to_i {
  case 1: return 2
  case 2: return 1
  case 3: return 0
  }
  return -1
}

func (h *Hanoi) move(from_i, to_i, from_depth int) {
  if from_depth == 1 && h.can_move(from_i, to_i) {
    h.move_disc(from_i, to_i)
    return
  }
  other_i := h.other_i(from_i, to_i)
  h.move(from_i, other_i, from_depth - 1)
  h.move_disc(from_i, to_i)
  h.move(other_i, to_i, from_depth - 1)
}

func (h *Hanoi) String() string {
  var b bytes.Buffer
  fmt.Fprintf(&b, "%v\n", h.towers[0])
  fmt.Fprintf(&b, "%v\n", h.towers[1])
  fmt.Fprintf(&b, "%v\n", h.towers[2])
  return b.String()
}

func main() {
  h := NewHanoi(5)
  fmt.Println(h)
  h.move(0, 2, 5)
}
