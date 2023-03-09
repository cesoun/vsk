package main

import (
	"github.com/cesoun/vsk/pkg/riot"
	"os"
	"os/signal"
)

func main() {
	rc := riot.NewClient()
	rc.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	rc.Stop()
}
