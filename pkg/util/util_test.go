package util

import (
	"errors"
	"testing"
)

func TestValidFile(t *testing.T) {

	negative := ValidFile("./../../test_audio/non-existant.mp3")
	if negative {
		t.Error(errors.New("unable to detect file does not exists"))
	}

	valid := ValidFile("./../../test_audio/loping_sting.mp3")
	if !valid {
		t.Error(errors.New("unable to read valid file"))
	}

}
