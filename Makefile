all: brightness.c
	go build

brightness.c:
	gcc -std=c99 -o brightness c/brightness.c -framework IOKit -framework ApplicationServices
