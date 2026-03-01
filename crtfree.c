// static void exit(void);
// {
//     __asm volatile (
//         "mov  %%rdi, %%rsi\n"     // arg2 = stack
//         "mov  $0x50f00, %%edi\n"  // arg1 = clone flags
//         "mov  $56, %%eax\n"       // SYS_clone
//         "syscall\n"
//         "mov  %%rsp, %%rdi\n"     // entry point argument
//         "ret\n"
//         : : : "rax", "rcx", "rsi", "rdi", "r11", "memory"
//     );
// }
void crt_exit(void) {
    __asm__ (
        "movl $0, %ebx\n\t"
        "movl $1, %eax\n\t"
        "int $0x80\n\t"
    );
}
void _start(void)
{
    crt_exit();
}
