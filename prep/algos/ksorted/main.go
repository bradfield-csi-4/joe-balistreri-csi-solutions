package main

import (
  "bytes"
  "fmt"
)

type ListNode struct {
    Val int
    Next *ListNode
}

func (ln *ListNode) String() string {
  if ln == nil {
    return "nil"
  }
  return fmt.Sprintf("%d ->", ln.Val)
}

func mergeKLists(lists []*ListNode) *ListNode {
  return nil
}


type minHeap struct {
    nodes []*ListNode
    count int
}

func (mh *minHeap) pop() *ListNode {
  if mh.count == 0 {
    return nil
  }
  returnValue := mh.nodes[1]

  mh.nodes[1] = mh.nodes[mh.count]
  mh.nodes = mh.nodes[:mh.count]
  mh.count--

  mh.propagate_down(1)

  return returnValue
}

func (mh *minHeap) propagate_down(i int) {
  child_a, child_b := i * 2, i * 2 + 1
  if child_a > mh.count {
    return
  }
  if child_b > mh.count {
    if mh.nodes[child_a].Val < mh.nodes[i].Val {
      mh.nodes[i], mh.nodes[child_a] = mh.nodes[child_a], mh.nodes[i]
    }
    return
  }
  if mh.nodes[child_a].Val <= mh.nodes[child_b].Val && mh.nodes[child_a].Val < mh.nodes[i].Val {
    mh.nodes[i], mh.nodes[child_a] = mh.nodes[child_a], mh.nodes[i]
    mh.propagate_down(child_a)
  } else if mh.nodes[child_b].Val <= mh.nodes[child_a].Val && mh.nodes[child_b].Val < mh.nodes[i].Val {
    mh.nodes[i], mh.nodes[child_b] = mh.nodes[child_b], mh.nodes[i]
    mh.propagate_down(child_b)
  }
}

func (mh *minHeap) add(ln *ListNode) {
  mh.count += 1
  mh.nodes = append(mh.nodes, ln)
  mh.propagate_up(mh.count)
}

func (mh *minHeap) propagate_up(i int) {
  if i == 1 {
    return
  }
  parent := i / 2
  if mh.nodes[i].Val >= mh.nodes[parent].Val {
    return
  }
  mh.nodes[i], mh.nodes[parent] = mh.nodes[parent], mh.nodes[i]
  mh.propagate_up(parent)
}

func newMinHeap() *minHeap {
    return &minHeap{
        count: 0,
        nodes: []*ListNode{nil},
    }
}

func (mh *minHeap) String() string {
  var b bytes.Buffer
  b.WriteString("[")
  for _, ln := range mh.nodes {
    b.WriteString(ln.String())
    b.WriteString(", ")
  }
  b.WriteString("]")
  return b.String()
}


func main() {
  fmt.Println("hello!")
  mh := newMinHeap()
  mh.add(&ListNode{Val: 10})
  mh.add(&ListNode{Val: 4})
  mh.add(&ListNode{Val: 4})
  mh.add(&ListNode{Val: 2})
  mh.add(&ListNode{Val: 27})
  mh.add(&ListNode{Val: 217})
  mh.add(&ListNode{Val: -1})
  mh.add(&ListNode{Val: 0})
  mh.add(&ListNode{Val: 5})
  fmt.Println(mh)

  for len(mh.nodes) > 1 {
    fmt.Println(mh.pop())
    fmt.Println(mh)
    fmt.Println("\n\n")
  }
}
