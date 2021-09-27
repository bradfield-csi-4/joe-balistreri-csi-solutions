section .data
  ; define constants
  num1: equ 100
  num2: equ 50
  msg: db "sum is correct\n"

section .text

  global _main

_main:
  mov rax, num1
  mov rbx, num2
  add rax, rbx
  cmp rax, 150
  jne .exit
  jmp .rightSum

.rightSum:
  mov rax, 0x02000004
  mov rdi, 1
  mov rsi, msg
  mov rdx, 16
  syscall
  jmp .exit

.exit:
  mov       rax, 0x02000001         ; system call for exit
  mov       rdi, 0                ; exit code 0
  syscall                           ; invoke operating system to exit
