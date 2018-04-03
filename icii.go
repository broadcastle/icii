package icii

import (
	"net"
	"net/http"
)

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
	Connection  net.Conn
	Req         *http.Request
}

// // Create a stream using the supplied information.
// //
// // Create returns an error if something went wrong during the creation process.
// func Create(name, description, host, mountpount, user, password string) (Stream, error) {

// 	s := Stream{
// 		Name:        name,
// 		Description: description,
// 	}

// 	connection, err := connect(host)
// 	if err != nil {
// 		return s, err
// 	}

// 	s.Connection = connection

// 	return s, nil
// }

// // File receives a filename to play on the stream.
// func (s Stream) File(filename string) error {

// 	data, err := GetData(filename)
// 	if err != nil {
// 		return err
// 	}

// 	return s.Send(data)
// }

// // Send the mp3 data to the stream.
// func (s Stream) Send(data Data) error {
// 	return nil
// }
