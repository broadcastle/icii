package icii

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func connect(host string) (net.Conn, error) {

	link, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}

	u := link.Hostname() + ":" + link.Port()
	t := time.Duration(10 * time.Second)

	conn, err := net.DialTimeout("tcp", u, t)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *Stream) sendBytes(b []byte, d Data) error {

	r := bytes.NewReader(b)

	return s.sendReader(r, d)

}

func (s *Stream) sendReader(r io.Reader, d Data) error {

	// Delete the forward slash if it exists.
	if s.Mount[0:1] == "/" {
		s.Mount = s.Mount[1:]
	}

	host := s.Host + "/" + s.Mount

	if _, err := url.ParseRequestURI(host); err != nil {
		return err
	}

	client := &http.Client{Timeout: time.Duration(10 * time.Second)}

	req, err := http.NewRequest("PUT", host, r)
	if err != nil {
		return err
	}

	req.SetBasicAuth(s.User, s.Password)
	req.Header.Set("content-type", "audio/mpeg")
	req.Header.Set("ice-public", "0")
	req.Header.Set("ice-name", s.Name)
	req.Header.Set("ice-description", s.Description)
	req.Header.Set("ice-genre", s.Genre)
	req.Header.Set("ice-website", s.Website)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	time.Sleep(time.Second * 10)

	fmt.Println(res)

	return nil

}
