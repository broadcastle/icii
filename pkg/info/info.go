package info

import (
	"errors"
	"io"
	"os"
	"time"

	"broadcastle.co/code/icii/pkg/util"
	"github.com/tcolgate/mp3"
)

// Data holds information about file to be played.
type Data struct {
	Bitrate      float64
	Channels     int
	Frames       []mp3.Frame
	FramesToRead int
	Reader       io.Reader
	SampleRate   int
	Samples      int
	Skipped      int
	Location     string
}

// GetData gets the information needed from a path.
func GetData(filename string) (Data, error) {

	var (
		data    Data
		frame   mp3.Frame
		skipped int
	)

	validFile := util.ValidFile(filename)
	if !validFile {
		return data, errors.New("not a valid file")
	}

	d, err := os.Open(filename)
	if err != nil {
		return data, err
	}

	track := mp3.NewDecoder(d)

	// Loop through all the MP3 frames.
	for {
		s := 0
		if err := track.Decode(&frame, &s); err != nil {
			// Return all errors except EOF. End the loop on a EOF.
			if err.Error() != "EOF" {
				return data, err
			}
			break
		}

		data.Frames = append(data.Frames, frame)

		if s > skipped {
			skipped = s
		}

	}

	// Determine how many channels the track has by using a frame.
	switch frame.Header().ChannelMode().String() {
	case "Stereo", "Dual Channel":
		data.Channels = 2
	case "Joint Stereo", "Single Channel":
		data.Channels = 1
	default:
		return Data{}, errors.New("unable to determine channels")
	}

	// Technical Assignments
	data.Bitrate = float64(frame.Header().BitRate()) / 1000
	data.SampleRate = int(frame.Header().SampleRate())
	data.Samples = frame.Samples()
	data.FramesToRead = (data.SampleRate / data.Samples) + 1
	data.Skipped = skipped
	data.Location = filename

	return data, nil

}

// TimeBetweenTracks calculates the time needed to clean the buffer between audio files.
func (d Data) TimeBetweenTracks(timeBegin time.Time) time.Duration {

	timeBetweenTracks := int(((float64(len(d.Frames))*float64(d.Samples))/float64(d.SampleRate))*1000) - int(float64((time.Now().Sub(timeBegin)).Seconds())*1000)

	return time.Duration(time.Millisecond) * time.Duration(timeBetweenTracks)
}
