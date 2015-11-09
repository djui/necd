package main

import (
	"os"
	"os/signal"

	"github.com/ian-kent/go-log/log"
)

func KeepAlive(f func()) {
	go f()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func AssertNoErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func Assert(cond bool, msg string) {
	if !cond {
		log.Fatal(msg)
	}
}
