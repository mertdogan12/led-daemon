package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mazznoer/colorgrad"
	"github.com/mertdogan12/led-daemon/config"
	"github.com/mertdogan12/pulse-simple"
	truecolor "github.com/wayneashleyberry/truecolor/pkg/color"
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

const BUFSIZE = 2004

func run(c *config.Config) error {
	c.Init(os.Args)

	// go uds.Run()
	// led.Run()

	ss := pulse.SampleSpec{
		Format:   pulse.SAMPLE_FLOAT32LE,
		Rate:     44100,
		Channels: 1,
	}

	stream, err := pulse.Capture("led", "music", "music", &ss)
	if err != nil {
		log.Fatal(err)
	}

	defer stream.Free()
	defer stream.Drain()

	fadeGrad := colorgrad.Rainbow()

	for {
		buf := make([]byte, BUFSIZE)

		_, err = stream.Read(buf)
		if err != nil {
			log.Fatal("Error while reading: ", err)
		}

		pcm := convToFloat(buf)

		for _, p := range pcm {
			min := float32(-0.5)
			max := float32(0.5)

			if p < min || p > max {
				fmt.Println(p)
				continue
			}

			difference := p - min
			value := difference
			out := make([]byte, int(value*100))

			for i := range out {
				out[i] = ' '
			}

			rgb := fadeGrad.At(float64(value))

			// fmt.Println("\033[2J")

			latency, err := stream.Latency()
			if err != nil {
				return err
			}
			fmt.Print(latency)

			truecolor.
				White().
				Background(uint8(rgb.R*100), uint8(rgb.G*100), uint8(rgb.B*100)).
				Print(fmt.Sprintf("%s%f\n", string(out), p))

			time.Sleep(time.Millisecond)
		}
	}
}

func convToFloat(inp []byte) []float32 {
	out := make([]float32, 0)
	tmp := make([]byte, 4)

	var i = 0
	for _, b := range inp {
		if i > 3 {
			i = 0

			fBits := binary.LittleEndian.Uint32(tmp)
			f := math.Float32frombits(fBits)
			out = append(out, f)
		}

		tmp[i] = b
		i++
	}

	return out
}
