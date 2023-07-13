package main

import (
	"log"

	"github.com/reiver/go-telnet"
)

func main() {
	var handler telnet.Handler = telnet.EchoHandler
	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		log.Fatal(err)
	}
}
