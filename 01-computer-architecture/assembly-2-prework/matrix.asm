section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex
	lea rbx, [rdx * 4]     ; compute row-width
	imul rbx, rcx					 ; multiply row-width by rindex
	lea rbx, [rbx + r8*4]  ; add col-index * elem size (4)
	mov rax, [rdi + rbx]   ; read from array start + offset
	ret
