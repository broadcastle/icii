package ice

import (
	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
)

// Stream information
type Stream struct {
	database.Stream
}

// Create a stream
func (s *Stream) Create() error {
	return nil
}

// Update a stream
func (s *Stream) Update(i interface{}) error {
	return nil
}

// Delete a stream
func (s *Stream) Delete() error {
	return nil
}

// Get stream information
func (s *Stream) Get() error {
	return nil
}

// Echo gets the stream struct from the echo context.
func (s *Stream) Echo(c echo.Context) error {
	return nil
}
