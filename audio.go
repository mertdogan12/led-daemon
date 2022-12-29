package main

import (
	"github.com/mertdogan12/pulse-simple"
)

func ReadAudio(lenght int) ([]byte, error) {
	out := make([]byte, lenght)

	ss := pulse.SampleSpec{
		Format:   pulse.SAMPLE_FLOAT32LE,
		Rate:     48000,
		Channels: 1,
	}
	stream, err := pulse.Capture("led", "music", "bluez_output.E0_9D_FA_E3_56_CA.1", &ss)

	if err != nil {
		return out, err
	}

	_, err = stream.Read(out)
	if err != nil {
		return out, err
	}

	return out, nil
}
