package icii

import (
	"io"

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
}

// GetData gets the information needed from a io.Reader.
func GetData(d io.Reader) (Data, error) {

	var (
		data    Data
		frame   mp3.Frame
		skipped int
	)

	track := mp3.NewDecoder(d)

	// Loop through all the MP3 frames.
	for {
		s := 0
		if err := track.Decode(&frame, &s); err != nil {
			// Return all errors except EOF. End the loop on a EOF.
			if err.Error() != "EOF" {
				return Data{}, err
			}
			break
		}

		data.Frames = append(data.Frames, frame)

		if s > skipped {
			skipped = s
		}

	}

	// Technical Assignments
	data.Bitrate = float64(frame.Header().BitRate()) / 1000
	data.SampleRate = int(frame.Header().SampleRate())
	data.Samples = frame.Samples()
	data.FramesToRead = (data.SampleRate / data.Samples) + 1
	data.Skipped = skipped

	// Determine how many channels the track has by using a frame.
	switch frame.Header().ChannelMode().String() {
	case "Stereo":
		data.Channels = 2
	case "Joint Stereo":
		data.Channels = 1
	case "Dual Channel":
		data.Channels = 2
	case "Single Channel":
		data.Channels = 1
	}

	return data, nil

}
