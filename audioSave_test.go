package main

import (
	"os"
	"testing"
)

func TestSaveAudio(t *testing.T) {
	t.Log("Recording audio")

	out, err := ReadAudio(8192)
	if err != nil {
		t.Fatal("Error while reading the audio: ", err)
	}

	var filename = "examples/audio/349.228hz.pcm"

	err = os.WriteFile(filename, out, 0644)
	if err != nil {
		t.Fatal("Error while saving the audio: ", err)
	}

	t.Log("Audio saved at ", filename)
}
