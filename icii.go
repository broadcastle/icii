package icii

import "net"

// Server holds the server information.
type Server struct {
	net.Conn
}

// Stream holds the stream information.
type Stream struct {
}
