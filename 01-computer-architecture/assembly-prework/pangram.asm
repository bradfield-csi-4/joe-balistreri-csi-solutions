section .text
global pangram
pangram:
	mov cl, [rdi] 			; load the next character of the string

	cmp cl, 0						; proceed to exit if string has terminated
	je .exit

	sub cl, 65

	mov rbx, 1					; move 1 into rbx
	sal rbx, cl					;
	or rax, rbx
	inc rdi
	jg pangram
.exit:
	mov rbx, rax
	shr rbx, 32 					; shift a-z into A-Z positions
	or rax, rbx
	and rax, 0b11111_11111_11111_11111_111111 ; clear the 38 high bits of rax
	sub rax, 0b11111_11111_11111_11111_111111
	js .condition_false
	mov rax, 1
	ret
.condition_false:
	mov rax, 0
	ret
