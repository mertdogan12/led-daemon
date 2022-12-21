package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mertdogan12/led-daemon/config"
	"github.com/mertdogan12/pulse-simple"
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

	// go uds.Run()
	// led.Run()

	ss := pulse.SampleSpec{
		Format:   pulse.SAMPLE_FLOAT32LE,
		Rate:     48000,
		Channels: 2,
	}
	stream, err := pulse.Capture("led", "music", "bluez_output.E0_9D_FA_E3_56_CA.1", &ss)

	if err != nil {
		log.Fatal(err)
	}

	defer stream.Free()
	defer stream.Drain()

	fmt.Println("left,right")

	out := make([]byte, 41943040/4)
	_, err = stream.Read(out)
	if err != nil {
		log.Fatal("Error while reading: ", err)
	}

	os.WriteFile("out.raw", out, 0644)

	return nil
}
