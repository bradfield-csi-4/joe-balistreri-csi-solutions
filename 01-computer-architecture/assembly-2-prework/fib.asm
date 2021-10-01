section .text
global fib
fib:
; 	cmp rdi, 1				; compare the arg to 1
; 	jle .base_case  	; exit if <= 1
;
; 	; fib(n-1)
; 	dec rdi
; 	push rdi
;
; 	call fib
;
; 	pop rdi
; 	mov rbx, rax
;
; 	push rbx
;
; 	; fib(n-2)
; 	dec rdi
; 	call fib
; 	pop rbx
; 	add rax, rbx
; 	ret
;
; .base_case:
; 	mov rax, rdi ;
; 	ret

	xor rax, rax
	mov rbx, 1
.top:
	cmp rdi, 0
	je .exit
	; rax = curr, rbx = next, rcx = tmp
	mov rcx, rax
	mov rax, rbx
	add rbx, rcx
	dec rdi
	jmp .top

.exit:
	ret
