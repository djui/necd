all: brightness necd

necd: *.go
	go build

brightness:
	gcc -std=c99 -o brightness c/brightness.c -framework IOKit -framework ApplicationServices
