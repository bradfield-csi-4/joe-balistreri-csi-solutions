section .text
global sum_to_n
sum_to_n:
; O(n) solution
; .start:
; 					add       rax, rdi    ; add the input to the total stored in rax
; 					dec 			rdi					; decrement the input value
; 					cmp       rdi, 0	    ; compare the input value to 0
; 					jg				.start      ; if >0, loop again
; 					ret 									; return


; O(1) solution
					lea       rax, [rdi + 1]
					imul			rdi
					sar  			rax, 1
					ret
