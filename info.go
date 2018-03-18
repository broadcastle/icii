package icii

import (
	"io"

	"github.com/tcolgate/mp3"
)

// Info holds information about file to be played.
type Info struct {
	Bitrate    int64
	Samples    int
	SampleRate int
	Channels   int
	Frames     []mp3.Frame
	Skipped    int
	Reader     io.Reader
}

// GetInfo retrieves the information from a file.
func GetInfo(filename string) (Info, error) {

	var info Info

	return info, nil
}
