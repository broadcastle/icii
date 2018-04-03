package frames

//// Taken from github.com/stunndard/goicy

import (
	"errors"
	"io"
	"log"
	"os"
)

// Get retrieves the frames from f.
func Get(f os.File, framesToRead int) (buf []byte, err error) {

	framesRead, bytesRead, bytesToRead := 0, 0, 0
	firstHeader, secondHeader := make([]byte, 4), make([]byte, 4)

	for framesRead < framesToRead {

		bytesRead, err = f.Read(firstHeader)
		if err != nil && err != io.EOF {
			return
		}

		if bytesRead < len(firstHeader) {
			break
		}

		if _, valid := isValidFrameHeader(firstHeader); !valid {
			if _, err := f.Seek(-3, 1); err != nil {
				return nil, err
			}
		}

		frameLength := getFrameSize(firstHeader)

		if frameLength == 0 || frameLength > 5000 {
			if _, err := f.Seek(-3, 1); err != nil {
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
			return nil, err
		}

		br, err := f.Seek(int64(bytesToRead), 1)
		if err != nil {
			return nil, err
		}

		bytesRead = int(br - posFirst)

		bbr, err := f.Read(secondHeader)
		if err != nil {
			return nil, err
		}

		if _, valid := isValidFrameHeader(secondHeader); !valid {
			if _, err := f.Seek(-3, 1); err != nil {
				return nil, err
			}
		}

		bytesRead = bytesRead + int(bbr)

		if _, err := f.Seek(int64(-bytesRead), 1); err != nil {
			return nil, err
		}

		buf = append(buf, firstHeader...)
		buf2 := make([]byte, bytesToRead)

		bytesRead, err = f.Read(buf2)
		if err != nil {
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

	f := frame{
		mpeg:    mpegver,
		layer:   layer,
		sri:     srindex,
		bri:     brindex,
		padding: padding,
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

	return f.size
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

type frame struct {
	br      uint32 // bitrate
	bri     byte   // bitrate index
	layer   byte   // layer version
	mpeg    byte   // mpeg version
	size    int    // frame size
	sr      uint32 // sample rate
	sri     byte   // sample rate index
	padding int    // padding
	version int    // used to determine mpeg & layer combo
}

func (f *frame) findVersion() error {

	x := 0

	if f.layer > 3 || f.layer < 1 || f.mpeg > 3 || f.mpeg == 1 {
		return errors.New("could not find version, missing mpeg and layer info")
	}

	switch f.mpeg {
	case 3:
		x = 2
	case 2:
		x = 1
	default:
		x = 0
	}

	f.version = (x * 3) + int(f.layer)

	return nil

}

func (f *frame) findBitrate() error {

	var b byte

	brtable := [...]uint32{
		0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0,
		0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0,
		0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0,
		0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0,
		0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}

	switch f.version {
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

	f.br = brtable[f.bri+b] * 1000

	return nil

}

func (f *frame) findSampleRate() error {

	srtable := [...]uint32{
		44100, 48000, 32000, 0, // mpeg1
		22050, 24000, 16000, 0, // mpeg2
		11025, 12000, 8000, 0} // mpeg2.5

	var n byte

	switch f.version {
	case 9, 8, 7: // mpeg 0 | layers 1,2,3
		n = 0
	case 6, 5, 4: // mpeg 2 | layers 1,2,3
		n = 4
	case 3, 2, 1: // mpeg 3 | layers 1,2,3
		n = 8
	default:
		return errors.New("unable to find sample rate: frame version is missing")
	}

	f.sr = srtable[f.sri+n]
	return nil

}

func (f *frame) calculateSize() error {

	if f.br == 0 || f.sr == 0 || f.padding == 0 {
		return errors.New("unable to calculate size: frame information is missing")
	}

	switch f.version {
	case 9:
		f.size = (int(12*f.br/f.sr) * 4) + (f.padding * 4)
	case 8, 7:
		f.size = int(144*f.br/f.sr) + f.padding
	case 6, 3:
		f.size = (int(12*f.br/f.sr) * 4) + (f.padding * 4)
	case 5, 2:
		f.size = int(144*f.br/f.sr) + f.padding
	case 4, 1:
		f.size = int(72*f.br/f.sr) + f.padding
	default:
		return errors.New("unable to calculate size: frame version is missing")
	}

	return nil
}
