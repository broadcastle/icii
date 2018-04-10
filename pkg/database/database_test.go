package database

import "testing"

func TestConnect(t *testing.T) {

	_, err := Config{Temp: true}.Connect()
	if err != nil {
		t.Error(err)
	}

	_, err = Config{}.Connect()
	if err == nil {
		t.Error("did not catch empty config")
	}

	_, err = Config{Postgres: true}.Connect()
	if err == nil {
		t.Error("did not catch empty config")
	}

}
