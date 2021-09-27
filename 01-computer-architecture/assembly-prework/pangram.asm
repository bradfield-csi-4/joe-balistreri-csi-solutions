section .text
global pangram
pangram:
	dec rdi
.top:
	inc rdi
	mov cl, [rdi]

	; exit if we've hit the end of the string
	cmp cl, 0
	je .exit

	; go back to the top if the value is <65
	sub cl, 65
	js .top

.store:
	mov rbx, 1
	sal rbx, cl
	or rax, rbx
	jg .top
.exit:
	mov rbx, rax
	shr rbx, 32 					; shift a-z into A-Z positions
	or rax, rbx
	sal rax, 38						; clear the 38 high bits of rax
	shr rax, 38
	sub rax, 0b11111_11111_11111_11111_111111
	js .real_exit
	mov rax, 1
	ret
.real_exit:
	mov rax, 0
	ret
