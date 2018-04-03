package util

import (
	"os"

	filetype "gopkg.in/h2non/filetype.v1"
)

// ValidFile returns whether a file is valid or not.
func ValidFile(filename string) bool {

	file, err := os.Open(filename)
	if err != nil {
		// log.Println(err)
		return false
	}

	header := make([]byte, 261)
	file.Read(header)

	return filetype.IsMIME(header, "audio/mpeg")
}
