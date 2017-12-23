package main

import (
	"os"
	"os/signal"
	"syscall"
)

func createSigCh() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	return ch
}
