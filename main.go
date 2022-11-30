package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mertdogan12/led-daemon/config"
	"github.com/mertdogan12/led-daemon/internal/uds"
)

func main() {
	c := &config.Config{}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChan:
			log.Printf("Got Signal to exit")
			os.Exit(1)
		}
	}()

	if err := run(c); err != nil {
		os.Exit(1)
	}
}

func run(c *config.Config) error {
	c.Init(os.Args)

	uds.Run()

	return nil
}
