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

	; store the value if it's in the right range (A-Z)
	cmp cl, 26
	jl .store

	; subtract 32 to account for a-z
	sub cl, 32
	; if it goes negative, skip because we're in 91-96
	js .top
	; if >= 26, we're above a-z
	cmp cl, 26
	jge .top

.store:
	mov rbx, 1
	sal rbx, cl
	or rax, rbx
	jg .top
.exit:
	sub rax, 0b11111_11111_11111_11111_111111
	js .real_exit
	mov rax, 1
	ret
.real_exit:
	mov rax, 0
	ret
