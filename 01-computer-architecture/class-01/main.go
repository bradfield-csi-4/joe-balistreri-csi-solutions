package main

import (
  "fmt"
)

func main() {
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
  fmt.Println(example_memory[0])
}

func compute(memory *[256]byte) {
  // start program counter at the beginning of the instructions
  registers := [3]byte{0x08}

  for {
    pc := registers[0]
    switch i := memory[pc]; i {
    case 0x01: // Load
      reg, addr := memory[pc+1],  memory[pc+2]
      registers[reg] = memory[addr]
    case 0x02: // Store
      reg, addr := memory[pc+1], memory[pc+2]
      memory[addr] = registers[reg]
    case 0x03: // Add
      reg1, reg2 := memory[pc+1], memory[pc+2]
      registers[1] = registers[reg1] + registers[reg2]
    case 0x04: // Sub
      reg1, reg2 := memory[pc+1], memory[pc+2]
      registers[1] = registers[reg1] - registers[reg2]
    case 0xff: // Halt
      return
    default:
      panic(fmt.Sprintf("unknown instruction %x", i))
    }

    registers[0] += 3
  }
}
