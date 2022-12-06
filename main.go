package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mertdogan12/led-daemon/config"
	"github.com/mesilliac/pulse-simple"
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
		Format:   pulse.SAMPLE_S16NE,
		Rate:     96000,
		Channels: 2,
	}
	stream, err := pulse.Capture("led", "music", &ss)

	if err != nil {
		log.Fatal(err)
	}

	defer stream.Free()
	defer stream.Drain()

	out := make([]byte, 2)
	for {
		_, err = stream.Read(out)
		if err != nil {
			log.Fatal("Error while reading: ", err)
		}

		data := int16(binary.LittleEndian.Uint16(out))
		fmt.Println(data, out)

		time.Sleep(100 * time.Millisecond)
	}
}
