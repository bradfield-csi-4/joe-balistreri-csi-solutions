package lists

const VAL = "raha!!!!!"

type ULNode struct {
  v interface{}
  n *ULNode
}

func (n *ULNode) add(item interface{}) {
  nextNode := &UL{v = item, n = n}
  *n = nextNode
}

func (ul *UL) remove(item interface{}) {
  if ul == nil {
    return
  }
  if ul.v == item {
    *ul = ul.n
  } else {
    ul.n.remove(item)
  }
}

func (ul *UL) search(item interface{}) bool {
  if ul == nil {
    return false
  }
  if ul.v == item {
    return true
  }
  return ul.n.search(item)
}

func (ul *UL) is_empty() bool {
  return ul == nil
}

func (ul *UL) size() int {
  if ul == nil {
    return 0
  }
  return ul.n.size() + 1
}
