section .text
global binary_convert
binary_convert:
	mov bl, [rdi]
	sub bl, 48			 				; subtract the value of the 0 ascii character
	js .exit								; if negative, bl was the null terminator

	inc rdi 								; move to the next position
	lea rax, [rbx, rax*2]   ; double rax and add the value in rbx
	jmp binary_convert
.exit:
	ret
