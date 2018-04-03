package stream

import (
	"log"
	"os"
	"time"

	"broadcastle.co/code/icii/pkg/frames"
	"broadcastle.co/code/icii/pkg/info"
)

// Buffer is how many seconds should be buffered.
var Buffer = 3

// File streams the current file
func File(filename string) error {

	// GET FILE INFORMATION
	i, err := info.GetData(filename)
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	framesSent := 0
	timeBegin := time.Now()

	// LOOP THROUGH FRAMES
	for framesSent < len(i.Frames) {

		sendBegin := time.Now()

		_, err := frames.Get(*file, i.FramesToRead)
		if err != nil {
			log.Println(err)
			return err
		}

		// SEND THE FILE

		framesSent = framesSent + i.FramesToRead

		timeElapsed := int(float64((time.Now().Sub(timeBegin)).Seconds()) * 1000)
		timeSent := int(float64(framesSent) * float64(i.Samples) / float64(i.SampleRate) * 1000)

		bufferSent := 0
		if timeSent > timeElapsed {
			bufferSent = timeSent - timeElapsed
		}

		sendLag := int(float64((time.Now().Sub(sendBegin)).Seconds()) * 1000)

		x := 975

		switch {
		case bufferSent < (Buffer - 100):
			x = 900
		case bufferSent > Buffer:
			x = 1100
		}

		timePause := x - sendLag

		time.Sleep(time.Duration(time.Millisecond) * time.Duration(timePause))

	}

	sleep := i.TimeBetweenTracks(timeBegin)

	time.Sleep(sleep)

	return nil
}
