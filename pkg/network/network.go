package network

import (
	"encoding/base64"
	"errors"
	"net"
	"strconv"
	"time"

	"broadcastle.co/code/icii/pkg/config"
)

// Connected returns if were connected to the host.
var Connected = false
var csock net.Conn

// Connect connects to the host:port and returns a connection.
func Connect(host string, port int) (net.Conn, error) {
	h := host + ":" + strconv.Itoa(int(port))
	sock, err := net.Dial("tcp", h)
	if err != nil {
		Connected = false
	}
	return sock, err
}

// Send buf to the sock.
func Send(sock net.Conn, buf []byte) error {
	n, err := sock.Write(buf)
	if err != nil {
		Connected = false
		return err
	}
	if n != len(buf) {
		Connected = false
		return errors.New("Send() error")
	}
	return nil
}

// Recv returns data from sock.
func Recv(sock net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)

	n, err := sock.Read(buf)
	//fmt.Println(n, err, string(buf), len(buf))
	if err != nil {
		// logger.Log(err.Error(), logger.LOG_ERROR)
		return nil, err
	}
	return buf[0:n], err
}

// Close gets rid of sock.
func Close(sock net.Conn) {
	Connected = false
	sock.Close()
}

// ConnectServer connects to the server and returns a connection.
func ConnectServer(host string, port int, br float64, sr, ch int) (net.Conn, error) {
	var sock net.Conn

	if Connected {
		return csock, nil
	}

	if config.Cfg.ServerType == "shoutcast" {
		port++
	}
	// logger.Log("Connecting to "+config.Cfg.ServerType+" at "+host+":"+strconv.Itoa(port)+"...", logger.LOG_DEBUG)
	sock, err := Connect(host, port)

	if err != nil {
		Connected = false
		return sock, err
	}

	//fmt.Println("connected ok")
	time.Sleep(time.Second)

	headers := ""
	bitrate := int(br)
	samplerate := sr
	channels := ch
	contenttype := "audio/mpeg"

	headers = "SOURCE /" + config.Cfg.Mount + " HTTP/1.0\r\n" +
		"Content-Type: " + contenttype + "\r\n" +
		"Authorization: Basic " + base64.StdEncoding.EncodeToString([]byte("source:"+config.Cfg.Password)) + "\r\n" +
		"User-Agent: goicy/" + config.Version + "\r\n" +
		"ice-name: " + config.Cfg.StreamName + "\r\n" +
		"ice-public: 0\r\n" +
		"ice-url: " + config.Cfg.StreamURL + "\r\n" +
		"ice-genre: " + config.Cfg.StreamGenre + "\r\n" +
		"ice-description: " + config.Cfg.StreamDescription + "\r\n" +
		"ice-audio-info: bitrate=" + strconv.Itoa(bitrate) +
		";channels=" + strconv.Itoa(channels) +
		";samplerate=" + strconv.Itoa(samplerate) + "\r\n" +
		"\r\n"

	if err := Send(sock, []byte(headers)); err != nil {
		// logger.Log("Error sending headers", logger.LOG_ERROR)
		Connected = false
		return sock, err
	}

	time.Sleep(time.Second)
	resp, err := Recv(sock)
	if err != nil {
		Connected = false
		return sock, err
	}
	if string(resp[9:12]) != "200" {
		Connected = false
		return sock, errors.New("Invalid Icecast response: " + string(resp))
	}

	// logger.Log("Server connect successful", logger.LOG_INFO)
	Connected = true
	csock = sock

	return sock, nil
}
