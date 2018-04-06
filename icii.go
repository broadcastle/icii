package main

import (
	"log"

	"broadcastle.co/code/icii/pkg/stream"
)

func main() {

	s, err := stream.Create("Sample Stream", "A sample radio station", "Music", "http://192.168.1.227:8080", "http://192.168.1.227:8080", "mount.mp3", "source", "hackme")
	if err != nil {
		log.Panic(err)
	}

	if err := s.File("test_audio/loping_sting.mp3"); err != nil {
		log.Panic(err)
	}
}
