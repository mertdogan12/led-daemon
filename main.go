package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mertdogan12/led-daemon/config"
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

	log.Printf("Hallo world")

	time.Sleep(8 * time.Second)
}

func run(c *config.Config) error {
	c.Init(os.Args)

	return nil
}
