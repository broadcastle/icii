package stream

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"broadcastle.co/code/icii/pkg/config"
	"broadcastle.co/code/icii/pkg/logger"
	"broadcastle.co/code/icii/pkg/mpeg"
	"broadcastle.co/code/icii/pkg/network"
)

// Abort tells us if we should abort.
var Abort bool

// File streams filename to icecast.
func File(filename string) error {

	var sock net.Conn

	cleanUp := func(err error) {
		log.Println(err)
		network.Close(sock)
	}

	logger.Log("Checking file: "+filename+"...", logger.LOG_INFO)

	i, err := mpeg.GetInfo(filename)
	if err != nil {
		return err
	}

	sock, err = network.ConnectServer(config.Cfg.Host, config.Cfg.Port, i.BitRate, i.SampleRate, i.Channels)
	if err != nil {
		logger.Log("Cannot connect to server", logger.LOG_ERROR)
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		cleanUp(err)
		return err
	}

	defer f.Close()

	// if config.Cfg.UpdateMetadata {
	// 	go metadata.GetTagsFFMPEG(filename)
	// 	// cuefile := util.Basename(filename) + ".cue"
	// 	// cuesheet.Load(cuefile)
	// }

	logger.Log("Streaming file: "+filename+"...", logger.LOG_INFO)
	logger.TermLn("CTRL-C to stop", logger.LOG_INFO)

	// get number of frames to read in 1 iteration
	timeBegin := time.Now()

	// OLD BUT WORKS
	mpeg.SeekTo1StFrame(*f)
	framesSent := 0

	for framesSent < i.NumberOfFrames {
		sendBegin := time.Now()

		lbuf, err := mpeg.GetFrames(*f, i.FramesToRead)
		if err != nil {
			logger.Log("Error reading data stream", logger.LOG_ERROR)
			cleanUp(err)
			return err
		}

		if err := network.Send(sock, lbuf); err != nil {
			cleanUp(err)
			logger.Log("Error sending data stream", logger.LOG_ERROR)
			return err
		}

		framesSent = framesSent + i.FramesToRead

		timeElapsed := int(float64((time.Now().Sub(timeBegin)).Seconds()) * 1000)
		timeSent := int(float64(framesSent) * float64(i.SPF) / float64(i.SampleRate) * 1000)

		bufferSent := 0
		if timeSent > timeElapsed {
			bufferSent = timeSent - timeElapsed
		}

		if timeElapsed > 1500 {
			logger.Term("Frames: "+strconv.Itoa(framesSent)+"/"+strconv.Itoa(i.NumberOfFrames)+"  Time: "+
				strconv.Itoa(timeElapsed/1000)+"/"+strconv.Itoa(timeSent/1000)+"s  Buffer: "+
				strconv.Itoa(bufferSent)+"ms  Frames/Bytes: "+strconv.Itoa(i.FramesToRead)+"/"+strconv.Itoa(len(lbuf)), logger.LOG_INFO)
		}

		timePause := sendLag(sendBegin, bufferSent, config.Cfg.BufferSize)

		if Abort {
			err := errors.New("aborted by user")
			cleanUp(err)
			return err
		}

		time.Sleep(time.Duration(time.Millisecond) * time.Duration(timePause))
	}

	// pause to clear up the buffer
	pause := i.TimeBetweenTracks(timeBegin)
	logger.Log("Pausing for "+pause.String()+"ms...", logger.LOG_DEBUG)
	time.Sleep(pause)

	return nil
}

// func getBuffer(bufs io.Reader...) ([]byte, error) {
// 	var b []byte
// 	for _, d := range bufs {
// 		f, err := ioutil.ReadAll(d)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return b, nil
// }

func sendLag(b time.Time, bufferSent int, bufferSize int) int {

	lag := int(float64((time.Now().Sub(b)).Seconds()) * 1000)

	switch {
	case bufferSent < (bufferSize - 100):
		return 800 - lag
	case bufferSent > bufferSize:
		return 1100 - lag
	default:
		return 975 - lag
	}

}
