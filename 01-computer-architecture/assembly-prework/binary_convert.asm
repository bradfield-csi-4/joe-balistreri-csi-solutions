section .text
global binary_convert
binary_convert:
.start:
.top:
	mov bl, [rdi]
	cmp bl, byte 0 				; compare to null terminator
	je .exit								; if <, it's the null terminator

	inc rdi 								; move to the next position
	sal rax, 1							; multiply rax by 2
	cmp bl, byte 48  			; compare to zero ascii character
	je .top

	inc rax
	jmp .top
.exit:
	ret
