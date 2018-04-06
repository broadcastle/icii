package frames

//// Taken from github.com/stunndard/goicy

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

// Get retrieves the frames from f.
func Get(f os.File, framesToRead int) (buf []byte, err error) {

	framesRead, bytesRead, bytesToRead := 0, 0, 0
	firstHeader, secondHeader := make([]byte, 4), make([]byte, 4)

	for framesRead < framesToRead {

		bytesRead, err = f.Read(firstHeader)
		if err != nil && err != io.EOF {
			log.Println("GetRead")
			return
		}

		if bytesRead < len(firstHeader) {
			break
		}

		if _, valid := isValidFrameHeader(firstHeader); !valid {
			if _, err := f.Seek(-3, 1); err != nil {
				log.Println("Get: isValidFrameHeader")
				return nil, err
			}
		}

		frameLength := getFrameSize(firstHeader)

		if frameLength == 0 || frameLength > 5000 {
			if _, err := f.Seek(-3, 1); err != nil {
				log.Println("Get: frameLength")
				return nil, err
			}
			continue
		}

		bytesToRead = 0
		if frameLength > len(firstHeader) {
			bytesToRead = frameLength - len(firstHeader)
		}

		posFirst, err := f.Seek(0, 1)
		if err != nil {
			log.Println("Get: 2")
			return nil, err
		}

		br, err := f.Seek(int64(bytesToRead), 1)
		if err != nil {
			log.Println("Get: 3")
			return nil, err
		}

		bytesRead = int(br - posFirst)

		bbr, err := f.Read(secondHeader)
		if err != nil {
			log.Println("Get: 4")
			return nil, err
		}

		if _, valid := isValidFrameHeader(secondHeader); !valid {
			if _, err := f.Seek(-3, 1); err != nil {
				log.Println("Get: 5")
				return nil, err
			}
		}

		bytesRead = bytesRead + int(bbr)

		f.Seek(int64(-bytesRead), 1)
		// if _, err := f.Seek(int64(-bytesRead), 1); err != nil {
		// 	log.Println("Get: 6")
		// 	return nil, err
		// }

		buf = append(buf, firstHeader...)
		buf2 := make([]byte, bytesToRead)

		bytesRead, err = f.Read(buf2)
		if err != nil {
			log.Println("Get: 7")
			return nil, err
		}

		buf = append(buf, buf2[0:bytesRead]...)

		if bytesRead < bytesToRead {
			break
		}

		framesRead = framesRead + 1

	}

	return
}

func getFrameSize(header []byte) int {

	mpegver := byte((header[1] & 0x18) >> 3)
	if mpegver == 1 || mpegver > 3 {
		return 0
	}

	layer := byte((header[1] & 0x06) >> 1)
	if layer == 0 || layer > 3 {
		return 0
	}
	srindex := byte((header[2] & 0x0C) >> 2)
	if srindex >= 3 {
		return 0
	}

	padding := int(header[2]&0x02) >> 1
	brindex := byte((header[2] & 0x0F0) >> 4)

	f := Frame{
		Mpeg:    mpegver,
		Layer:   layer,
		Sri:     srindex,
		Bri:     brindex,
		Padding: padding,
	}

	if err := f.findVersion(); err != nil {
		log.Println(err)
		return 0
	}

	if err := f.findBitrate(); err != nil {
		log.Println(err)
		return 0
	}

	if err := f.findSampleRate(); err != nil {
		log.Println(err)
		return 0
	}

	if err := f.calculateSize(); err != nil {
		log.Println(err)
		return 0
	}

	return f.Size
}

func isValidFrameHeader(header []byte) (int, bool) {

	if (header[0] != 0x0FF) && ((header[1] & 0x0E0) != 0x0E0) {
		return 0, false
	}

	// get and check the mpeg version
	mpegver := (uint16(header[1]) & 0x18) >> 3
	if mpegver == 1 || mpegver > 3 {
		return 0, false
	}

	// get and check mpeg layer
	layer := (header[1] & 0x06) >> 1
	if layer == 0 || layer > 3 {
		return 0, false
	}

	// get and check bitreate index
	brindex := (header[2] & 0x0F0) >> 4
	if brindex > 15 {
		return 0, false
	}

	// get and check the 'sampling_rate_index':
	srindex := (header[2] & 0x0C) >> 2
	if srindex >= 3 {
		return 0, false
	}

	return 0, true
}

// Frame has info
type Frame struct {
	Br      uint32 // bitrate
	Bri     byte   // bitrate index
	Layer   byte   // layer version
	Mpeg    byte   // mpeg version
	Size    int    // frame size
	Sr      uint32 // sample rate
	Sri     byte   // sample rate index
	Padding int    // padding
	Version int    // used to determine mpeg & layer combo
}

func (f *Frame) findVersion() error {

	x := 0

	if f.Layer > 3 || f.Layer < 1 || f.Mpeg > 3 || f.Mpeg == 1 {
		return errors.New("could not find version, missing mpeg and layer info")
	}

	switch f.Mpeg {
	case 3:
		x = 3
	case 2:
		x = 2
	default:
		x = 1
	}

	f.Version = (x * 3) + int(f.Layer)

	log.Println(f.Mpeg)
	log.Println(f.Layer)
	log.Println(f.Version)
	// log.Printf("mpeg: %x, %layer: %x, version: %x", f.Mpeg, f.Layer, f.Version)

	return nil

}

func (f *Frame) findBitrate() error {

	var b byte

	brtable := [...]uint32{
		0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0,
		0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0,
		0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0,
		0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0,
		0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}

	switch f.Version {
	case 9:
		b = 0
	case 8:
		b = 16
	case 7:
		b = 32
	case 6, 3:
		b = 48
	case 5, 4, 2, 1:
		b = 64
	default:
		return errors.New("unable to find bitrate: frame version is missing")
	}

	f.Br = brtable[f.Bri+b] * 1000

	return nil

}

func (f *Frame) findSampleRate() error {

	srtable := [...]uint32{
		44100, 48000, 32000, 0, // mpeg1
		22050, 24000, 16000, 0, // mpeg2
		11025, 12000, 8000, 0} // mpeg2.5

	var n byte

	switch f.Version {
	case 9, 8, 7: // mpeg 0 | layers 1,2,3
		n = 0
	case 6, 5, 4: // mpeg 2 | layers 1,2,3
		n = 4
	case 3, 2, 1: // mpeg 3 | layers 1,2,3
		n = 8
	default:
		return errors.New("unable to find sample rate: frame version is missing")
	}

	f.Sr = srtable[f.Sri+n]
	return nil

}

func (f *Frame) calculateSize() error {

	if f.Br == 0 || f.Sr == 0 || f.Padding == 0 {

		bitrate := strconv.Itoa(int(f.Br))
		samplerate := strconv.Itoa(int(f.Sr))
		padding := strconv.Itoa(f.Padding)

		log.Printf("bitrate %v, sample rate %v, padding %v", bitrate, samplerate, padding)
		// return errors.New("unable to calculate size: frame information is missing")
	}

	switch f.Version {
	case 9:
		f.Size = (int(12*f.Br/f.Sr) * 4) + (f.Padding * 4)
	case 8, 7:
		f.Size = int(144*f.Br/f.Sr) + f.Padding
	case 6, 3:
		f.Size = (int(12*f.Br/f.Sr) * 4) + (f.Padding * 4)
	case 5, 2:
		f.Size = int(144*f.Br/f.Sr) + f.Padding
	case 4, 1:
		f.Size = int(72*f.Br/f.Sr) + f.Padding
	default:
		return errors.New("unable to calculate size: frame version is missing")
	}

	return nil
}
