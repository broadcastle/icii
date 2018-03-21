package icii

import (
	"net"
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
