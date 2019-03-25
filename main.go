package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := NewServer()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		s.Shutdown()
	}()

	s.Run()
}