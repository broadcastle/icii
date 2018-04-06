package stream

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"broadcastle.co/code/icii/pkg/frames"
	"broadcastle.co/code/icii/pkg/info"
)

// Buffer is how many seconds should be buffered.
var Buffer = 3

// Stream holds the server information.
type Stream struct {
	Name        string
	Description string
	Mount       string
	User        string
	Password    string
	Host        string
	Genre       string
	Website     string
}

// Create a stream using the supplied information.
//
// Create returns an error if something went wrong during the creation process.
func Create(name, description, genre, host, website, mountpoint, user, password string) (Stream, error) {

	if mountpoint[0:1] == "/" {
		mountpoint = mountpoint[1:]
	}

	s := Stream{
		Name:        name,
		Description: description,
		Mount:       mountpoint,
		User:        user,
		Password:    password,
		Host:        host,
		Website:     website,
	}

	return s, nil
}

// File sends the file.
func (s Stream) File(filename string) error {

	// GET FILE INFORMATION
	i, err := info.GetData(filename)
	if err != nil {
		log.Println("GetData error")
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println("os.OpenError in s.File")
		return err
	}

	framesSent := 0
	timeBegin := time.Now()

	// LOOP THROUGH FRAMES
	for framesSent < len(i.Frames) {

		sendBegin := time.Now()

		buffer, err := frames.Get(*file, i.FramesToRead)
		if err != nil {
			log.Println("frames.Get")
			return err
		}

		// SEND THE FILE

		if err := s.sendBytes(buffer, i); err != nil {
			log.Println("s.sendBytes")
			return err
		}

		framesSent = framesSent + i.FramesToRead

		bufferSent := i.BufferSent(framesSent, timeBegin)

		pause := timePause(bufferSent, sendBegin)

		time.Sleep(pause)

	}

	sleep := i.TimeBetweenTracks(timeBegin)

	time.Sleep(sleep)

	return nil
}

// SendBytes sends the byte to the stream.
func (s Stream) sendBytes(b []byte, d info.Data) error {

	link := s.Host + "/" + s.Mount

	_, err := url.ParseRequestURI(link)
	if err != nil {
		return err
	}

	client := &http.Client{}
	r := bytes.NewReader(b)

	req, err := http.NewRequest("PUT", link, r)
	if err != nil {
		return err
	}

	audioInfo := "samplerate=" + strconv.Itoa(d.SampleRate) + ";channels=" + strconv.Itoa(d.Channels)

	req.SetBasicAuth(s.User, s.Password)
	req.Header.Set("content-type", "audio/mpeg")
	req.Header.Set("ice-public", "0")
	req.Header.Set("ice-name", s.Name)
	req.Header.Set("ice-description", s.Description)
	req.Header.Set("ice-genre", s.Genre)
	req.Header.Set("ice-website", s.Website)
	req.Header.Set("ice-audio-info", audioInfo)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(resp)

	defer resp.Body.Close()

	return nil

}

func timePause(sent int, send time.Time) time.Duration {

	x := 975

	switch {
	case sent < (Buffer - 100):
		x = 900
	case sent > Buffer:
		x = 1100
	}

	lag := float64((time.Now().Sub(send)).Seconds()) * 1000
	pause := x - int(lag)

	return time.Duration(time.Millisecond) * time.Duration(pause)

}
