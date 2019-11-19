package main

import (
	"fmt"
	"os"
	"os/signal"
)

func waitForSignal() {
	for {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		// Block until a signal is received.
		s := <-c
		fmt.Println("Got signal:", s)

		if s == os.Kill {
			// how can one send a 'kill' signal to process on windows?
			dumpGoRoutinesInfo()
		}

		break
	}
}
