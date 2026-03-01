all: spirefix crtfree
spirefix: spirefix.c
	cc spirefix.c -o spirefix -lX11
crtfree: crtfree.c
	cc crtfree.c -nostdlib -o crtfree -Wall -Wextra
clean:
	rm spirefix
