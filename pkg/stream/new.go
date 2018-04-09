package stream

import (
	"encoding/base64"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"broadcastle.co/code/icii/pkg/mpeg"
)

// Config holds the information for a stream.
type Config struct {
	Host        string
	Port        int
	Mount       string
	User        string
	Password    string
	Name        string
	URL         string
	Genre       string
	Description string
	Abort       bool
	Connection  net.Conn
	BufferSize  int
	Connected   bool
}

var version = 1

// File streams a file from filename.
func (c Config) File(filename string) error {

	info, err := mpeg.GetInfo(filename)
	if err != nil {
		return err
	}

	if err := c.Connect(info); err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return c.Stop(err)
	}

	defer file.Close()

	timeBegin := time.Now()
	framesSent := 0

	mpeg.SeekTo1StFrame(*file)

	for framesSent < info.NumberOfFrames {

		sendBegin := time.Now()

		// Get the frames to be sent.
		buf, err := mpeg.GetFrames(*file, info.FramesToRead)
		if err != nil {
			return c.Stop(err)
		}

		// Send the frames.
		if err := c.Send(buf); err != nil {
			return c.Stop(err)
		}

		framesSent = framesSent + info.FramesToRead

		pause := timePause(timeBegin, sendBegin, framesSent, info.SPF, info.SampleRate, c.BufferSize)

		if c.Abort {
			err := errors.New("aborted by user")
			return c.Stop(err)
		}

		time.Sleep(pause)

	}

	return nil
}

// Connect connects to the server in the config.
func (c *Config) Connect(info mpeg.Info) error {
	// return errors.New("c.Connect is not written")

	if c.Connected {
		return nil
	}

	host := c.Host + ":" + strconv.Itoa(c.Port)
	sock, err := net.Dial("tcp", host)
	if err != nil {
		return c.Stop(err)
	}

	c.Connection = sock

	time.Sleep(time.Second)

	headers := "SOURCE /" + c.Mount + " HTTP/1.0\r\n" +
		"Content-Type: audio/mpeg\r\n" +
		"Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte(c.User+":"+c.Password)) + "\r\n" +
		"User-Agent: goicy/" + strconv.Itoa(version) + "\r\n" +
		"ice-name: " + c.Name + "\r\n" +
		"ice-public: 0\r\n" +
		"ice-url: " + c.URL + "\r\n" +
		"ice-genre: " + c.Genre + "\r\n" +
		"ice-description: " + c.Description + "\r\n" +
		"ice-audio-info: bitrate=" + strconv.Itoa(int(info.BitRate)) +
		";channels=" + strconv.Itoa(info.Channels) +
		";samplerate=" + strconv.Itoa(info.SampleRate) + "\r\n\r\n"

	if err := c.Send([]byte(headers)); err != nil {
		return c.Stop(err)
	}

	time.Sleep(time.Second)

	res, err := c.Receive()
	if err != nil {
		return c.Stop(err)
	}

	if string(res[9:12]) != "200" {
		return c.Stop(errors.New("invalid icecast response: " + string(res)))
	}

	c.Connected = true

	return nil
}

// Send the buffer to the Connection in c.
func (c *Config) Send(buffer []byte) error {

	n, err := c.Connection.Write(buffer)
	if err != nil {
		return c.Stop(err)
	}

	if n != len(buffer) {
		return c.Stop(errors.New("c.Send() error"))
	}

	return nil
}

// Receive returns the response from a c.Connection.
func (c Config) Receive() ([]byte, error) {
	// return nil, errors.New("c.Receive is not written")

	buffer := make([]byte, 1024)
	n, err := c.Connection.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[0:n], nil

}

// Stop the stream from broadcasting.
func (c *Config) Stop(err error) error {

	if err == nil {
		err = errors.New("stream was stopped without a error")
	}

	log.Println(err)

	c.Connection.Close()
	c.Connected = false

	return err

}

func timePause(timeBegin time.Time, sendBegin time.Time, framesSent int, spf int, samplerate int, bufferSize int) time.Duration {

	elapsed := float64((time.Now().Sub(timeBegin)).Seconds()) * 1000
	lag := int(float64((time.Now().Sub(sendBegin)).Seconds()) * 1000)
	sent := float64(framesSent) * float64(spf) / float64(samplerate) * 1000

	pause := 0
	bufferSent := 0

	if sent > elapsed {
		bufferSent = int(sent) - int(elapsed)
	}

	switch {
	case bufferSent < (bufferSize - 100):
		pause = 800 - lag
	case bufferSent > bufferSize:
		pause = 1100 - lag
	default:
		pause = 975 - lag
	}

	return time.Duration(time.Millisecond) * time.Duration(pause)

}
