package ice

import (
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"

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

	var results uint

	if err := db.Model(&Stream{}).Where("station_id = ? ", s.StationID).Count(&results).Error; err != nil {
		// log.Printf("stream.Delete(): %v", err)
		// log.WithFields(log.Fields{
		// 	"context": "database creation",
		// }).Warn(err)
		log.Info(err)
		return err
	}

	if int(results) < 2 {
		log.Warn("unable to delete stream")
		return errors.New("unable to delete last stream")
	}

	return db.Delete(&s).Error

}

// Get stream information
func (s *Stream) Get() error {
	return db.Where(&s).First(&s).Error
}

// Echo gets the stream struct from the echo context.
func (s *Stream) Echo(c echo.Context) error {

	i := c.Param("stream")

	id, err := strconv.Atoi(i)
	if err != nil {
		return err
	}

	st := c.Param("station")
	sid, err := strconv.Atoi(st)
	if err != nil {
		return err
	}

	s.ID = uint(id)
	s.StationID = uint(sid)

	return s.Get()
}
