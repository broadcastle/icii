package icii

import (
	"io"

	"github.com/tcolgate/mp3"
)

// Data holds information about file to be played.
type Data struct {
	Bitrate    int64
	Samples    int
	SampleRate int
	Channels   int
	Frames     []mp3.Frame
	Skipped    int
	Reader     io.Reader
}

// GetData gets the information needed from a io.Reader.
func GetData(d io.Reader) (Data, error) {

	var data Data

	return data, nil

}
