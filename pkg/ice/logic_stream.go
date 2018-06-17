package ice

import (
	"errors"
	"strconv"

	"broadcastle.co/code/icii/pkg/database"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

// Stream information
type Stream struct {
	database.Stream
}

// Create a stream
func (s *Stream) Create() error {

	if err := db.Create(&s).Error; err != nil {
		logrus.Warn(err)
		return err
	}

	log.Infof("icii create a new stream with id: %x", s.ID)

	return nil

}

// Update a stream
func (s *Stream) Update(i interface{}) error {

	info := i.(Stream)

	if err := db.Model(&s).Updates(info).Error; err != nil {
		logrus.Warn(err)
		return err
	}

	logrus.Infof("icii updated stream #%x with new information", s.ID)

	return nil
}

// Delete a stream
func (s *Stream) Delete() error {

	var results uint

	if err := db.Model(&Stream{}).Where("station_id = ? ", s.StationID).Count(&results).Error; err != nil {
		logrus.Warn(err)
		return err
	}

	if int(results) < 2 {
		logrus.Warn("unable to delete stream")
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
