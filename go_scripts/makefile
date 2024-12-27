all: autogamma

autogamma:
	go build autogamma.go

clean:
	rm -f autogamma

install:
	mkdir -p ~/.local/bin
	mv autogamma ~/.local/bin/
