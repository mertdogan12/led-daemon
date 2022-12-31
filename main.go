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

func run(c *config.Config) error {
	c.Init(os.Args)

	// go uds.Run()
	// led.Run()

	ss := pulse.SampleSpec{
		Format:   pulse.SAMPLE_FLOAT32LE,
		Rate:     48000,
		Channels: 1,
	}
	stream, err := pulse.Capture("led", "music", "", &ss)

	if err != nil {
		log.Fatal(err)
	}

	defer stream.Free()
	defer stream.Drain()

	tmp := make([]byte, 4)
	// pcmData := make([]byte, 1024)
	// _, err = stream.Read(pcmData)
	// if err != nil {
	// 	log.Fatal("Error while reading: ", err)
	// }

	// pcm := convToFloat(pcmData)

	fadeGrad := colorgrad.Rainbow()
	outFile, err := os.Create("out.raw")
	if err != nil {
		panic(err)
	}

	for {
		_, err = stream.Read(tmp)
		if err != nil {
			log.Fatal("Error while reading: ", err)
		}

		tmpBits := binary.LittleEndian.Uint32(tmp)
		tmpFloat := math.Float32frombits(tmpBits)

		outFile.Write(tmp)

		// pcm := pcm[1:]
		// pcm = append(pcm, tmpFloat)

		// _, max := calcAvrMax(pcm)
		min := float32(-0.5)
		// average := float32(0)
		max := float32(0.5)

		if tmpFloat < min || tmpFloat > max {
			continue
		}

		difference := tmpFloat - min
		value := difference
		out := make([]byte, int(value*100))

		for i := range out {
			out[i] = ' '
		}

		rgb := fadeGrad.At(float64(value))

		// fmt.Print("\033[H\033[2J")

		truecolor.
			White().
			Background(uint8(rgb.R*100), uint8(rgb.G*100), uint8(rgb.B*100)).
			Print(fmt.Sprintf("%s%f\n", string(out), tmpFloat))

		time.Sleep(10 * time.Millisecond)
	}
}

func calcAvrMax(inp []float32) (float32, float32) {
	var sum float32 = 0
	var max float32 = 0.01

	for _, f := range inp {
		sum += f

		if f > max {
			f = max
		}
	}

	return sum / float32(len(inp)), max
}

func convToFloat(inp []byte) []float32 {
	out := make([]float32, len(inp)/4)
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
