package ice

import (
	"broadcastle.co/code/icii/pkg/stream"
	"github.com/sirupsen/logrus"
)

// Play will play the track in the stream.
func (s Stream) Play(t Track) error {

	config := stream.Config{
		Host:        s.Host,
		Port:        s.Port,
		Mount:       s.Mount,
		User:        s.User,
		Password:    s.Password,
		Name:        s.Name,
		URL:         s.URL,
		Genre:       s.Genre,
		Description: s.Description,
	}

	if err := config.Play(t.Location); err != nil {
		logrus.Warn(err)
		return err
	}

	return nil
}
