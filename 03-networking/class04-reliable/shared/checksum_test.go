package shared

import (
  "testing"
)


func TestChecksum(t *testing.T) {
  h := Header{}

  sum := SumBytes(h)
  if sum != 0 {
    t.Errorf("SumBytes empty header - got %d, want 0", sum)
  }

  h.Length = 2
  h.Data = []byte{0, 2}

  sum = SumBytes(h)
  if sum != 4 {
    t.Errorf("SumBytes small header - got %d, want 4", sum)
  }

  h2 := NewHeader([]byte{1,2,3,4,5})
  if h2.Length != 7 {
    t.Errorf("NewHeader length - got %d, want 7", h2.Length)
  }

  sum = SumBytes(h2)
  if sum != ValidSum {
    t.Errorf("Sumbytes, new header - got %d, want %d", sum, ValidSum)
  }
}
