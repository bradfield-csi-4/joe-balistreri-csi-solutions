package compute

import (
  "fmt"
)

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
