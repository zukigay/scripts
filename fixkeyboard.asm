section .data
    ;; file name is null terminated
    open_error_s db 'error opening file',10 ; 19
    write_error_s db 'error writing to file',10 ; 22
    loop_s db 'looping',10 ; 8
    filename db '/sys/module/hid_apple/parameters/fnmode', 0
    dataToWrite db '1'
secs dq 5,0
section .bss
    fd resb 8
section .text
    global _start

open_error:
    mov rax, 1
    mov rdi, 2
    mov rsi, open_error_s
    mov rdx, 19
    syscall
    mov rax, 60
    mov rdi, 1
    syscall
write_error:
    mov rax, 1
    mov rdi, 2
    mov rsi, write_error_s
    mov rdx, 22
    syscall
    mov rax, 60
    mov rdi, 2
    syscall
_start:
    mov rax, 2
    mov rdi, filename
    mov rsi, 1 ; for O_WRONLY
    mov rdx, 0o664 ; the mode
    syscall
    ; check for error
    cmp rax,0
    jle open_error
    mov [fd], rax ; put return value into fd

    ; write to file
    mov rax, 1
    mov rdi, [fd]
    mov rsi, dataToWrite
    mov rdx, 1 ; writing just one char
    syscall
    ; check for error
    cmp rax,0
    jle write_error

    ; close file
    mov rax, 3
    mov rdi, [fd]
    syscall

    mov rax, 35
    lea rdi, secs
    xor rsi, rsi
    syscall
    mov rax, 1
    mov rdi, 2
    mov rsi, loop_s 
    mov rdx, 8
    syscall
    jmp _start

    ; exit with 0 error code
    mov rax, 60
    mov rdi, 0
    syscall
