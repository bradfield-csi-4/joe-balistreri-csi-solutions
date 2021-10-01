default rel

section .text
global volume
volume:
  mulss xmm0, xmm0
  mulss xmm0, xmm1
  mulss xmm0, [pi]
  divss xmm0, [three]
 	ret

section .data
pi: dd 3.141592
three: dd 3.0
