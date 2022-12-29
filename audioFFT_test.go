package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/mjibson/go-dsp/fft"
)

func TestAnalyseAuidio(t *testing.T) {
	var filename = "examples/audio/349.228hz.pcm"

	t.Log("Reading file ", filename)

	out, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf(
			"Error while reading the file %s: %s\n\n%s",
			"You need to run the audio save test first, Run: go test audioSave_test",
			filename,
			err,
		)
	}

	if len(out) == 0 {
		t.Fatal("You need to runt the audio save test first, Run: go test audioSave_test.go")
	}

	t.Log("Converting file data to float64")

	data := make([]complex128, 0)
	tmp := make([]byte, 4)

	f2, err := os.Create("examples/pcm/349.228hz.csv")
	if err != nil {
		t.Fatal("File couldn't be created, ", err)
	}

	var i = 0
	for j, b := range out {
		if i > 3 {
			i = 0

			fBits := binary.LittleEndian.Uint32(tmp)
			f := math.Float32frombits(fBits)
			data = append(data, complex128(complex(f, 0)))

			f2.WriteString(fmt.Sprintf("%d,%f\n", j/3, f))
		}

		tmp[i] = b
		i++
	}

	t.Log("Converted file data saved at: examples/pcm/349.228hz.csv")

	f, err := os.Create("examples/fft/349.228hz.csv")
	if err != nil {
		t.Fatal("File couldn't be created, ", err)
	}

	w := bufio.NewWriter(f)

	t.Log("Creating fft data")

	x := fft.FFT(data)
	multi := 48000 / len(x)

	for i, y := range x {
		_, err := w.WriteString(fmt.Sprintf("%d,%f\n", i*multi, real(y)))
		if err != nil {
			t.Fatal("Error while writeing to out.raw, ", err)
		}
	}

	t.Log("FFT data saved at examples/fft/349.228hz.csv")
}
