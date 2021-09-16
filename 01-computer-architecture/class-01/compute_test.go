package compute

import "testing"


func TestComputeClassExample(t *testing.T) {
  example_memory := [256]byte{
    0x00, // output data
    0x03, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, // input data
    0x01, 0x01, 0x01, // load r1 0x01
    0x01, 0x02, 0x02, // load r2 0x02
    0x03, 0x01, 0x02, // add  r1 r2
    0x02, 0x01, 0x00, // store r1 0x00
    0xff, // halt
  }

  compute(&example_memory)

  if example_memory[0] != 8 {
    t.Errorf("Computation was incorrect. Got %x, want: %x", example_memory[0], 8)
  }
}
