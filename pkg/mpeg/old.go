package mpeg

import (
	"io"
	"os"
)

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

//// NEED TO REMOVE

// SeekTo1StFrame ...
func SeekTo1StFrame(f os.File) int64 {

	buf := make([]byte, 100000)
	f.ReadAt(buf, 0)

	// skip ID3V2 at the beginning of file
	var ID3Length int64
	for id3 := string(buf[0:3]); id3 == "ID3"; {
		//major := byte(buf[4])
		//minor := byte(buf[5])
		//flags := buf[6]
		ID3Length = ID3Length + (int64(buf[6])<<21 | int64(buf[7])<<14 | int64(buf[8])<<7 | int64(buf[9])) + 10
		f.ReadAt(buf, ID3Length)
		id3 = string(buf[0:3])
	}

	pos := int64(-1)

	for i := 0; i < len(buf); i++ {
		if (buf[i] == 0xFF) && ((buf[i+1] & 0xE0) == 0xE0) {
			if len(buf)-i < 10 {
				break
			}
			mpxHeader := buf[i : i+4]
			if _, ok := isValidFrameHeader(mpxHeader); ok {
				if framelength := getFrameSize(mpxHeader); framelength > 0 {
					if i+framelength+4 > len(buf) {
						break
					}
					mpxHeader = buf[i+framelength : i+framelength+4]
					if _, ok := isValidFrameHeader(mpxHeader); ok {
						pos = int64(i) + ID3Length
						f.Seek(pos, 0)
						break
					}
				}
			}
		}
	}
	return pos
}

// REMOVE

func reSync(f os.File) bool {

	f.Seek(-3, 1)
	return false

}

// GetFrames ...
func GetFrames(f os.File, framesToRead int) ([]byte, error) {
	var buf []byte
	var err error

	framesRead, bytesRead := 0, 0
	headers := make([]byte, 4)
	headers2 := make([]byte, 4)
	numBytesToRead := 0

	for framesRead < framesToRead {

		bytesRead, err = f.Read(headers)
		if err != nil && err != io.EOF {
			return nil, err
		}

		//input file has ended
		if bytesRead < len(headers) {
			break
		}
		framelength := getFrameSize(headers)

		numBytesToRead = 0
		if framelength > len(headers) {
			numBytesToRead = framelength - len(headers) // + crc
		}

		oldpos, _ := f.Seek(0, 1)
		br, _ := f.Seek(int64(numBytesToRead), 1)
		bytesRead = int(br - oldpos)

		bbr, _ := f.Read(headers2)
		bytesRead = bytesRead + int(bbr)

		f.Seek(int64(-bytesRead), 1)

		// copy frame header to out buffer
		buf = append(buf, headers...)

		// read raw frame data
		lbuf := make([]byte, numBytesToRead)
		bytesRead, err = f.Read(lbuf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		buf = append(buf, lbuf[0:bytesRead]...)
		if bytesRead < numBytesToRead {
			break
		}

		framesRead = framesRead + 1
	}
	return buf, nil
}

// func GetFrames(f os.File, framesToRead int) ([]byte, error) {
// 	var buf []byte
// 	var err error

// 	framesRead, bytesRead := 0, 0
// 	headers := make([]byte, 4)
// 	headers2 := make([]byte, 4)
// 	// inSync := true
// 	// inSync := false
// 	numBytesToRead := 0

// 	for framesRead < framesToRead {

// 		bytesRead, err = f.Read(headers)
// 		if err != nil && err != io.EOF {
// 			return nil, err
// 		}

// 		//input file has ended
// 		if bytesRead < len(headers) {
// 			break
// 		}

// 		// // Check if frame header is valid, if not, go back.
// 		// if _, ok := isValidFrameHeader(headers); !ok {
// 		// 	// if inSync {
// 		// 	// 	pos, _ := f.Seek(0, 1)
// 		// 	// 	logger.Log("Bad MPEG frame at offset "+strconv.Itoa(int(pos-4))+
// 		// 	// 		", resyncing...", logger.LOG_DEBUG)
// 		// 	// }
// 		// 	// inSync = false
// 		// 	f.Seek(-3, 1)
// 		// 	continue
// 		// }

// 		framelength := getFrameSize(headers)
// 		// if framelength == 0 || framelength > 5000 {
// 		// 	// if inSync {
// 		// 	// 	pos, _ := f.Seek(0, 1)
// 		// 	// 	logger.Log("Bad MPEG frame at offset "+strconv.Itoa(int(pos-4))+
// 		// 	// 		", resyncing...", logger.LOG_DEBUG)
// 		// 	// }
// 		// 	// inSync = false
// 		// 	f.Seek(-3, 1)
// 		// 	continue
// 		// }

// 		numBytesToRead = 0
// 		if framelength > len(headers) {
// 			numBytesToRead = framelength - len(headers) // + crc
// 		}

// 		oldpos, _ := f.Seek(0, 1)
// 		br, _ := f.Seek(int64(numBytesToRead), 1)
// 		bytesRead = int(br - oldpos)

// 		bbr, _ := f.Read(headers2)
// 		bytesRead = bytesRead + int(bbr)

// 		f.Seek(int64(-bytesRead), 1)
// 		// if _, ok := isValidFrameHeader(headers2); !ok {
// 		// 	// if inSync {
// 		// 	// 	pos, _ := f.Seek(0, 1)
// 		// 	// 	logger.Log("Bad MPEG frame at offset "+strconv.Itoa(int(pos-4))+", resyncing...", logger.LOG_DEBUG)
// 		// 	// }
// 		// 	// inSync = false
// 		// 	f.Seek(-3, 1)
// 		// 	continue
// 		// }

// 		// // from now on, frame is considered valid
// 		// if !inSync {
// 		// 	pos, _ := f.Seek(0, 1)
// 		// 	logger.Log("Resynced at offset "+strconv.Itoa(int(pos-4)), logger.LOG_DEBUG)
// 		// }
// 		// inSync = true

// 		// copy frame header to out buffer
// 		buf = append(buf, headers...)

// 		// read raw frame data
// 		lbuf := make([]byte, numBytesToRead)
// 		bytesRead, err = f.Read(lbuf)
// 		if err != nil {
// 			if err != io.EOF {
// 				return nil, err
// 			}
// 		}
// 		buf = append(buf, lbuf[0:bytesRead]...)
// 		if bytesRead < numBytesToRead { // the // input // file // has // ended
// 			break
// 		}
// 		framesRead = framesRead + 1
// 	}
// 	return buf, nil
// }
