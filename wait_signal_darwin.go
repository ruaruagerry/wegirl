package main

import (
	"fmt"
	"os"
	"os/signal"
	"wegirl/servercfg"
	"syscall"
)

func waitForSignal() {
	for {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)

		// Block until a signal is received.
		s := <-c
		fmt.Println("Got signal:", s)

		if s == syscall.SIGUSR1 {
			dumpGoRoutinesInfo()
			continue
		}

		if s == syscall.SIGUSR2 {
			servercfg.ReLoadConfigFile()
			continue
		}

		break
	}
}
