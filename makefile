all: spirefix crtfree fixkeyboard
fixkeyboard: fixkeyboard.asm
	nasm -f elf64 fixkeyboard.asm -o fixkeyboard.o
	ld fixkeyboard.o -o fixkeyboard
spirefix: spirefix.c
	cc spirefix.c -o spirefix -lX11
crtfree: crtfree.c
	cc crtfree.c -nostdlib -o crtfree -Wall -Wextra
clean:
	rm spirefix
