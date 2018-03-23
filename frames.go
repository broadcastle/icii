package icii

//// Taken from github.com/stunndard/goicy

import (
	"io"
	"os"
)

func getFrames(f os.File, framesToRead int) (buf []byte, err error) {

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
	var sr, bitrate uint32
	var res int

	brtable := [...]uint32{
		0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0,
		0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0,
		0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0,
		0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0,
		0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}
	srtable := [...]uint32{
		44100, 48000, 32000, 0, // mpeg1
		22050, 24000, 16000, 0, // mpeg2
		11025, 12000, 8000, 0} // mpeg2.5

	// get and check the mpeg version
	mpegver := byte((header[1] & 0x18) >> 3)
	if mpegver == 1 || mpegver > 3 {
		return 0
	}

	// get and check mpeg layer
	layer := byte((header[1] & 0x06) >> 1)
	if layer == 0 || layer > 3 {
		return 0
	}

	brindex := byte((header[2] & 0x0F0) >> 4)

	if mpegver == 3 && layer == 3 {
		// mpeg1, layer1
		bitrate = brtable[brindex]
	}
	if mpegver == 3 && layer == 2 {
		// mpeg1, layer2
		bitrate = brtable[brindex+16]
	}
	if mpegver == 3 && layer == 1 {
		// mpeg1, layer3
		bitrate = brtable[brindex+32]
	}
	if (mpegver == 2 || mpegver == 0) && layer == 3 {
		// mpeg2, 2.5, layer1
		bitrate = brtable[brindex+48]
	}
	if (mpegver == 2 || mpegver == 0) && (layer == 2 || layer == 1) {
		//mpeg2, layer2 or layer3
		bitrate = brtable[brindex+64]
	}
	bitrate = bitrate * 1000
	padding := int(header[2]&0x02) >> 1

	// get and check the 'sampling_rate_index':
	srindex := byte((header[2] & 0x0C) >> 2)
	if srindex >= 3 {
		return 0
	}
	if mpegver == 3 {
		sr = srtable[srindex]
	}
	if mpegver == 2 {
		sr = srtable[srindex+4]
	}
	if mpegver == 0 {
		sr = srtable[srindex+8]
	}

	switch mpegver {
	case 3: // mpeg1
		if layer == 3 { // layer1
			res = (int(12*bitrate/sr) * 4) + (padding * 4)
		}
		if layer == 2 || layer == 1 {
			// layer 2 & 3
			res = int(144*bitrate/sr) + padding
		}

	case 2, 0: //mpeg2, mpeg2.5
		if layer == 3 { // layer1
			res = (int(12*bitrate/sr) * 4) + (padding * 4)
		}
		if layer == 2 { // layer2
			res = int(144*bitrate/sr) + padding
		}
		if layer == 1 { // layer3
			res = int(72*bitrate/sr) + padding
		}
	}
	return res
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
