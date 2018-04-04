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

		bufferSent := i.BufferSent(framesSent, timeBegin)

		pause := timePause(bufferSent, sendBegin)

		time.Sleep(pause)

	}

	sleep := i.TimeBetweenTracks(timeBegin)

	time.Sleep(sleep)

	return nil
}

func timePause(sent int, send time.Time) time.Duration {

	x := 975

	switch {
	case sent < (Buffer - 100):
		x = 900
	case sent > Buffer:
		x = 1100
	}

	lag := float64((time.Now().Sub(send)).Seconds()) * 1000
	pause := x - int(lag)

	return time.Duration(time.Millisecond) * time.Duration(pause)

}
