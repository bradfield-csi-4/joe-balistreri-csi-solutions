section .text
global pangram
pangram:
	mov cl, [rdi]
	inc rdi

	; exit if we've hit the end of the string
	cmp cl, 0
	je .exit

	; go back to the top if the value is <65
	sub cl, 65
	js pangram

.store:
	mov rbx, 1
	sal rbx, cl
	or rax, rbx
	jg pangram
.exit:
	mov rbx, rax
	shr rbx, 32 					; shift a-z into A-Z positions
	or rax, rbx
	and rax, 0b11111_11111_11111_11111_111111 ; clear the 38 high bits of rax
	sub rax, 0b11111_11111_11111_11111_111111
	js .real_exit
	mov rax, 1
	ret
.real_exit:
	mov rax, 0
	ret
