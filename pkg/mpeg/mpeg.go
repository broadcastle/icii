package mpeg

// FIXED

import (
	"os"
	"time"

	"github.com/tcolgate/mp3"
)

// Info holds the information.
// Don't print the whole thing. It's going to be bad.
type Info struct {
	BitRate        float64
	SPF            int
	SampleRate     int
	Channels       int
	NumberOfFrames int
	Skipped        int
	FramesToRead   int
	Frames         []mp3.Frame
}

// GetInfo get the needed information about filename.
func GetInfo(filename string) (info Info, e error) {

	r, err := os.Open(filename)
	if err != nil {
		e = err
		return
	}

	d := mp3.NewDecoder(r)
	skipped := 0

	var first mp3.Frame

	for {
		s := 0
		if err := d.Decode(&first, &s); err != nil {
			if err.Error() != "EOF" {
				e = err
				return
			}
			break
		}

		info.Frames = append(info.Frames, first)

		if s > skipped {
			skipped = s
		}

	}

	info.BitRate = float64(first.Header().BitRate()) / 1000
	info.SampleRate = int(first.Header().SampleRate())
	info.SPF = first.Samples()
	info.NumberOfFrames = len(info.Frames)
	info.Skipped = skipped
	info.FramesToRead = (info.SampleRate / info.SPF) + 1

	switch first.Header().ChannelMode().String() {
	case "Stereo":
		info.Channels = 2
	case "Joint Stereo":
		info.Channels = 1
	case "Dual Channel":
		info.Channels = 2
	case "Single Channel":
		info.Channels = 1
	}

	return
}

// TimeBetweenTracks returns the time to pause between tracks.
func (i Info) TimeBetweenTracks(began time.Time) time.Duration {

	f := float64(i.NumberOfFrames)
	spf := float64(i.SPF)
	sr := float64(i.SampleRate)
	s := float64((time.Now().Sub(began)).Seconds())

	timeBetweenTracks := int(((f*spf)/sr)*1000) - int(s*1000)
	return time.Duration(time.Millisecond) * time.Duration(timeBetweenTracks)

}
