package icii

import (
	"io"
	"net"
	"os"
)

// Stream holds the server information.
type Stream struct {
	Name        string
	Description string
	Connection  net.Conn
}

// Create a stream using the supplied information.
//
// Create returns an error if something went wrong during the creation process.
func Create(name, description, host, mountpount string) (Stream, error) {

	s := Stream{
		Name:        name,
		Description: description,
	}

	connection, err := connect(host)
	if err != nil {
		return s, err
	}

	s.Connection = connection

	return s, nil
}

// File receives a filename to play on the stream.
func (s Stream) File(filename string) error {

	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	return s.Reader(r)
}

// Reader receive data from a io.Reader and plays it.
func (s Stream) Reader(r io.Reader) error {

	data, err := GetData(r)
	if err != nil {
		return err
	}

	return s.Send(data)
}

// Send the mp3 data to the stream.
func (s Stream) Send(data Data) error {
	return nil
}
