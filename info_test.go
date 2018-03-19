package icii

import (
	"os"
	"testing"
)

func TestGetData(t *testing.T) {

	file, err := os.Open("test_audio/loping_sting.mp3")
	if err != nil {
		t.Error(err)
	}

	data, err := GetData(file)
	if err != nil {
		t.Error(err)
	}

	if data.Skipped != 2206 {
		t.Errorf("skipped wrong amount of frames\nwas %x instead of %x", data.Skipped, 2206)
	}

	if int(data.Bitrate) != 320 {
		t.Errorf("wrong bitrate\nwas %x instead of %x", int(data.Bitrate), 320)
	}

	if data.SampleRate != 44100 {
		t.Errorf("wrong sample rate\nwas %x instead of %x", data.SampleRate, 44100)
	}

}
